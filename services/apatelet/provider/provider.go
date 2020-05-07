// Package provider implements the virtual kubelet provider emulate to facilitate emulating pods.
package provider

import (
	"errors"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
	"math/rand"
	"sync"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"bytes"
	"context"
	"io"
	"io/ioutil"

	vkprov "github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// VKProvider implements the virtual-kubelet provider interface
type VKProvider struct {
	store     *store.Store
	Pods      map[types.UID]*corev1.Pod
	PodLock   sync.RWMutex
	resources *normalization.NodeResources
}

// CreateProvider returns the provider but with the vk type instead of our own.
func CreateProvider(resources *normalization.NodeResources, store *store.Store) vkprov.Provider {
	return &VKProvider{
		resources: resources,
		store:     store,
		Pods:      make(map[types.UID]*corev1.Pod),
	}
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *VKProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.runLatency(ctx); err != nil {
		return err
	}

	_, err := podAndNodeResponse(podNodeResponse{
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
func (p *VKProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
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

func updateMap(p *VKProvider, pod *corev1.Pod) func() (interface{}, error) {
	return func() (interface{}, error) {
		p.PodLock.Lock()
		p.Pods[pod.UID] = pod
		p.PodLock.Unlock()
		return nil, nil
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *VKProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	if err := p.runLatency(ctx); err != nil {
		return err
	}

	_, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.Lock()
			delete(p.Pods, pod.UID)
			p.PodLock.Unlock()
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
func (p *VKProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}

	pod, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.RLock()
			defer p.PodLock.RUnlock()
			for _, element := range p.Pods {
				if element.Namespace == namespace && element.Name == name {
					return element, nil
				}
			}
			return nil, throw.Exception("unable to find pod")
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

	return pod.(*corev1.Pod), err
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
func (p *VKProvider) GetPodStatus(ctx context.Context, namespace string, name string) (*corev1.PodStatus, error) {
	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}

	pod, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{ctx, p, func() (interface{}, error) {
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
					Phase:   corev1.PodRunning,
					Message: "Emulating pod successfully",
				}, nil
			}

			return &corev1.PodStatus{
				Phase:   podStatusToPhase(status),
				Message: "Emulating pod successfully",
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

	return pod.(*corev1.PodStatus), err
}

// GetPods retrieves a list of all pods running.
func (p *VKProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	if err := p.runLatency(ctx); err != nil {
		return nil, err
	}
	pod, err := nodeResponse(responseArgs{ctx, p, func() (interface{}, error) {
		var arr []*corev1.Pod

		for _, element := range p.Pods {
			arr = append(arr, element)
		}
		return arr, nil
	},
	},
		nodeResponseArgs{
			nodeResponseFlag:   events.NodeGetPodsResponse,
			nodePercentageFlag: events.NodeGetPodsResponsePercentage,
		},
	)

	return pod.([]*corev1.Pod), err
}

// GetContainerLogs retrieves the log of a specific container.
func (p *VKProvider) GetContainerLogs(context.Context, string, string, string, api.ContainerLogOpts) (io.ReadCloser, error) {
	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte("This container is emulated by Apate"))), nil
}

// RunInContainer retrieves the log of a specific container.
func (p *VKProvider) RunInContainer(context.Context, string, string, string, []string, api.AttachIO) error {
	// There is no actual process running in the containers, so we can't do anything.
	return nil
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *VKProvider) ConfigureNode(_ context.Context, v *corev1.Node) {
	var cpu resource.Quantity
	cpu.Set(p.resources.CPU)

	var mem resource.Quantity
	mem.Set(p.resources.Memory)

	var pods resource.Quantity
	pods.Set(p.resources.MaxPods)

	var storage resource.Quantity
	storage.Set(p.resources.Storage)

	var ephemeralStorage resource.Quantity
	ephemeralStorage.Set(p.resources.EphemeralStorage)

	v.Status.Capacity = corev1.ResourceList{
		corev1.ResourceCPU:              cpu,
		corev1.ResourceMemory:           mem,
		corev1.ResourcePods:             pods,
		corev1.ResourceStorage:          storage,
		corev1.ResourceEphemeralStorage: ephemeralStorage,
	}
}

func (p *VKProvider) runLatency(ctx context.Context) error {
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
		return errors.New("NodeAddedLatencyEnabled is not a bool")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}
