package kubernetes

import (
	"os"
	"path"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kind"
)

// Builder allows for the creation of KubernetesClusters.
type Builder struct {
	name               string
	manager            Manager
	kubeConfigLocation string
	managerConfigPath  string
}

// New is used to create a new Builder with all fields empty.
func New() Builder {
	return Builder{}
}

// Default is used to create a new Builder with all fields set to default values.
func Default() (c Builder) {
	c.name = "Apate"
	c.manager = &kind.KinD{}
	c.kubeConfigLocation = env.ControlPlaneEnv().KubeConfigLocation
	return c
}

// WithName is used to give the cluster that is to be built a name.
func (b *Builder) WithName(name string) *Builder {
	b.name = name
	return b
}

// WithManagerConfig is used to set the path to the config for the manager, if applicable
func (b *Builder) WithManagerConfig(path string) *Builder {
	b.managerConfigPath = path
	return b
}

// WithConfigLocation is used to give the cluster that is to be built a name.
func (b *Builder) WithConfigLocation(kubeConfigLocation string) *Builder {
	b.kubeConfigLocation = kubeConfigLocation
	return b
}

// WithCreator is used to enable the cluster to be built with a different
// cluster manager.
func (b *Builder) WithCreator(creator Manager) *Builder {
	b.manager = creator
	return b
}

func (b *Builder) unmanaged() (Cluster, error) {
	config, err := kubeconfig.FromPath(b.kubeConfigLocation)
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
// Makes sure that old clusters with the same name as this one are deleted.
func (b *Builder) ForceCreate() (ManagedCluster, error) {
	if b.name == "" {
		return ManagedCluster{}, errors.New("trying to create a kind cluster with an empty name (\"\")")
	}

	if err := b.manager.DeleteCluster(b.name); err != nil {
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
	if b.name == "" {
		return ManagedCluster{}, errors.New("trying to create a cluster with an empty name (\"\")")
	}

	// TODO: Should this happen here? we use this dir for other things right?
	// 		 Currently it's always fine because we always start by creating a managed cluster, but that won't be true forever
	if _, err := os.Stat(b.kubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(b.kubeConfigLocation), os.ModePerm); err != nil {
			return ManagedCluster{}, errors.Wrapf(err, "failed to create directory for kubeconfig (%v)", path.Dir(b.kubeConfigLocation))
		}
	}

	err := b.manager.CreateCluster(b.name, b.kubeConfigLocation, b.managerConfigPath)
	if err != nil {
		err = errors.Wrap(err, "failed to create kind cluster")

		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		if deleteClusterError := b.manager.DeleteCluster(b.name); deleteClusterError != nil {
			err = errors.Wrapf(deleteClusterError, "failed to delete kind cluster to clean up earlier failure: (%v)", err)
		}
		return ManagedCluster{}, err
	}

	kubernetesCluster, err := b.unmanaged()
	if err != nil {
		return ManagedCluster{}, errors.Wrap(err, "failed to create unmanaged Kubernetes cluster")
	}

	return ManagedCluster{
		kubernetesCluster,
		b.manager,
		b.name,
	}, nil
}
