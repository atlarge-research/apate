// Package kubernetes provides an interface to manage a kubernetes cluster with the help of
// kind en kubernetes' client-go modules. Use the Builder to create a new cluster.
package kubernetes

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// Cluster object can be used to interact with kubernetes clusters.
// It abstracts away calls to the kubernetes client-go api.
type Cluster struct {
	clientSet  *kubernetes.Clientset
	KubeConfig *kubeconfig.KubeConfig
	manager    *ClusterManager
}

// Shutdown performs clean up on the cluster before shutdown.
func (c *Cluster) Shutdown() error {
	return errors.Wrap((*c.manager).Shutdown(c), "shutting down cluster failed")
}

// ClusterManager contains utilities to handle connecting / creating a cluster and shutting it down if necessary.
type ClusterManager interface {
	// GetKubeConfig returns the KubeConfig by this cluster. This is also the place where one can create the cluster if necessary.
	GetKubeConfig() (*kubeconfig.KubeConfig, error)

	// Shutdown is called when the program exits. This can be used to clean up a cluster if necessary.
	Shutdown(*Cluster) error
}

// ClusterManagerHandler contains utilities to create cluster objects based on a certain manager.
type ClusterManagerHandler struct {
	clusterManager ClusterManager
}

// NewClusterManagerHandler creates a new ClusterManagerHandler, depending on what environment variables have been set
func NewClusterManagerHandler() ClusterManagerHandler {
	if len(env.ControlPlaneEnv().KubeConfig) == 0 {
		return ClusterManagerHandler{clusterManager: &KinDClusterManager{}}
	}
	return ClusterManagerHandler{clusterManager: &UnmanagedClusterManager{}}
}

// NewCluster creates a new cluster object by calling GetKubeConfig on its manager
func (cmh *ClusterManagerHandler) NewCluster() (*Cluster, error) {
	if _, err := os.Stat(env.ControlPlaneEnv().KubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(env.ControlPlaneEnv().KubeConfigLocation), os.ModePerm); err != nil {
			return nil, errors.Wrapf(err, "failed to create directory for kubeconfig (%v)", path.Dir(env.ControlPlaneEnv().KubeConfigLocation))
		}
	}

	config, err := cmh.clusterManager.GetKubeConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kube config")
	}

	res, err := cmh.NewClusterFromKubeConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kind cluster from Kubeconfig")
	}

	return res, nil
}

// NewClusterFromKubeConfig creates a new KubernetesCluster from a location of a kubeconfig.
func (cmh *ClusterManagerHandler) NewClusterFromKubeConfig(kubeConfig *kubeconfig.KubeConfig) (*Cluster, error) {
	restconfig, err := kubeConfig.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rest config from Kubeconfig")
	}

	clientSet, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes cluster from config")
	}

	return &Cluster{
		clientSet,
		kubeConfig,
		&cmh.clusterManager,
	}, nil
}
