package cluster

import (
	"errors"
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	c.kubeConfigLocation = os.TempDir() + "/apate/scenario"
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

// ForceCreate creates a new cluster based on the state of the Builder.
// Makes sure that old clusters with the same name as this one are deleted.
func (b *Builder) ForceCreate() (KubernetesCluster, error) {
	if b.name == "" {
		return KubernetesCluster{}, errors.New("trying to create a cluster with an empty name (\"\")")
	}

	if err := b.manager.DeleteCluster(b.name); err != nil {
		return KubernetesCluster{}, err
	}
	return b.Create()
}

// Create creates a new cluster based on the state of the Builder.
func (b *Builder) Create() (KubernetesCluster, error) {
	if b.name == "" {
		return KubernetesCluster{}, errors.New("trying to create a cluster with an empty name (\"\")")
	}

	if _, err := os.Stat(b.kubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(b.kubeConfigLocation), os.ModePerm); err != nil {
			return KubernetesCluster{}, err
		}
	}

	err := b.manager.CreateCluster(b.name, b.kubeConfigLocation)
	if err != nil {
		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		if err1 := b.manager.DeleteCluster(b.name); err1 != nil {
			err = err1
		}
		return KubernetesCluster{}, err
	}

	config, err := getConfigForContext(b.manager.ClusterContext(b.name), b.kubeConfigLocation)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		if err1 := b.manager.DeleteCluster(b.name); err1 != nil {
			err = err1
		}
		return KubernetesCluster{}, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		if err1 := b.manager.DeleteCluster(b.name); err1 != nil {
			err = err1
		}
		return KubernetesCluster{}, err
	}

	return KubernetesCluster{
		name:      b.name,
		clientSet: clientSet,
		manager:   b.manager,
	}, nil
}

// Gets a kubernetes client scenario for the context given.
func getConfigForContext(context string, kubeConfigLocation string) (*rest.Config, error) {
	// Create a default scenario rules struct
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	rules.ExplicitPath = kubeConfigLocation

	// Override with defaults (this call might not be necessary since the defaults are already set above?)
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	// But set the context to our own context while overriding
	overrides.CurrentContext = context

	// Now create the actual scenario
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
