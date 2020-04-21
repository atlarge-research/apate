package cluster

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// The ClusterBuilder creates a new cluster object used to manage a cluster.
type ClusterBuilder struct {
	name string
}

// The New function is used to create a new ClusterBuilder with all fields empty.
func New() (c ClusterBuilder) {
	return c
}

// The Default function is used to create a new ClusterBuilder with all fields set to default values.
func Default() (c ClusterBuilder) {
	c.name = "Apate"
	return c
}

// The WithName function is used to give the cluster that is to be built a name.
func (b * ClusterBuilder) WithName(name string) * ClusterBuilder {
	b.name = name
	return b
}

// Creates a new cluster based on the state of the ClusterBuilder.
// Makes sure that old clusters with the same name as this one are deleted.
func (b *ClusterBuilder) ForceCreate() (Cluster, error) {
	if b.name == "" {
		return Cluster{}, errors.New("Tying to create a cluster with an empty name (\"\")")
	}

	DeleteCluster(b.name)
	return b.Create()
}

// Creates a new cluster based on the state of the ClusterBuilder.
func (b *ClusterBuilder) Create() (Cluster, error) {
	if b.name == "" {
		return Cluster{}, errors.New("Tying to create a cluster with an empty name (\"\")")
	}

	err := CreateCluster(b.name)
	if err != nil {
		// If something went wrong, there still could be a built cluster we can't interact with.
		// delete the cluster to be safe for the next run, otherwise ForceCreate would be necessary
		DeleteCluster(b.name)
		return Cluster{}, err
	}

	config, err := getConfigForContext("kind-" + b.name)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		DeleteCluster(b.name)
		return Cluster{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// If something went wrong, delete the cluster for the next run,
		// otherwise ForceCreate would be necessary
		DeleteCluster(b.name)
		return Cluster{}, err
	}

	return Cluster {
		name: b.name,
		clientset: clientset,
	}, nil
}

// Gets a kubernetes client configuration for the context given.
func getConfigForContext(context string) (*rest.Config, error) {
	// Create a default config rules struct
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

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

