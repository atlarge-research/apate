package cluster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCluster(t *testing.T) {
	// Delete it before to be safe
	DeleteCluster("Apate")

	clusterbuilder := New()
	cluster, err := clusterbuilder.WithName("Apate").Create()

	assert.NoError(t, err)

	cluster.Delete()
}


func TestForceCreateCluster(t *testing.T) {
	// Delete it before to be safe
	DeleteCluster("Apate")

	clusterbuilder := New()
	// Create a cluster
	_, err := clusterbuilder.WithName("Apate").Create()
	assert.NoError(t, err)

	// Now create another one. This should error
	_, err = clusterbuilder.WithName("Apate").Create()
	assert.Error(t, err)

	// Now force create one. This should not error but instead delete the old one.
	cluster, err := clusterbuilder.WithName("Apate").ForceCreate()
	assert.NoError(t, err)

	// TODO: There's currently no way to test if the old cluster was actually deleted (but it kinda has to be)
	//		 but more importantly: it might be useful to have a way that other Cluster structs with the same name
	//		 as a force created one are marked as invalid so they can't be used to interact with the cluster anymore.

	cluster.Delete()
}