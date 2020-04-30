package cluster

type ManagedCluster struct {
	KubernetesCluster
	manager Manager
	name string
}

func (c ManagedCluster) Delete() error {
	return c.manager.DeleteCluster(c.name)
}
