package kubernetes

import (
	"os"
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kind"

	"github.com/stretchr/testify/assert"
)

// I specifically give clusters the name of their test,
// so even if tests are ran in parallel there won't be a problem.

func TestCreateCluster(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	clusterInterface := kind.KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster("TestCreateCluster"))

	clusterBuilder := Default()
	cluster, err := clusterBuilder.WithCreator(&clusterInterface).WithName("TestCreateCluster").Create()
	assert.NotNil(t, cluster)
	assert.NoError(t, err)

	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)
	assert.NoError(t, cluster.Delete())
}

func TestCreateClusterNoFolder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	clusterInterface := kind.KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster("TestCreateClusterNoFolder"))

	_ = os.RemoveAll("/tmp/TestCreateClusterNoFolder_e2e")

	clusterBuilder := Default()
	cluster, err := clusterBuilder.
		WithCreator(&clusterInterface).
		WithName("TestCreateClusterNoFolder").
		WithConfigLocation("/tmp/TestCreateClusterNoFolder_e2e").
		Create()

	assert.NotNil(t, cluster)
	assert.NoError(t, err)

	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)
	assert.NoError(t, cluster.Delete())

	_ = os.RemoveAll("/tmp/TestCreateClusterNoFolder_e2e")
}

func TestForceCreateCluster(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	clusterInterface := kind.KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster("TestForceCreateCluster"))

	clusterBuilder := Default()
	// Create a cluster
	cluster, err := clusterBuilder.WithCreator(&clusterInterface).WithName("TestForceCreateCluster").Create()
	assert.NotNil(t, cluster)
	assert.NoError(t, err)

	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)

	// Now create another one. This should error
	_, err = clusterBuilder.WithName("TestForceCreateCluster").Create()
	assert.Error(t, err)

	// Now force create one. This should not error but instead delete the old one.
	cluster, err = clusterBuilder.WithName("TestForceCreateCluster").ForceCreate()
	assert.NotNil(t, cluster)
	assert.NoError(t, err)

	nodes, err = cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 1, nodes)

	assert.NoError(t, cluster.Delete())
}
