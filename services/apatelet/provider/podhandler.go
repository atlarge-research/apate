package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net"
	"strings"
)

type Provider struct {
	Pods      map[types.UID]*corev1.Pod
	resources *normalization.NodeResources
	cfg       provider.InitConfig
	nodeInfo  cluster.NodeInfo
}

// CreateProvider returns the provider but with the vk type instead of our own.
func NewProvider(resources *normalization.NodeResources, cfg provider.InitConfig, nodeInfo cluster.NodeInfo) provider.Provider {
	return &Provider{
		resources: resources,
		cfg:       cfg,
		nodeInfo:  nodeInfo,
	}
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(_ context.Context, pod *corev1.Pod) error {
	fmt.Println("CreatePod called")
	p.Pods[pod.UID] = pod
	return nil
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(_ context.Context, pod *corev1.Pod) error {
	fmt.Println("UpdatePod called")
	p.Pods[pod.UID] = pod
	return nil
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(_ context.Context, pod *corev1.Pod) error {
	fmt.Println("DeletePod called")
	delete(p.Pods, pod.UID)
	return nil
}

// GetPod retrieves a pod by name.
func (p *Provider) GetPod(_ context.Context, namespace, name string) (*corev1.Pod, error) {
	fmt.Println("GetPod called")
	// TODO: think about better structure for p.Pods
	for _, element := range p.Pods {
		if element.Namespace == namespace && element.Name == name {
			return element, nil
		}
	}

	return nil, errors.New("unable to find pod")
}

// GetPodStatus retrieves the status of a pod by name.
func (p *Provider) GetPodStatus(context.Context, string, string) (*corev1.PodStatus, error) {
	fmt.Println("GetPodStatus called")
	return &corev1.PodStatus{}, nil
}

// GetPods retrieves a list of all pods running.
func (p *Provider) GetPods(context.Context) ([]*corev1.Pod, error) {
	fmt.Println("GetPods called")
	// TODO: Improve
	var arr []*corev1.Pod

	for _, element := range p.Pods {
		arr = append(arr, element)
	}

	return arr, nil
}

// GetContainerLogs retrieves the log of a specific container.
func (p *Provider) GetContainerLogs(context.Context, string, string, string, api.ContainerLogOpts) (io.ReadCloser, error) {
	fmt.Println("GetContainerLogs called")
	// We return empty string as the emulated containers don't have a log.
	return ioutil.NopCloser(bytes.NewReader([]byte("This container is emulated by Apate"))), nil
}

// RunInContainer retrieves the log of a specific container.
func (p *Provider) RunInContainer(context.Context, string, string, string, []string, api.AttachIO) error {
	// There is no actual process running in the containers, so we can't do anything.
	fmt.Println("RunInContainer called")
	return nil
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *Provider) ConfigureNode(_ context.Context, node *corev1.Node) {
	node.Spec = p.spec()
	node.ObjectMeta = p.objectMeta()
	node.Status = p.nodeStatus()

	// TODO: https://github.com/virtual-kubelet/azure-aci/blob/4ad70d2ccfbc90b24e48bf7fd4e76f293dc3cca5/provider/aci.go#L1172
}

func (p *Provider) nodeStatus() corev1.NodeStatus {
	return corev1.NodeStatus{
		NodeInfo: corev1.NodeSystemInfo{
			Architecture:   "amd64",
			KubeletVersion: p.nodeInfo.Version,
		},
		DaemonEndpoints: p.nodeDaemonEndpoints(),
		Addresses:       p.addresses(),
		Capacity:        p.capacity(),
	}
}

func (p *Provider) objectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: p.nodeInfo.Name,
		Labels: map[string]string{
			"type":                   p.nodeInfo.NodeType,
			"kubernetes.io/role":     p.nodeInfo.Role,
			"kubernetes.io/hostname": p.nodeInfo.Name,
		},
	}
}

func (p *Provider) spec() corev1.NodeSpec {
	taints := make([]corev1.Taint, 0)
	return corev1.NodeSpec{
		Taints: taints,
	}
}

func (p *Provider) addresses() []corev1.NodeAddress {
	return []corev1.NodeAddress{
		{
			Type:    "InternalIP",
			Address: getExternalAddress(),
		},
		{
			Type:    "ExternalIP",
			Address: getExternalAddress(),
		},
	}
}

// TODO: Make better, rename, basically everything except code
func getExternalAddress() string {
	// Check for external IP override
	override := container.RetrieveFromEnvironment(container.ControlPlaneExternalIP, container.ControlPlaneExternalIPDefault)
	if override != container.ControlPlaneExternalIPDefault {
		return override
	}

	// Check for IP in interface addresses
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return "error"
	}

	// Get first 172.17.0.0/16 address, if any
	for _, address := range addresses {
		if strings.Contains(address.String(), container.DockerAddressPrefix) {
			ip := strings.Split(address.String(), "/")[0]

			return ip
		}
	}

	// Default to localhost
	return "localhost"
}

func (p *Provider) nodeDaemonEndpoints() corev1.NodeDaemonEndpoints {
	return corev1.NodeDaemonEndpoints{
		KubeletEndpoint: corev1.DaemonEndpoint{
			Port: p.cfg.DaemonPort,
		},
	}
}

func (p *Provider) capacity() corev1.ResourceList {
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

	return corev1.ResourceList{
		corev1.ResourceCPU:              cpu,
		corev1.ResourceMemory:           mem,
		corev1.ResourcePods:             pods,
		corev1.ResourceStorage:          storage,
		corev1.ResourceEphemeralStorage: ephemeralStorage,
	}
}
