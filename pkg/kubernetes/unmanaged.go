package kubernetes

import (
	"github.com/pkg/errors"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/env"
	"github.com/atlarge-research/apate/pkg/kubernetes/kubeconfig"
)

// UnmanagedClusterManager is a struct which implements ClusterManager by writing the given kubeconfig to a file.
type UnmanagedClusterManager struct{}

// GetKubeConfig writes the given kubeconfig and returns it as a struct.
func (u *UnmanagedClusterManager) GetKubeConfig() (*kubeconfig.KubeConfig, error) {
	kubeConfig, err := kubeconfig.FromBytes([]byte(env.ControlPlaneEnv().KubeConfig), env.ControlPlaneEnv().KubeConfigLocation, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kube config from env variable")
	}
	return kubeConfig, nil
}

// Shutdown currently does nothing.
func (u *UnmanagedClusterManager) Shutdown(cluster *Cluster) error {
	err := nodeconfigv1.UpdateInKubernetes(cluster.KubeConfig, true)
	if err != nil {
		return errors.Wrap(err, "failed to remove node CRD")
	}

	err = podconfigv1.UpdateInKubernetes(cluster.KubeConfig, true)
	if err != nil {
		return errors.Wrap(err, "failed to remove pod CRD")
	}

	return errors.Wrap(cluster.RemoveAllApateletsFromCluster(), "failed to remove all apatelets from cluster")
}
