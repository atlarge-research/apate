package cluster

import (
	"control_plane/cluster/clustercreationinterface"
	"github.com/stretchr/testify/assert"
	"testing"
)

// I specifically give clusters the name of their test,
// so even if tests are ran in parallel there won't be a problem.

func TestCreateCluster_e2e(t *testing.T) {
	clusterInterface := clustercreationinterface.Kind{}

	// Delete it before to be safe
	clusterInterface.DeleteCluster("TestCreateCluster")

	clusterbuilder := New()
	cluster, err := clusterbuilder.WithClusterCreationInterface(&clusterInterface).WithName("TestCreateCluster").Create()

	assert.NoError(t, err)

	cluster.Delete()
}


func TestForceCreateCluster_e2e(t *testing.T) {
	clusterInterface := clustercreationinterface.Kind{}

	// Delete it before to be safe
	clusterInterface.DeleteCluster("TestForceCreateCluster")

	clusterbuilder := New()
	// Create a cluster
	_, err := clusterbuilder.WithClusterCreationInterface(&clusterInterface).WithName("TestForceCreateCluster").Create()
	assert.NoError(t, err)


	// Now create another one. This should error
	_, err = clusterbuilder.WithName("TestForceCreateCluster").Create()
	assert.Error(t, err)

	// Now force create one. This should not error but instead delete the old one.
	cluster, err := clusterbuilder.WithName("TestForceCreateCluster").ForceCreate()
	assert.NoError(t, err)

	// TODO: There's currently no way to test if the old cluster was actually deleted (but it kinda has to be)
	//		 but more importantly: it might be useful to have a way that other Cluster structs with the same name
	//		 as a force created one are marked as invalid so they can't be used to interact with the cluster anymore.

	cluster.Delete()
}