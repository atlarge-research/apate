// Package provider implements the virtual kubelet provider emulate to facilitate emulating pods.
package provider

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"sync"

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
	}
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *VKProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {

	_, err := magicPodAndNode(magicPodNodeArgs{
		magicArgs: magicArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.Lock()
			p.Pods[pod.UID] = pod
			p.PodLock.Unlock()
			return nil, nil
		}},
		magicPodArgs: magicPodArgs{
			pod.Name,
			events.PodCreatePodResponse,
			events.PodCreatePodResponsePercentage,
		},
		magicNodeArgs: magicNodeArgs{
			events.NodeCreatePodResponse,
			events.NodeCreatePodResponsePercentage,
		},
	})

	return err
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *VKProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := magicPodAndNode(magicPodNodeArgs{
		magicArgs: magicArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.Lock()
			p.Pods[pod.UID] = pod
			p.PodLock.Unlock()
			return nil, nil
		}},
		magicPodArgs: magicPodArgs{
			pod.Name,
			events.PodUpdatePodResponse,
			events.PodUpdatePodResponsePercentage,
		},
		magicNodeArgs: magicNodeArgs{
			events.NodeUpdatePodResponse,
			events.NodeUpdatePodResponsePercentage,
		},
	})

	return err
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *VKProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := magicPodAndNode(magicPodNodeArgs{
		magicArgs: magicArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.Lock()
			delete(p.Pods, pod.UID)
			p.PodLock.Unlock()
			return nil, nil
		}},
		magicPodArgs: magicPodArgs{
			pod.Name,
			events.PodDeletePodResponse,
			events.PodDeletePodResponsePercentage,
		},
		magicNodeArgs: magicNodeArgs{
			events.NodeDeletePodResponse,
			events.NodeDeletePodResponsePercentage,
		},
	})

	return err
}

// GetPod retrieves a pod by name.
func (p *VKProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod, err := magicPodAndNode(magicPodNodeArgs{
		magicArgs: magicArgs{ctx, p, func() (interface{}, error) {
			p.PodLock.RLock()
			defer p.PodLock.RUnlock()
			for _, element := range p.Pods {
				if element.Namespace == namespace && element.Name == name {
					return element, nil
				}
			}
			return nil, wError("unable to find pod")
		}},
		magicPodArgs: magicPodArgs{
			name,
			events.PodGetPodResponse,
			events.PodGetPodResponsePercentage,
		},
		magicNodeArgs: magicNodeArgs{
			events.NodeGetPodResponse,
			events.NodeGetPodResponsePercentage,
		},
	})

	return pod.(*corev1.Pod), err
}

// GetPodStatus retrieves the status of a pod by name.
func (p *VKProvider) GetPodStatus(ctx context.Context, namespace string, name string) (*corev1.PodStatus, error) {
	runningStatus := corev1.PodStatus{
		Phase:   corev1.PodRunning,
		Message: "Emulating pod successfully",
	}

	pod, err := magicPodAndNode(magicPodNodeArgs{
		magicArgs: magicArgs{ctx, p, func() (interface{}, error) {
			return &runningStatus, nil
		}},
		magicPodArgs: magicPodArgs{
			name,
			events.PodGetPodStatusResponse,
			events.PodGetPodStatusResponsePercentage,
		},
		magicNodeArgs: magicNodeArgs{
			events.NodeGetPodStatusResponse,
			events.NodeGetPodStatusResponsePercentage,
		},
	})

	return pod.(*corev1.PodStatus), err
}

// GetPods retrieves a list of all pods running.
func (p *VKProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	pod, err := magicNode(magicArgs{ctx, p, func() (interface{}, error) {
		var arr []*corev1.Pod

		for _, element := range p.Pods {
			arr = append(arr, element)
		}
		return arr, nil
	},
	},
		magicNodeArgs{
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
