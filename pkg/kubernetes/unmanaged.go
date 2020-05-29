package kubernetes

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/pkg/errors"
)

type UnmanagedClusterManager struct {}

func (u UnmanagedClusterManager) GetKubeConfig() (*kubeconfig.KubeConfig, error) {
	kubeConfig, err := kubeconfig.FromBytes([]byte(env.ControlPlaneEnv().KubeConfig), env.ControlPlaneEnv().KubeConfigLocation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kube config from env variable")
	}
	return kubeConfig, nil
}

func (u UnmanagedClusterManager) Shutdown() error {
	return nil // noop, we want the cluster to still live afterwards
}
