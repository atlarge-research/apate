package cluster

import (
	"errors"
	"os"
	"path"
)

// Builder allows for the creation of KubernetesClusters.
type Builder struct {
	name               string
	manager            Manager
	kubeConfigLocation string
}

// New is used to create a new Builder with all fields empty.
func New() Builder {
	return Builder{}
}

// Default is used to create a new Builder with all fields set to default values.
func Default() (c Builder) {
	c.name = "Apate"
	c.manager = &KinD{}
	c.kubeConfigLocation = os.TempDir() + "/apate/config"
	return c
}

// WithName is used to give the cluster that is to be built a name.
func (b *Builder) WithName(name string) *Builder {
	b.name = name
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

func (b *Builder) unManaged() (KubernetesCluster, error) {
	config, err := GetKubeConfig(b.kubeConfigLocation)
	if err != nil {
		return KubernetesCluster{}, err
	}

	return KubernetesClusterFromKubeConfig(config)
}

// ForceCreate creates a new cluster based on the state of the Builder.
// Makes sure that old clusters with the same name as this one are deleted.
func (b *Builder) ForceCreate() (ManagedCluster, error) {
	if b.name == "" {
		return ManagedCluster{}, errors.New("trying to create a cluster with an empty name (\"\")")
	}

	if err := b.manager.DeleteCluster(b.name); err != nil {
		return ManagedCluster{}, err
	}
	return b.Create()
}

// Create creates a new cluster based on the state of the Builder.
func (b *Builder) Create() (ManagedCluster, error) {
	if b.name == "" {
		return ManagedCluster{}, errors.New("trying to create a cluster with an empty name (\"\")")
	}

	if _, err := os.Stat(b.kubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(b.kubeConfigLocation), os.ModePerm); err != nil {
			return ManagedCluster{}, err
		}
	}

	err := b.manager.CreateCluster(b.name, b.kubeConfigLocation)
	if err != nil {
		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		if err1 := b.manager.DeleteCluster(b.name); err1 != nil {
			err = err1
		}
		return ManagedCluster{}, err
	}

	kubernetesCluster, err := b.unManaged()
	if err != nil {
		return ManagedCluster{}, err
	}

	return ManagedCluster{
		kubernetesCluster,
		b.manager,
		b.name,
	}, nil
}
