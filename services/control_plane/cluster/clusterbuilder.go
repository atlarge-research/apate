package cluster

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
)

// The ClusterBuilder creates a new cluster object used to manage a cluster.
type ClusterBuilder struct {
	name               string
	creator            ClusterCreator
	kubeConfigLocation string
}

// The New function is used to create a new ClusterBuilder with all fields empty.
func New() (c ClusterBuilder) {
	return c
}

// The Default function is used to create a new ClusterBuilder with all fields set to default values.
func Default() (c ClusterBuilder) {
	c.name = "Apate"
	c.creator = &KinD{}
	c.kubeConfigLocation = os.TempDir() + "/apate/config"
	return c
}

// The WithName function is used to give the cluster that is to be built a name.
func (b *ClusterBuilder) WithName(name string) *ClusterBuilder {
	b.name = name
	return b
}

// The WithConfigLocation function is used to give the cluster that is to be built a name.
func (b *ClusterBuilder) WithConfigLocation(kubeconfiglocation string) *ClusterBuilder {
	b.kubeConfigLocation = kubeconfiglocation
	return b
}

// The WithCreator function is used to enable the cluster to be built with a different
// cluster creator.
func (b *ClusterBuilder) WithCreator(creator ClusterCreator) *ClusterBuilder {
	b.creator = creator
	return b
}

// Creates a new cluster based on the state of the ClusterBuilder.
// Makes sure that old clusters with the same name as this one are deleted.
func (b *ClusterBuilder) ForceCreate() (KubernetesCluster, error) {
	if b.name == "" {
		return KubernetesCluster{}, errors.New("Tying to create a cluster with an empty name (\"\")")
	}

	b.creator.DeleteCluster(b.name)
	return b.Create()
}

// Creates a new cluster based on the state of the ClusterBuilder.
func (b *ClusterBuilder) Create() (KubernetesCluster, error) {
	if b.name == "" {
		return KubernetesCluster{}, errors.New("Tying to create a cluster with an empty name (\"\")")
	}

	if _, err := os.Stat(b.kubeConfigLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(b.kubeConfigLocation), os.ModePerm); err != nil {
			return KubernetesCluster{}, err
		}
	}

	err := b.creator.CreateCluster(b.name, b.kubeConfigLocation)
	if err != nil {
		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		b.creator.DeleteCluster(b.name)
		return KubernetesCluster{}, err
	}

	config, err := getConfigForContext(b.creator.ClusterContext(b.name), b.kubeConfigLocation)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		b.creator.DeleteCluster(b.name)
		return KubernetesCluster{}, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		b.creator.DeleteCluster(b.name)
		return KubernetesCluster{}, err
	}

	return KubernetesCluster{
		name:                     b.name,
		clientSet:                clientSet,
		clusterCreationInterface: b.creator,
	}, nil
}

// Gets a kubernetes client configuration for the context given.
func getConfigForContext(context string, kubeConfigLocation string) (*rest.Config, error) {
	// Create a default config rules struct
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	rules.ExplicitPath = kubeConfigLocation

	// Override with defaults (this call might not be necessary since the defaults are already set above?)
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	// But set the context to our own context while overriding
	overrides.CurrentContext = context

	// Now create the actual configuration
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
