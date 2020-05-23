package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	// TODO: Is this the proper way to check ctx?
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return p.createOrUpdate(ctx, pod, events.PodCreatePodResponse, events.NodeCreatePodResponse)
}

func (p *Provider) createOrUpdate(ctx context.Context, pod *corev1.Pod, pf events.PodEventFlag, nf events.NodeEventFlag) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Create or Update)")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, updateMap(p, pod)},
		getPodLabelByPod(pod),
		pf,
		nf,
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

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	return p.createOrUpdate(ctx, pod, events.PodUpdatePodResponse, events.NodeUpdatePodResponse)
}

func updateMap(p *Provider, pod *corev1.Pod) func() (interface{}, error) {
	return func() (interface{}, error) {
		p.pods.AddPod(*pod)
		return nil, nil
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	// TODO: Is this the proper way to check ctx?
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Delete Pod)")
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
	// TODO: Is this the proper way to check ctx?
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Get Pod)")
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

	return nil, fmt.Errorf("invalid pod: %v", pod)
}

func podStatusToPhase(status interface{}) corev1.PodPhase {
	switch status {
	case scenario.PodStatusPending:
		return corev1.PodPending
	case scenario.ResponseUnset:
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

// GetPodStatus retrieves the status of a pod by label.
func (p *Provider) GetPodStatus(ctx context.Context, ns string, name string) (*corev1.PodStatus, error) {
	// TODO: Is this the proper way to check ctx?
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Get Pod Status)")
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

	return nil, fmt.Errorf("invalid podstatus: %v", pod)
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	// TODO: Is this the proper way to check ctx?
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Get Pods)")
		log.Println(err)
		return nil, err
	}

	pod, err := nodeResponse(responseArgs{ctx, p, func() (interface{}, error) {
		return p.pods.GetAllPods(), nil
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

	return nil, fmt.Errorf("invalid pods: %v", pod)
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
		return errors.Wrap(ctx.Err(), "context cancelled while running latency")
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
