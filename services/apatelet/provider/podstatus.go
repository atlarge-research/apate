package provider

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// GetPodStatus retrieves the status of a pod by label.
func (p *Provider) GetPodStatus(ctx context.Context, ns string, name string) (*corev1.PodStatus, error) {
	if p.Environment.DebugEnabled {
		log.Printf("GetPodStatus for %s/%s\n", ns, name)
	}

	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "context cancelled in GetPodStatus")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency in GetPodStatus")
		log.Println(err)
		return nil, err
	}

	pod, ok := p.Pods.GetPodByName(ns, name)
	if !ok {
		return nil, nil
	}

	podStatus, err := podResponse(responseArgs{ctx: ctx, provider: p, action: func() (interface{}, error) {
		status, err := (*p.Store).GetPodFlag(pod, events.PodStatus)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get pod status flag while getting pod status")
		}

		limitExceeded, err := p.doesPodExceedLimit(pod)
		if err != nil {
			return nil, errors.Wrap(err, "failed to determine if limit is exceeded while getting pod status")
		}

		if limitExceeded {
			return p.podFailed("Pod used too many resources and was then killed"), nil
		}

		switch status {
		case scenario.PodStatusPending:
			return p.podPending(), nil
		case scenario.PodStatusUnset:
			fallthrough // act as a normal pod
		case scenario.PodStatusRunning:
			return p.podRunning(), nil
		case scenario.PodStatusSucceeded:
			return p.podSucceeded(), nil
		case scenario.PodStatusFailed:
			return p.podFailed("Emulated pod has failed"), nil
		case scenario.PodStatusUnknown:
			fallthrough
		default:
			return p.podUnknown(), nil
		}
	}},
		pod,
		events.PodGetPodStatusResponse,
	)

	if IsExpected(err) {
		return nil, err
	}

	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to execute pod and node response (Get Pod Status)")
	}

	if status, ok := podStatus.(*corev1.PodStatus); ok {
		return status.DeepCopy(), nil
	}

	return nil, errors.Errorf("invalid podstatus %v", pod)
}

func (p *Provider) podPending() *corev1.PodStatus {
	return &corev1.PodStatus{
		Phase:   corev1.PodPending,
		Message: "Pod is awaiting further emulation instructions",
		Conditions: []corev1.PodCondition{
			{
				Type:               corev1.PodScheduled,
				Status:             corev1.ConditionTrue,
				LastProbeTime:      metav1.Time{Time: time.Now()},
				LastTransitionTime: metav1.Time{Time: time.Now()},
				Message:            "Pod is scheduled pod...",
			},
		},
	}
}

func (p *Provider) podUnknown() *corev1.PodStatus {
	return &corev1.PodStatus{
		Phase:   corev1.PodUnknown,
		Message: "Unknown emulated pod status",
	}
}

func (p *Provider) podRunning() *corev1.PodStatus {
	return &corev1.PodStatus{
		Phase:   corev1.PodRunning,
		Message: "Emulating pod successfully",
		Conditions: []corev1.PodCondition{
			{
				Type:               corev1.PodReady,
				Status:             corev1.ConditionTrue,
				LastProbeTime:      metav1.Time{Time: time.Now()},
				LastTransitionTime: metav1.Time{Time: time.Now()},
				Message:            "Emulating pod...",
			},
		},
	}
}

func (p *Provider) podSucceeded() *corev1.PodStatus {
	return &corev1.PodStatus{
		Phase:   corev1.PodSucceeded,
		Message: "Pod has completed successfully",
	}
}

func (p *Provider) podFailed(reason string) *corev1.PodStatus {
	return &corev1.PodStatus{
		Phase:   corev1.PodFailed,
		Message: reason,
		Conditions: []corev1.PodCondition{
			{
				Type:               corev1.PodReady,
				Status:             corev1.ConditionFalse,
				LastProbeTime:      metav1.Time{Time: time.Now()},
				LastTransitionTime: metav1.Time{Time: time.Now()},
				Message:            "Failed pod...",
			},
		},
	}
}

func (p *Provider) doesPodExceedLimit(pod *corev1.Pod) (bool, error) {
	limits := p.getPodResourceLimits(pod)

	podResourcesFlag, err := (*p.Store).GetPodFlag(pod, events.PodResources)
	if err != nil {
		return false, errors.Wrap(err, "failed to get pod resources flag while getting pod status")
	}

	podResources, ok := podResourcesFlag.(stats.PodStats)
	if !ok {
		return false, errors.Wrapf(err, "unable to convert '%v' to PodStats", podResourcesFlag)
	}

	usageCores := podResources.CPU.UsageNanoCores
	usageMemory := podResources.Memory.UsageBytes
	usageEphemeralStorage := podResources.EphemeralStorage.UsedBytes

	podExceedsPodLimit := *usageCores > limits.cpu || *usageMemory > limits.memory || *usageEphemeralStorage > limits.ephemeralStorage

	// If the total amount of all pods resources exceed the resources on the node, just kill the current one
	// TODO implement k8s OOM handling (much more complicated)
	nodeStats := p.Stats.statsSummary.Node

	totalLimitExceeded := *nodeStats.CPU.UsageNanoCores > uint64(p.Resources.CPU) ||
		*nodeStats.Memory.UsageBytes > uint64(p.Resources.Memory) ||
		*nodeStats.Fs.UsedBytes > uint64(p.Resources.EphemeralStorage)

	return podExceedsPodLimit || totalLimitExceeded, nil
}

func (p *Provider) getPodResourceLimits(pod *corev1.Pod) resources {
	totalCPU := uint64(0)
	totalMem := uint64(0)
	totalEphemeralStorage := uint64(0)

	for _, c := range pod.Spec.Containers {
		limits := c.Resources.Limits
		totalCPU += uint64(limits.Cpu().Value())
		totalMem += uint64(limits.Memory().Value())
		totalEphemeralStorage += uint64(limits.StorageEphemeral().Value())
	}

	return resources{
		totalCPU,
		totalMem,
		totalEphemeralStorage,
	}
}
