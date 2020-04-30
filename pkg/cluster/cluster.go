// Package cluster provides an interface to manage a kubernetes cluster with the help of
// kind en kubernetes' client-go modules. Use the Builder to create a new cluster.
package cluster

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// A KubernetesCluster object can be used to interact with kubernetes clusters.
// It abstracts away calls to the kubernetes client-go api.
type KubernetesCluster struct {
	clientSet *kubernetes.Clientset

	KubeConfig KubeConfig
}

// KubernetesClusterFromLocation Creates a new KubernetesCluster from a location of a configuration file.
func KubernetesClusterFromLocation(kubeConfig KubeConfig) (KubernetesCluster, error) {
	restconfig, err := kubeConfig.GetConfig()
	if err != nil {
		return KubernetesCluster{}, err
	}

	clientSet, err := kubernetes.NewForConfig(restconfig)

	if err != nil {
		return KubernetesCluster{}, err
	}

	return KubernetesCluster{
		clientSet,
		kubeConfig,
	}, nil
}

// ManagedCluster creates a managed cluster from an unmanaged cluster.
// If you know the name and manager type of a cluster, you can make an unmanaged cluster become managed,
// and you are for example able to delete it.
func (c KubernetesCluster) ManagedCluster(name string, manager Manager) ManagedCluster {
	return ManagedCluster{
		c,
		manager,
		name,
	}
}

// GetNumberOfPods returns the number of pods in the cluster, or an error if it couldn't get these.
func (c KubernetesCluster) GetNumberOfPods(namespace string) (int, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(pods.Items), nil
}

// RemoveNodeFromCluster removes a node with a given name from the cluster.
func (c KubernetesCluster) RemoveNodeFromCluster(nodename string) error {
	return c.clientSet.CoreV1().Nodes().Delete(nodename, &metav1.DeleteOptions{})
}
