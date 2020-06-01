package kubernetes

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/pod"

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
	restConfig, err := cluster.KubeConfig.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get restconfig from kubeconfig")
	}

	err = deletePodConfigurations(restConfig)
	if err != nil {
		return errors.Wrap(err, "failed to remove all pod configurations")
	}

	err = deleteNodeConfigurations(restConfig)
	if err != nil {
		return errors.Wrap(err, "failed to remove all node configurations")
	}

	return errors.Wrap(cluster.RemoveAllApateletsFromCluster(), "failed to remove all apatelets from cluster")
}

func deletePodConfigurations(restConfig *rest.Config) error {
	podCRDClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		return errors.Wrap(err, "failed to get podclient from rest config for pod informer")
	}

	err = podCRDClient.Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete pod configurations")
	}
	return nil
}

func deleteNodeConfigurations(restConfig *rest.Config) error {
	nodeCRDClient, err := node.NewForConfig(restConfig, "default")
	if err != nil {
		return errors.Wrap(err, "failed to get nodeclient from rest config for pod informer")
	}
	err = nodeCRDClient.Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete node configurations")
	}
	return nil
}
