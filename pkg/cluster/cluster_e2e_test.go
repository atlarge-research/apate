package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// I specifically give clusters the name of their test,
// so even if tests are ran in parallel there won't be a problem.

func TestCreateCluster_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	clusterInterface := KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster("TestCreateCluster"))

	clusterBuilder := Default()
	cluster, err := clusterBuilder.WithCreator(&clusterInterface).WithName("TestCreateCluster").Create()

	assert.NoError(t, err)

	assert.NoError(t, cluster.Delete())
}

func TestForceCreateCluster_e2e(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	clusterInterface := KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster("TestForceCreateCluster"))

	clusterBuilder := Default()
	// Create a cluster
	_, err := clusterBuilder.WithCreator(&clusterInterface).WithName("TestForceCreateCluster").Create()
	assert.NoError(t, err)

	// Now create another one. This should error
	_, err = clusterBuilder.WithName("TestForceCreateCluster").Create()
	assert.Error(t, err)

	// Now force create one. This should not error but instead delete the old one.
	cluster, err := clusterBuilder.WithName("TestForceCreateCluster").ForceCreate()
	assert.NoError(t, err)

	// TODO: There's currently no way to test if the old cluster was actually deleted (but it kinda has to be)
	//		 but more importantly: it might be useful to have a way that other Cluster structs with the same name
	//		 as a force created one are marked as invalid so they can't be used to interact with the cluster anymore.

	assert.NoError(t, cluster.Delete())
}
