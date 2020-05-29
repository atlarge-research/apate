package kubernetes

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kind"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	clusterbuilder := Default()

	assert.Equal(t, env.ControlPlaneEnv().KinDClusterName, clusterbuilder.Name)
	assert.Equal(t, env.ControlPlaneEnv().ManagerConfigLocation, clusterbuilder.ManagerConfigPath)
	assert.Equal(t, env.ControlPlaneEnv().KubeConfigLocation, clusterbuilder.KubeConfigLocation)
	assert.Equal(t, &kind.KinD{}, clusterbuilder.Manager)
}
