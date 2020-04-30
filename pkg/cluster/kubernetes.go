package cluster

import (
	"context"

	"github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	nodeType string
	role     string
	name     string
	version  string
}

// NewNode create a new NodeInfo
func NewNode(nodeType string, role string, name string, version string) *NodeInfo {
	return &NodeInfo{
		nodeType: nodeType,
		role:     role,
		name:     name,
		version:  version,
	}
}

// GetConfigForContext returns a kubernetes client configuration for the context given.
func GetConfigForContext(context string, kubeConfigLocation string) (*rest.Config, error) {
	// Create a default config rules struct
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	rules.ExplicitPath = kubeConfigLocation

	// Override with defaults (this call might not be necessary since the defaults are already set above?)
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	// But set the context to our own context while overriding
	if context != "" {
		overrides.CurrentContext = context
	}

	// You could override server ip/host here with this:
	// overrides.ClusterInfo.Server

	// Now create the actual configuration
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfigForContext returns a kubernetes client configuration from just the path to the kubeconfig.
func GetConfigForLocation(kubeConfigLocation string) (*rest.Config, error) {
	return GetConfigForContext("", kubeConfigLocation)
}

// CreateKubernetesNode creates a kubernetes api object representing a node
func CreateKubernetesNode(ctx context.Context, node NodeInfo, provider provider.Provider) *corev1.Node {
	taints := make([]corev1.Taint, 0)

	kubernetesNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: node.name,
			Labels: map[string]string{
				"type":                   node.nodeType,
				"kubernetes.io/role":     node.role,
				"kubernetes.io/hostname": node.name,
			},
		},
		Spec: corev1.NodeSpec{
			Taints: taints,
		},
		Status: corev1.NodeStatus{
			NodeInfo: corev1.NodeSystemInfo{
				Architecture:   "amd64",
				KubeletVersion: node.version,
			},
		},
	}

	provider.ConfigureNode(ctx, kubernetesNode)
	return kubernetesNode
}
