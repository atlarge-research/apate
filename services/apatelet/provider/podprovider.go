package provider

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"time"

	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/pkg/errors"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

const requestLoggingEnabled = false

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if requestLoggingEnabled {
		log.Printf("CreatePod %s/%s\n", pod.Namespace, pod.Name)
	}

	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in CreatePod")
	}

	return p.createOrUpdate(ctx, pod, events.PodCreatePodResponse)
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if requestLoggingEnabled {
		log.Printf("UpdatePod %s/%s\n", pod.Namespace, pod.Name)
	}

	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in UpdatePod")
	}

	return p.createOrUpdate(ctx, pod, events.PodUpdatePodResponse)
}

func (p *Provider) createOrUpdate(ctx context.Context, pod *corev1.Pod, pf events.PodEventFlag) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Create or Update)")
		log.Println(err)
		return err
	}

	_, err := podResponse(
		responseArgs{ctx, p, updateMap(p, pod)},
		getPodLabelByPod(pod),
		pf,
	)

	if IsExpected(err) {
		return err
	}

	err = errors.Wrap(err, "failed to execute pod and node response (Create or Update)")
	if err != nil {
		log.Println(err)
	}

	return err
}

func updateMap(p *Provider, pod *corev1.Pod) func() (interface{}, error) {
	return func() (interface{}, error) {
		p.Pods.AddPod(pod)
		return nil, nil
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	if requestLoggingEnabled {
		log.Printf("DeletePod %s/%s\n", pod.Namespace, pod.Name)
	}

	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in DeletePod")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency DeletePod")
		log.Println(err)
		return err
	}

	_, err := podResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			p.Pods.DeletePod(pod)
			return nil, nil
		}},
		getPodLabelByPod(pod),
		events.PodDeletePodResponse,
	)

	if IsExpected(err) {
		return err
	}

	err = errors.Wrap(err, "failed to execute pod and node response (Delete Pod)")
	if err != nil {
		log.Println(err)
	}

	return err
}

// GetPod retrieves a pod by label.
func (p *Provider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	if requestLoggingEnabled {
		log.Printf("GetPod %s/%s\n", namespace, name)
	}
	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "context cancelled in GetPod")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency GetPod")
		log.Println(err)
		return nil, err
	}

	label := p.getPodLabelByName(namespace, name)

	pod, err := podResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			pod, _ := p.Pods.GetPodByName(namespace, name)
			return pod, nil
		}},
		label,
		events.PodGetPodResponse,
	)

	if IsExpected(err) {
		return nil, err
	}

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pod)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if p, ok := pod.(*corev1.Pod); ok {
		return p, nil
	}

	return nil, errors.Errorf("invalid pod %v", pod)
}

// GetPodStatus retrieves the status of a pod by label.
func (p *Provider) GetPodStatus(ctx context.Context, ns string, name string) (*corev1.PodStatus, error) {
	if requestLoggingEnabled {
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

	label := p.getPodLabelByName(ns, name)

	pod, err := podResponse(responseArgs{ctx: ctx, provider: p, action: func() (interface{}, error) {
		status, err := (*p.Store).GetPodFlag(label, events.PodStatus)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get pod status flag while getting pod status")
		}

		limitExceeded, err := p.doesPodExceedLimit(ctx, ns, name, label)
		if err != nil {
			return nil, errors.Wrap(err, "failed to determine if limit is exceeded while getting pod status")
		}

		if limitExceeded {
			p.Pods.DeletePodByName(ns, name)

			return &corev1.PodStatus{
				Phase:   podStatusToPhase(scenario.PodStatusFailed),
				Message: "Pod used too many resources and was then killed",
				Conditions: []corev1.PodCondition{
					{
						Type:               corev1.PodReady,
						Status:             corev1.ConditionFalse,
						LastProbeTime:      metav1.Time{Time: time.Now()},
						LastTransitionTime: metav1.Time{Time: time.Now()},
						Message:            "Failed pod...",
					},
				},
			}, nil
		}

		return &corev1.PodStatus{
			Phase:   podStatusToPhase(status),
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
			ContainerStatuses: []corev1.ContainerStatus{}, //TODO: Implement this properly
		}, nil
	}},
		label,
		events.PodGetPodStatusResponse,
	)

	if IsExpected(err) {
		return nil, err
	}

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pod Status)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if status, ok := pod.(*corev1.PodStatus); ok {
		return status, nil
	}

	return nil, errors.Errorf("invalid podstatus %v", pod)
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	if requestLoggingEnabled {
		log.Printf("GetPods called\n")
	}

	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "context cancelled in GetPods")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency in GetPods")
		log.Println(err)
		return nil, err
	}

	pod, err := nodeResponse(responseArgs{ctx, p, func() (interface{}, error) {
		return p.Pods.GetAllPods(), nil
	}},
		events.NodeGetPodsResponse,
	)

	if IsExpected(err) {
		return nil, err
	}

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pods)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if pods, ok := pod.([]*corev1.Pod); ok {
		return pods, nil
	}

	return nil, errors.Errorf("invalid pods %v", pod)
}

