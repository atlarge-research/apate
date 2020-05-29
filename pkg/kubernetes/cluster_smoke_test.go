package kubernetes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kind"
)

func setup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	dir := os.Getenv("CI_PROJECT_DIR")
	if len(dir) == 0 {
		dir = "../../"
	}

	initEnv := env.ControlPlaneEnv()
	initEnv.ManagerConfigLocation = dir + "/config/gitlab-kind.yml"
	env.SetEnv(initEnv)
}

// I specifically give clusters the name of their test,
// so even if tests are ran in parallel there won't be a problem.

func TestCreateCluster(t *testing.T) {
	setup(t)

	clusterBuilder := Default()
	clusterBuilder.Name = "TestCreateCluster"
	c, err := clusterBuilder.ForceCreate()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err := c.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)
	assert.NoError(t, c.Delete())
}

func TestCreateClusterNoFolder(t *testing.T) {
	path := "/tmp/TestCreateClusterNoFolder_e2e"
	setup(t)

	_ = os.RemoveAll(path)

	clusterBuilder := Default()
	clusterBuilder.Name = "TestCreateClusterNoFolder"
	clusterBuilder.KubeConfigLocation = path + "/config.yml"
	c, err := clusterBuilder.ForceCreate()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	f, err := os.Stat(path)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err := c.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)
	assert.NoError(t, c.Delete())

	_ = os.RemoveAll(path)
}

func TestForceCreateCluster(t *testing.T) {
	name := "TestForceCreateCluster"
	setup(t)

	clusterInterface := kind.KinD{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.DeleteCluster(name))

	clusterBuilder := Default()
	clusterBuilder.Name = name

	// Create a cluster
	c, err := clusterBuilder.Create()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err := c.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)

	// Now create another one. This should error
	cb := Default()
	cb.Name = name
	c, err = cb.Create()
	assert.Error(t, err)

	// Now force create one. This should not error but instead delete the old one.
	cb = Default()
	cb.Name = name
	c, err = cb.ForceCreate()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err = c.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 1, nodes)

	assert.NoError(t, c.Delete())
}
