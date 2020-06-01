package kubernetes

import (
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
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
	return errors.Wrap(cluster.RemoveAllApateletsFromCluster(), "failed to remove all apatelets from cluster")
}
