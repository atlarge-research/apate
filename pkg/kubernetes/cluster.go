// Package kubernetes provides an interface to manage a kubernetes cluster with the help of
// kind en kubernetes' client-go modules. Use the Builder to create a new cluster.
package kubernetes

import (
	"log"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

const (
	nodeDeletionGracePeriod int64 = 0
)

// Cluster object can be used to interact with kubernetes clusters.
// It abstracts away calls to the kubernetes client-go api.
type Cluster struct {
	clientSet *kubernetes.Clientset

	KubeConfig *kubeconfig.KubeConfig
}

// ClusterFromKubeConfig Creates a new KubernetesCluster from a location of a configuration file.
func ClusterFromKubeConfig(kubeConfig *kubeconfig.KubeConfig) (Cluster, error) {
	restconfig, err := kubeConfig.GetConfig()
	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to get rest config from Kubeconfig")
	}

	clientSet, err := kubernetes.NewForConfig(restconfig)

	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to create kubernetes cluster from config")
	}

	return Cluster{
		clientSet,
		kubeConfig,
	}, nil
}

// ManagedCluster creates a managed cluster from an unmanaged cluster.
// If you know the name and manager type of a cluster, you can make an unmanaged cluster become managed,
// and you are for example able to delete it.
func (c Cluster) ManagedCluster(name string, manager Manager) ManagedCluster {
	return ManagedCluster{
		c,
		manager,
		name,
	}
}

// GetNumberOfPods returns the number of pods in the cluster, or an error if it couldn't get these.
func (c Cluster) GetNumberOfPods(namespace string) (int, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return -1, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(pods.Items), nil
}

// RemoveNodesFromCluster removes all apatelet nodes with the given uuids from the cluster.
func (c Cluster) RemoveNodesFromCluster(uuids []uuid.UUID) error {
	gracePeriod := nodeDeletionGracePeriod

	lbl := make([]string, 0, len(uuids))

	for _, id := range uuids {
		lbl = append(lbl, id.String())
	}
	// Sort for determinism.
	sort.StringSlice(lbl).Sort()
	labelSelector := nodeconfigv1.NodeIDLabel + " in (" + strings.Join(lbl, ",") + ")"
	log.Printf("deleting %s\n", labelSelector)
	return errors.Wrap(c.clientSet.CoreV1().Nodes().DeleteCollection(&metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	}, metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         10000,
	}), "failed to remove nodes from cluster")
}

// RemoveNodeFromCluster removes a node with a given name from the cluster.
func (c Cluster) RemoveNodeFromCluster(nodename string) error {
	gracePeriod := nodeDeletionGracePeriod
	return errors.Wrap(c.clientSet.CoreV1().Nodes().Delete(nodename, &metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	}), "failed to remove node from cluster")
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
