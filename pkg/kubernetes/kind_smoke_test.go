package kubernetes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
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

	clusterBuilder := NewClusterManagerHandler()
	c, err := clusterBuilder.NewCluster()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err := c.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)
	assert.NoError(t, c.Shutdown())
}

func TestForceCreateCluster(t *testing.T) {
	setup(t)

	clusterInterface := KinDClusterManager{}

	// Delete it before to be safe
	assert.NoError(t, clusterInterface.Shutdown(nil))

	clusterBuilder := NewClusterManagerHandler()

	// Create a cluster
	c, err := clusterBuilder.NewCluster()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err := c.GetNumberOfNodes()
	assert.NoError(t, err)

	assert.Equal(t, 1, nodes)

	// Now force create one. This should not error but instead delete the old one.
	cb := NewClusterManagerHandler()
	c, err = cb.NewCluster()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	nodes, err = c.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 1, nodes)

	assert.NoError(t, c.Shutdown())
}
