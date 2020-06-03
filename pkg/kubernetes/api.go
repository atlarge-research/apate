package kubernetes

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
)

// GetNumberOfPods returns the number of pods in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfPods(namespace string) (int, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return -1, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(pods.Items), nil
}

// RemoveNodeFromCluster removes a node with a given name from the cluster.
func (c *Cluster) RemoveNodeFromCluster(nodename string) error {
	return errors.Wrap(c.
		clientSet.
		CoreV1().
		Nodes().
		Delete(nodename, &metav1.DeleteOptions{}),
		"failed to remove node from cluster")
}

// RemoveAllApateletsFromCluster removes all apatelets from the Kubernetes cluster.
func (c *Cluster) RemoveAllApateletsFromCluster() error {
	return errors.Wrap(c.
		clientSet.
		CoreV1().
		Nodes().
		DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: nodeconfigv1.EmulatedLabel + "=" + nodeconfigv1.EmulatedLabelValue,
		}),
		"failed to remove all apatelets from cluster")
}

// GetNumberOfNodes returns the number of nodes in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfNodes() (int, error) {
	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(nodes.Items), nil
}
