package cluster

// A ManagedCluster inherits all it's methods from a
// KubernetesCluster, but is managed. This means it can be deleted.
// A ManagedCluster is guaranteed to be created by Apate, and can
// therefore also be safely be deleted by Apate.
type ManagedCluster struct {
	KubernetesCluster
	manager Manager
	name string
}

// Delete destroys a (managed) kubernetes cluster
func (c ManagedCluster) Delete() error {
	return c.manager.DeleteCluster(c.name)
}
