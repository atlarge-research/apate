package kubernetes

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNumberOfPods returns the number of pods in the cluster, or an error if it couldn't get these.
func (c Cluster) GetNumberOfPods(namespace string) (int, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return -1, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(pods.Items), nil
}

// RemoveNodeFromCluster removes a node with a given name from the cluster.
func (c Cluster) RemoveNodeFromCluster(nodename string) error {
	return errors.Wrap(c.clientSet.CoreV1().Nodes().Delete(nodename, &metav1.DeleteOptions{}), "failed to remove node from cluster")
}

// GetNumberOfPendingPods will return the number of pods in the pending state.
func (c Cluster) GetNumberOfPendingPods(namespace string) (int, error) {
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
func (c Cluster) GetNumberOfNodes() (int, error) {
	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(nodes.Items), nil
}
