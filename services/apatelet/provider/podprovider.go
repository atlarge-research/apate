package provider

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, updateMap(p, pod)},
		getPodLabelByPod(pod),
		events.PodCreatePodResponse,
		events.NodeCreatePodResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Create Pod)")
	if err != nil {
		log.Println(err)
	}

	return err
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, updateMap(p, pod)},
		getPodLabelByPod(pod),
		events.PodUpdatePodResponse,
		events.NodeUpdatePodResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Update Pod)")
	if err != nil {
		log.Println(err)
	}

	return err
}

func updateMap(p *Provider, pod *corev1.Pod) func() (interface{}, error) {
	return func() (interface{}, error) {
		p.pods.AddPod(*pod)
		return nil, nil
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			p.pods.DeletePod(pod)
			return nil, nil
		}},
		getPodLabelByPod(pod),
		events.PodDeletePodResponse,
		events.NodeDeletePodResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Delete Pod)")
	if err != nil {
		log.Println(err)
	}

	return err
}

// GetPod retrieves a pod by label.
func (p *Provider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return nil, err
	}

	label := p.getPodLabelByName(namespace, name)

	pod, err := podAndNodeResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			pod, _ := p.pods.GetPodByName(namespace, name)
			return pod, nil
		}},
		label,
		events.PodGetPodResponse,
		events.NodeGetPodResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pod)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pod.(*corev1.Pod), nil
}

func podStatusToPhase(status interface{}) corev1.PodPhase {
	switch status {
	case scenario.PodStatus_POD_STATUS_PENDING:
		return corev1.PodPending
	case scenario.PodStatus_POD_STATUS_UNSET:
		fallthrough // act as a normal pod
	case scenario.PodStatus_POD_STATUS_RUNNING:
		return corev1.PodRunning
	case scenario.PodStatus_POD_STATUS_SUCCEEDED:
		return corev1.PodSucceeded
	case scenario.PodStatus_POD_STATUS_FAILED:
		return corev1.PodFailed
	case scenario.PodStatus_POD_STATUS_UNKNOWN:
		fallthrough
	default:
		return corev1.PodUnknown
	}
}

// GetPodStatus retrieves the status of a pod by label.
func (p *Provider) GetPodStatus(ctx context.Context, ns string, name string) (*corev1.PodStatus, error) {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return nil, err
	}

	label := p.getPodLabelByName(ns, name)

	pod, err := podAndNodeResponse(responseArgs{ctx: ctx, provider: p, action: func() (interface{}, error) {
		status, err := (*p.store).GetPodFlag(label, events.PodStatus)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get pod flag while getting pod status")
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
		}, nil
	}},
		label,
		events.PodGetPodStatusResponse,
		events.NodeGetPodStatusResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pod Status)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pod.(*corev1.PodStatus), nil
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency")
		log.Println(err)
		return nil, err
	}

	pod, err := nodeResponse(responseArgs{ctx, p, func() (interface{}, error) {
		return p.pods.GetAllPods(), nil
	}},
		events.NodeGetPodsResponse,
	)

	err = errors.Wrap(err, "failed to execute pod and node response (Get Pods)")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pod.([]*corev1.Pod), nil
}

// GetContainerLogs retrieves the log of a specific container.
func (p *Provider) GetContainerLogs(context.Context, string, string, string, api.ContainerLogOpts) (io.ReadCloser, error) {
	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte("This container is emulated by Apate\n"))), nil
}

// RunInContainer retrieves the log of a specific container.
func (p *Provider) RunInContainer(context.Context, string, string, string, []string, api.AttachIO) error {
	// There is no actual process running in the containers, so we can't do anything.
	return nil
}

func (p *Provider) runLatency(ctx context.Context) error {
	val, err := (*p.store).GetNodeFlag(events.NodeAddedLatencyEnabled)
	if err != nil {
		return errors.Wrap(err, "failed to get node flag (enabled)")
	}

	y, ok := val.(bool)
	if !ok {
		return errors.New("NodeAddedLatencyEnabled is not a bool")
	}
	if !y {
		return nil
	}

	ims, err := (*p.store).GetNodeFlag(events.NodeAddedLatencyMsec)
	if err != nil {
		return errors.Wrap(err, "failed to get node flag (msec)")
	}

	ms, ok := ims.(int64)
	if !ok {
		return errors.New("NodeAddedLatencyMsec is not an int")
	}

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context failed while running latency")
	default:
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}

func (p *Provider) getPodLabelByName(ns string, name string) string {
	pod, ok := p.pods.GetPodByName(ns, name)
	if !ok {
		return ""
	}
	return getPodLabelByPod(pod)
}

func getPodLabelByPod(pod *corev1.Pod) string {
	return pod.Namespace + "/" + pod.Labels["apate"]
}