// GetContainerLogs retrieves the log of a specific container.
func (p *Provider) GetContainerLogs(_ context.Context, ns, name, _ string, _ api.ContainerLogOpts) (io.ReadCloser, error) {
	if requestLoggingEnabled {
		log.Printf("GetContainerLogs for %s/%s\n", ns, name)
	}

	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte("This container is emulated by Apate\n"))), nil
}

// RunInContainer runs a command in a specific container.
func (p *Provider) RunInContainer(_ context.Context, ns, name, _ string, _ []string, _ api.AttachIO) error {
	if requestLoggingEnabled {
		log.Printf("RunInContainer for %s/%s\n", ns, name)
	}

	// There is no actual process running in the containers, so we can't do anything.
	return nil
}

func (p *Provider) runLatency(ctx context.Context) error {
	durationFlag, err := (*p.Store).GetNodeFlag(events.NodeAddedLatency)
	if err != nil {
		return errors.Wrap(err, "failed to get node flag (msec)")
	}

	duration, ok := durationFlag.(time.Duration)
	if !ok {
		return errors.New("NodeAddedLatency is not a duration")
	}

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context cancelled while running latency")
	case <-time.After(duration):
		// Do the actual latency
		return nil
	}
}

func (p *Provider) getPodLabelByName(ns string, name string) string {
	pod, ok := p.Pods.GetPodByName(ns, name)
	if !ok {
		return ""
	}
	return getPodLabelByPod(pod)
}

func getPodLabelByPod(pod *corev1.Pod) string {
	label, ok := pod.Labels[podconfigv1.PodConfigurationLabel]
	if !ok {
		return ""
	}
	return pod.Namespace + "/" + label
}

func (p *Provider) doesPodExceedLimit(ctx context.Context, ns string, name string, label string) (bool, error) {
	currentResources, err := p.getCurrentPodResources(label)
	if err != nil {
		return false, errors.Wrap(err, "unable to determine current pod resources")
	}

	limits := p.getPodResourceLimits(ns, name)
	podExceedsPodLimit := currentResources.cpu > limits.cpu || currentResources.memory > limits.memory || currentResources.ephemeralStorage > limits.ephemeralStorage

	// If the total amount of all pods resources exceed the resources on the node, just kill the current one
	// TODO implement k8s OOM handling (much more complicated)
	statsSummary, err := p.GetStatsSummary(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to retrieve stats summary")
	}

	nodeStats := statsSummary.Node
	totalLimitExceeded := *nodeStats.CPU.UsageNanoCores > uint64(p.Resources.CPU) ||
		*nodeStats.Memory.UsageBytes > uint64(p.Resources.Memory) ||
		*nodeStats.Fs.UsedBytes > uint64(p.Resources.EphemeralStorage)

	return podExceedsPodLimit || totalLimitExceeded, nil
}

func (p *Provider) getCurrentPodResources(label string) (resources, error) {
	podResourcesFlag, err := (*p.Store).GetPodFlag(label, events.PodResources)
	if err != nil {
		return resources{}, errors.Wrap(err, "failed to get pod resources flag while getting pod status")
	}

	podResources, ok := podResourcesFlag.(*stats.PodStats)
	if !ok {
		return resources{}, errors.Wrapf(err, "unable to convert '%v' to PodStats", podResourcesFlag)
	}

	usageCores := uint64(0)
	if podResources.CPU != nil && podResources.CPU.UsageNanoCores != nil {
		usageCores = *podResources.CPU.UsageNanoCores
	}

	usageMemory := uint64(0)
	if podResources.Memory != nil && podResources.Memory.UsageBytes != nil {
		usageMemory = *podResources.Memory.UsageBytes
	}

	usageEphemeralStorage := uint64(0)
	if podResources.EphemeralStorage != nil && podResources.EphemeralStorage.UsedBytes != nil {
		usageEphemeralStorage = *podResources.EphemeralStorage.UsedBytes
	}

	return resources{
		cpu:              usageCores,
		memory:           usageMemory,
		ephemeralStorage: usageEphemeralStorage,
	}, nil
}

func (p *Provider) getPodResourceLimits(ns string, name string) resources {
	pod, ok := p.Pods.GetPodByName(ns, name)
	if !ok {
		return resources{}
	}

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

func podStatusToPhase(status interface{}) corev1.PodPhase {
	switch status {
	case scenario.PodStatusPending:
		return corev1.PodPending
	case scenario.PodStatusUnset:
		fallthrough // act as a normal pod
	case scenario.PodStatusRunning:
		return corev1.PodRunning
	case scenario.PodStatusSucceeded:
		return corev1.PodSucceeded
	case scenario.PodStatusFailed:
		return corev1.PodFailed
	case scenario.PodStatusUnknown:
		return corev1.PodUnknown
	default:
		return corev1.PodUnknown
	}
}
