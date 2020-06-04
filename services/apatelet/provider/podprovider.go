package provider

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"

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

	return p.createOrUpdate(ctx, pod, events.PodCreatePodResponse, true)
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context cancelled in UpdatePod")
	}

	return p.createOrUpdate(ctx, pod, events.PodUpdatePodResponse, false)
}

func (p *Provider) createOrUpdate(ctx context.Context, pod *corev1.Pod, pf events.PodEventFlag, updateStartTime bool) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Create or Update)")
		log.Println(err)
		return err
	}

	_, err := podResponse(
		responseArgs{ctx, p, updateMap(p, pod, updateStartTime)},
		pod,
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

	_, err := podResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			p.Pods.DeletePod(pod)
			return nil, nil
		}},
		pod,
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
	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "context cancelled in GetPod")
	}

	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency GetPod")
		log.Println(err)
		return nil, err
	}

	pod, ok := p.Pods.GetPodByName(namespace, name)
	if !ok {
		return nil, nil
	}

	_, err := podResponse(
		responseArgs{ctx, p, func() (interface{}, error) {
			return nil, nil
		}},
		pod,
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

	return pod.DeepCopy(), nil
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
