package kubernetes

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kind"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// Builder allows for the creation of KubernetesClusters.
type Builder struct {
	Name               string
	Manager            Manager
	KubeConfigLocation string
	ManagerConfigPath  string
}

// New is used to create a new Builder with all fields empty.
func New() Builder {
	return Builder{}
}

// Default is used to create a new Builder with all fields set to default values.
func Default() (c Builder) {
	c.Name = env.ControlPlaneEnv().KinDClusterName
	c.Manager = &kind.KinD{}
	c.KubeConfigLocation = env.ControlPlaneEnv().KubeConfigLocation
	c.ManagerConfigPath = env.ControlPlaneEnv().ManagerConfigLocation
	return c
}

func (b *Builder) unmanaged() (Cluster, error) {
	config, err := kubeconfig.FromPath(b.KubeConfigLocation)
	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to load Kubeconfig")
	}

	res, err := ClusterFromKubeConfig(config)
	if err != nil {
		return Cluster{}, errors.Wrap(err, "failed to create kind cluster from Kubeconfig")
	}

	return res, nil
}

// ForceCreate creates a new cluster based on the state of the Builder.
// Makes sure that old clusters with the same Name as this one are deleted.
func (b *Builder) ForceCreate() (ManagedCluster, error) {
	if b.Name == "" {
		return ManagedCluster{}, errors.New("trying to create a kind cluster with an empty Name (\"\")")
	}

	if err := b.Manager.DeleteCluster(b.Name); err != nil {
		return ManagedCluster{}, errors.Wrap(err, "failed to delete existing kind cluster")
	}
	res, err := b.Create()

	if err != nil {
		return ManagedCluster{}, errors.Wrap(err, "failed to create kind cluster")
	}

	return res, nil
}

// Create creates a new cluster based on the state of the Builder.
func (b *Builder) Create() (ManagedCluster, error) {
	if b.Name == "" {
		return ManagedCluster{}, errors.New("trying to create a cluster with an empty Name (\"\")")
	}

	// TODO: Should this happen here? we use this dir for other things right?
	// 		 Currently it's always fine because we always start by creating a managed cluster, but that won't be true forever
	if _, err := os.Stat(b.KubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(b.KubeConfigLocation), os.ModePerm); err != nil {
			return ManagedCluster{}, errors.Wrapf(err, "failed to create directory for kubeconfig (%v)", path.Dir(b.KubeConfigLocation))
		}
	}

	err := b.Manager.CreateCluster(b.Name, b.KubeConfigLocation, b.ManagerConfigPath)
	if err != nil {
		err = errors.Wrap(err, "failed to create kind cluster")

		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		if deleteClusterError := b.Manager.DeleteCluster(b.Name); deleteClusterError != nil {
			err = errors.Wrapf(deleteClusterError, "failed to delete kind cluster to clean up earlier failure (%v)", err)
		}
		return ManagedCluster{}, err
	}

	kubernetesCluster, err := b.unmanaged()
	if err != nil {
		return ManagedCluster{}, errors.Wrap(err, "failed to create unmanaged Kubernetes cluster")
	}

	return ManagedCluster{
		kubernetesCluster,
		b.Manager,
		b.Name,
	}, nil
}
