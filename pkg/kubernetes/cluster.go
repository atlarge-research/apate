// Package kubernetes provides an interface to manage a kubernetes cluster with the help of
// kind en kubernetes' client-go modules. Use the Builder to create a new cluster.
package kubernetes

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"os"
	"path"
)

// Cluster object can be used to interact with kubernetes clusters.
// It abstracts away calls to the kubernetes client-go api.
type Cluster struct {
	clientSet  *kubernetes.Clientset
	KubeConfig *kubeconfig.KubeConfig
	manager    *ClusterManager
}

// Shutdown does iets
func (c Cluster) Shutdown() error {
	return errors.Wrap((*c.manager).Shutdown(), "error shutting down cluster")
}

// ClusterManager does something
type ClusterManager interface {
	GetKubeConfig() (*kubeconfig.KubeConfig, error)
	Shutdown() error
}

// ClusterManagerHandler does something else
type ClusterManagerHandler struct {
	clusterManager ClusterManager
}

// NewClusterManagerHandler does yet another thing
func NewClusterManagerHandler() ClusterManagerHandler {
	if len(env.ControlPlaneEnv().KubeConfig) == 0 {
		return ClusterManagerHandler{KinDClusterManager{}}
	} else {
		return ClusterManagerHandler{UnmanagedClusterManager{}}
	}
}

// NewCluster does ite
func (cmh ClusterManagerHandler) NewCluster() (Cluster, error) {
	if _, err := os.Stat(env.ControlPlaneEnv().KubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(env.ControlPlaneEnv().KubeConfigLocation), os.ModePerm); err != nil {
			return Cluster{}, errors.Wrapf(err, "failed to create directory for kubeconfig (%v)", path.Dir(env.ControlPlaneEnv().KubeConfigLocation))
		}
	}

	config, err := cmh.clusterManager.GetKubeConfig()
	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to get kube config")
	}

	res, err := cmh.NewClusterFromKubeConfig(config)
	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to create kind cluster from Kubeconfig")
	}

	return res, nil
}

// NewClusterFromKubeConfig Creates a new KubernetesCluster from a location of a configuration file.
func (cmh ClusterManagerHandler) NewClusterFromKubeConfig(kubeConfig *kubeconfig.KubeConfig) (Cluster, error) {
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
		&cmh.clusterManager,
	}, nil
}
