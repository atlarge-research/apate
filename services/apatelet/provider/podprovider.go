package provider

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/pkg/errors"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in CreatePod")
	}

	return p.createOrUpdate(ctx, pod, events.PodCreatePodResponse, events.NodeCreatePodResponse, true)
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in UpdatePod")
	}

	return p.createOrUpdate(ctx, pod, events.PodUpdatePodResponse, events.NodeUpdatePodResponse, false)
}

func (p *Provider) createOrUpdate(ctx context.Context, pod *corev1.Pod, pf events.PodEventFlag, nf events.NodeEventFlag, updateStartTime bool) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Create or Update)")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, updateMap(p, pod, updateStartTime)},
		getPodLabelByPod(pod),
		pod,
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

func updateMap(p *Provider, pod *corev1.Pod, updateStartTime bool) func() (interface{}, error) {
	return func() (interface{}, error) {
		if updateStartTime {
			now := metav1.Now()
			pod.Status.StartTime = &now
		}
		p.Pods.AddPod(pod)
		return nil, nil
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in DeletePod")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency DeletePod")
		log.Println(err)
		return err
	}

	_, err := podAndNodeResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			p.Pods.DeletePod(pod)
			return nil, nil
		}},
		getPodLabelByPod(pod),
		pod,
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
	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "context cancelled in GetPod")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency GetPod")
		log.Println(err)
		return nil, err
	}

	pod, err := p.getPodByName(namespace, name)
	if err != nil {
		return nil, errors.Wrap(err, "pod not found")
	}
	label := getPodLabelByPod(pod)

	_, err = podAndNodeResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			return nil, nil
		}},
		label,
		pod,
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

	return pod, nil
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
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
func (p *Provider) GetContainerLogs(context.Context, string, string, string, api.ContainerLogOpts) (io.ReadCloser, error) {
	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte("This container is emulated by Apate\n"))), nil
}

// RunInContainer runs a command in a specific container.
func (p *Provider) RunInContainer(context.Context, string, string, string, []string, api.AttachIO) error {
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

func (p *Provider) getPodByName(ns string, name string) (*corev1.Pod, error) {
	pod, ok := p.Pods.GetPodByName(ns, name)
	if !ok {
		return nil, errors.New("f")
	}
	return pod, nil
}

func (p *Provider) getPodLabelByName(ns string, name string) string {
	pod, err := p.getPodByName(ns, name)
	if err != nil {
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
