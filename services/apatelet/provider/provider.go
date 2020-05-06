// Package provider implements the virtual kubelet provider emulate to facilitate emulating pods.
package provider

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"math/rand"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"

	vkprov "github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// VKProvider implements the virtual-kubelet provider interface
type VKProvider struct {
	store 	  *store.Store
	Pods      map[types.UID]*corev1.Pod
	resources *normalization.NodeResources
}

// CreateProvider returns the provider but with the vk type instead of our own.
func CreateProvider(resources *normalization.NodeResources, store * store.Store) vkprov.Provider {
	return &VKProvider{
		resources: resources,
		store: store,
	}
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *VKProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	iflag, err := (*p.store).GetPodFlag(pod.Name, events.PodCreatePodResponse)
	if err != nil {
		return err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return errors.New("invalid flag type")
	}
	
	iflagp, err := (*p.store).GetPodFlag(pod.Name, events.PodCreatePodResponsePercentage)
	if err != nil {
		return err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		p.Pods[pod.UID] = pod
		return nil
	}

	switch flag {
	case scenario.Response_NORMAL:
		p.Pods[pod.UID] = pod
		return nil
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil
	case scenario.Response_ERROR:
		return errors.New("expected error")
	default:
		return errors.New("invalid response")
	}
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *VKProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	iflag, err := (*p.store).GetPodFlag(pod.Name, events.PodUpdatePodResponse)
	if err != nil {
		return err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return errors.New("invalid flag type")
	}

	iflagp, err := (*p.store).GetPodFlag(pod.Name, events.PodUpdatePodResponsePercentage)
	if err != nil {
		return err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		p.Pods[pod.UID] = pod
		return nil
	}

	switch flag {
	case scenario.Response_NORMAL:
		p.Pods[pod.UID] = pod
		return nil
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil
	case scenario.Response_ERROR:
		return errors.New("expected error")
	default:
		return errors.New("invalid response")
	}
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *VKProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	iflag, err := (*p.store).GetPodFlag(pod.Name, events.PodDeletePodResponse)
	if err != nil {
		return err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return errors.New("invalid flag type")
	}

	iflagp, err := (*p.store).GetPodFlag(pod.Name, events.PodDeletePodResponsePercentage)
	if err != nil {
		return err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		delete(p.Pods, pod.UID)
		return nil
	}

	switch flag {
	case scenario.Response_NORMAL:
		delete(p.Pods, pod.UID)
		return nil
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil
	case scenario.Response_ERROR:
		return errors.New("expected error")
	default:
		return errors.New("invalid response")
	}
}

// GetPod retrieves a pod by name.
func (p *VKProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {

	iflag, err := (*p.store).GetPodFlag(name, events.PodDeletePodResponse)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.New("invalid flag type")
	}

	iflagp, err := (*p.store).GetPodFlag(name, events.PodDeletePodResponsePercentage)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		return getPod(p, namespace, name)
	}

	switch flag {
	case scenario.Response_NORMAL:
		// TODO: think about better structure for p.Pods
		return getPod(p, namespace, name)
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, errors.New("expected error")
	default:
		return nil, errors.New("invalid response")
	}
}

func getPod(p *VKProvider, namespace string, name string) (*corev1.Pod, error) {
	for _, element := range p.Pods {
		if element.Namespace == namespace && element.Name == name {
			return element, nil
		}
	}
	return nil, errors.New("unable to find pod")
}

// GetPodStatus retrieves the status of a pod by name.
func (p *VKProvider) GetPodStatus(ctx context.Context, namespace string, name string) (*corev1.PodStatus, error) {
	runningStatus := corev1.PodStatus{
		Phase:                 corev1.PodRunning,
		Message:               "Emulating pod successfully",
	}
	
	iflag, err := (*p.store).GetPodFlag(name, events.PodGetPodStatusResponse)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.New("invalid flag type")
	}

	iflagp, err := (*p.store).GetPodFlag(name, events.PodGetPodStatusResponsePercentage)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		return &runningStatus, nil
	}

	switch flag {
	case scenario.Response_NORMAL:
		return &runningStatus, nil
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, errors.New("expected error")
	default:
		return nil, errors.New("invalid response")
	}

}

// GetPods retrieves a list of all pods running.
func (p *VKProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	iflag, err := (*p.store).GetNodeFlag(events.NodeGetPodsResponse)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.New("invalid flag type")
	}

	iflagp, err := (*p.store).GetNodeFlag(events.NodeGetPodsResponsePercentage)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, errors.New("invalid percentage type")
	}

	if flagp < rand.Int31n(int32(100)) {
		return getPods(p), nil
	}

	switch flag {
	case scenario.Response_NORMAL:
		return getPods(p), nil
	case scenario.Response_TIMEOUT:
		<-ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, errors.New("expected error")
	default:
		return nil, errors.New("invalid response")
	}
}

// TODO: Improve?
func getPods(p *VKProvider) []*corev1.Pod {
	var arr []*corev1.Pod

	for _, element := range p.Pods {
		arr = append(arr, element)
	}
	return arr
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
