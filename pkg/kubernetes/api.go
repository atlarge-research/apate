package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return errors.Wrap(c.clientSet.CoreV1().Nodes().Delete(nodename, &metav1.DeleteOptions{}), "failed to remove node from cluster")
}

// RemoveAllApateletsFromCluster removes all apatelets from the Kubernetes cluster.
func (c *Cluster) RemoveAllApateletsFromCluster() error {
	apatelets, err := c.GetAllApatelets()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve all apatelets")
	}

	for _, node := range apatelets {
		err := c.RemoveNodeFromCluster(node.Name)
		if err != nil {
			return errors.Wrapf(err, "failed to remove node from cluster with name %v", node.Name)
		}
	}

	return nil
}

// GetAllApatelets returns all apatelets from the Kubernetes cluster.
func (c *Cluster) GetAllApatelets() ([]*corev1.Node, error) {
	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	apatelets := make([]*corev1.Node, 0)
	for i, node := range nodes.Items {
		if strings.HasPrefix(node.Name, "apatelet-") {
			apatelets = append(apatelets, &nodes.Items[i])
		}
	}
	return apatelets, nil
}

// GetNumberOfPendingPods will return the number of pods in the pending state.
func (c *Cluster) GetNumberOfPendingPods(namespace string) (int, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return -1, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	cnt := 0
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodPending {
			cnt++
		}
	}

	return cnt, nil
}

// GetNumberOfNodes returns the number of nodes in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfNodes() (int, error) {
	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(nodes.Items), nil
}
