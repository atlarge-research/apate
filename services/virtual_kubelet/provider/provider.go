package provider

import (
	"bytes"
	"context"
	"errors"
	vkprov "github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// VKProvider implements the virtual-kubelet provider interface
type VKProvider struct {
	Pods map[types.UID]*corev1.Pod
}

// Returns the provider but with the vk type instead of our own.
func CreateProvider() vkprov.Provider {
	return &VKProvider{}
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *VKProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	p.Pods[pod.UID] = pod
	return nil
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *VKProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	p.Pods[pod.UID] = pod
	return nil
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *VKProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	delete(p.Pods, pod.UID)
	return nil
}

// GetPod retrieves a pod by name.
func (p *VKProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {

	// TODO: think about better structure for p.Pods
	for _, element := range p.Pods {
		if element.Namespace == namespace && element.Name == name {
			return element, nil
		}
	}

	return nil, errors.New("unable to find pod")
}

// GetPodStatus retrieves the status of a pod by name.
func (p *VKProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	return &corev1.PodStatus{}, nil
}

// GetPods retrieves a list of all pods running.
func (p *VKProvider) GetPods(context.Context) ([]*corev1.Pod, error) {

	// TODO: Improve
	var arr []*corev1.Pod

	for _, element := range p.Pods {
		arr = append(arr, element)
	}

	return arr, nil
}

// GetContainerLogs retrieves the log of a specific container.
func (p *VKProvider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte(""))), nil
}

// RunInContainer retrieves the log of a specific container.
func (p *VKProvider) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach api.AttachIO) error {
	// There is no actual process running in the containers, so we can't do anything.
	return nil
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *VKProvider) ConfigureNode(ctx context.Context, v *corev1.Node) {

	var cpu resource.Quantity
	cpu.Set(1000)

	var mem resource.Quantity
	mem.Set(1000)

	v.Status.Capacity = corev1.ResourceList{
		"cpu":    cpu,
		"memory": mem,
		"pods":   resource.MustParse("1000"),
	}
}
