package provider

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
)

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	log.Println("Creating pod")

	if err := p.runLatency(ctx); err != nil {
		return err
	}

	find, exists, err := p.crdInformer.Find(pod.Labels["apate"])
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Found CRD %v", find)
	} else {
		log.Printf("No CRD found")
	}

	_, err = podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, updateMap(p, pod)},
		podResponseArgs: podResponseArgs{
			pod.Name,
			events.PodCreatePodResponse,
			events.PodCreatePodResponsePercentage,
		},
		nodeResponseArgs: nodeResponseArgs{
			events.NodeCreatePodResponse,
			events.NodeCreatePodResponsePercentage,
		},
	})

	return err
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	log.Println("Updating pod")

	if err := p.runLatency(ctx); err != nil {
		return err
	}

	_, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, updateMap(p, pod)},
		podResponseArgs: podResponseArgs{
			pod.Name,
			events.PodUpdatePodResponse,
			events.PodUpdatePodResponsePercentage,
		},
		nodeResponseArgs: nodeResponseArgs{
			events.NodeUpdatePodResponse,
			events.NodeUpdatePodResponsePercentage,
		},
	})

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
	log.Println("Delete pod")

	if err := p.runLatency(ctx); err != nil {
		return err
	}

	_, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, func() (interface{}, error) {
			p.pods.DeletePod(pod)
			return nil, nil
		}},
		podResponseArgs: podResponseArgs{
			pod.Name,
			events.PodDeletePodResponse,
			events.PodDeletePodResponsePercentage,
		},
		nodeResponseArgs: nodeResponseArgs{
			events.NodeDeletePodResponse,
			events.NodeDeletePodResponsePercentage,
		},
	})

	return err
}

// GetPod retrieves a pod by name.
func (p *Provider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	log.Println("Getting pod")

	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}

	pod, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, func() (interface{}, error) {
			return p.pods.GetPodByName(namespace, name), nil
		}},
		podResponseArgs: podResponseArgs{
			name,
			events.PodGetPodResponse,
			events.PodGetPodResponsePercentage,
		},
		nodeResponseArgs: nodeResponseArgs{
			events.NodeGetPodResponse,
			events.NodeGetPodResponsePercentage,
		},
	})

	if err != nil {
		return nil, err
	}

	return pod.(*corev1.Pod), nil
}

func podStatusToPhase(status interface{}) corev1.PodPhase {
	switch status {
	case scenario.PodStatus_POD_PENDING:
		return corev1.PodPending
	case scenario.PodStatus_POD_RUNNING:
		return corev1.PodRunning
	case scenario.PodStatus_POD_SUCCEEDED:
		return corev1.PodSucceeded
	case scenario.PodStatus_POD_FAILED:
		return corev1.PodFailed
	case scenario.PodStatus_POD_UNKNOWN:
		fallthrough
	default:
		return corev1.PodUnknown
	}
}

// GetPodStatus retrieves the status of a pod by name.
func (p *Provider) GetPodStatus(ctx context.Context, namespace string, name string) (*corev1.PodStatus, error) {
	log.Println("Getting pod status")

	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}

	pod, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx: ctx, provider: p, action: func() (interface{}, error) {
			status, err := (*p.store).GetPodFlag(name, events.PodUpdatePodStatus)
			if err != nil {
				return nil, throw.Exception(err.Error())
			}

			ipercent, err := (*p.store).GetPodFlag(name, events.PodUpdatePodStatusPercentage)
			if err != nil {
				return nil, throw.Exception(err.Error())
			}

			percent, ok := ipercent.(int32)
			if !ok {
				return nil, throw.Exception("cast error")
			}

			if percent < rand.Int31n(int32(100)) {
				return &corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionTrue,
							LastProbeTime:      metav1.Time{Time: time.Now()},
							LastTransitionTime: metav1.Time{Time: time.Now()},
							Message:            "Emulating pod...",
						},
					},
					Message: "Emulating pod successfully",
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
			}, nil
		}},
		podResponseArgs: podResponseArgs{
			name,
			events.PodGetPodStatusResponse,
			events.PodGetPodStatusResponsePercentage,
		},
		nodeResponseArgs: nodeResponseArgs{
			events.NodeGetPodStatusResponse,
			events.NodeGetPodStatusResponsePercentage,
		},
	})

	if err != nil {
		return nil, err
	}

	return pod.(*corev1.PodStatus), nil
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	log.Println("Getting pods")

	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}
	pod, _, err := nodeResponse(responseArgs{ctx, p, func() (interface{}, error) {
		return p.pods.GetAllPods(), nil
	},
	},
		nodeResponseArgs{
			nodeResponseFlag:   events.NodeGetPodsResponse,
			nodePercentageFlag: events.NodeGetPodsResponsePercentage,
		},
	)

	if err != nil {
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
		return err
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
		return err
	}

	ms, ok := ims.(int64)
	if !ok {
		return errors.New("NodeAddedLatencyMsec is not an int")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}
