package kubernetes

import (
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
)

// ClusterAPI are API calls for a cluster connection
type ClusterAPI interface {
	GetNodes() (*corev1.NodeList, error)
	GetPods(namespace string) (*corev1.PodList, error)

	GetNumberOfNodes() (int, error)
	GetNumberOfReadyNodes() (int, error)
	GetNumberOfPods(namespace string) (int, error)

	RemoveAllApateletsFromCluster() error
	RemoveNodesFromCluster(uuids []uuid.UUID) error
	RemoveNodeFromCluster(nodename string) error
}

// GetNodes returns the number of nodes in the cluster, or an error if it couldn't get these.
func (c Cluster) GetNodes() (*corev1.NodeList, error) {
	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve nodes list from kubernetes")
	}

	return nodes, nil
}

// GetPods gets a list of pods from kubernetes using the specified namespace
func (c Cluster) GetPods(namespace string) (*corev1.PodList, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve pod list from kubernetes")
	}

	return pods, nil
}

const nodeDeletionGracePeriod int64 = 0

// GetNumberOfPods returns the number of pods in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfPods(namespace string) (int, error) {
	pods, err := c.GetPods(namespace)
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
	nodes, err := c.GetNodes()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	return len(nodes.Items), nil
}

// GetNumberOfReadyNodes returns the number of nodes in the cluster, or an error if it couldn't get these.
func (c *Cluster) GetNumberOfReadyNodes() (int, error) {
	nodes, err := c.GetNodes()

	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve pods list from kubernetes")
	}

	i := 0
	for _, node := range nodes.Items {
		for _, c := range node.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				i++
			}
		}
	}

	return i, nil
}
