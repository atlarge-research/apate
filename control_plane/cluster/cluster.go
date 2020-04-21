// Package cluster provides an interface to manage a kubernetes cluster with the help of
// kind en kubernetes' client-go modules. Use the ClusterBuilder to create a new cluster.
package cluster

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// A Cluster object can be used to interact with kubernetes clusters.
// It abstracts away calls to the kubernetes client-go api.
type Cluster struct {
	name string
	clientset *kubernetes.Clientset
}

// Used to delete a cluster
func (c *Cluster) Delete() {
	DeleteCluster(c.name)
}

// Returns the number of pods in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfPods() (int, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(pods.Items),  nil
}
