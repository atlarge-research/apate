package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

func TestSpawnControlPlaneDocker(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E")
	}

	ctx, cancel := context.WithCancel(context.Background())

	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)

	if hasContainer, c := hasContainer(ctx, t, cli); hasContainer {
		second := 1 * time.Second
		err = cli.ContainerStop(ctx, c.ID, &second)
		assert.NoError(t, err)

		time.Sleep(second)

		err = cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
		assert.NoError(t, err)
	}

	err = container.SpawnControlPlaneContainer(ctx, env.PullIfNotLocal, env.ControlPlaneEnv())
	assert.NoError(t, err)

	hasContainer, c := hasContainer(ctx, t, cli)
	assert.True(t, hasContainer)
	assert.NotNil(t, c)

	assert.Len(t, c.Ports, 1)
	assert.Equal(t, uint16(8085), c.Ports[0].PrivatePort)
	assert.Equal(t, uint16(8085), c.Ports[0].PublicPort)
	assert.Equal(t, "0.0.0.0", c.Ports[0].IP)
	assert.Equal(t, "tcp", c.Ports[0].Type)

	assert.Len(t, c.Names, 1)
	assert.Equal(t, "/"+env.ControlPlaneContainerName, c.Names[0])

	inspect, err := cli.ContainerInspect(ctx, c.ID)
	assert.NoError(t, err)

	envVars := inspect.Config.Env
	assert.Contains(t, envVars, "CP_EXTERNAL_IP=auto")
	assert.Contains(t, envVars, "CP_LISTEN_ADDRESS=0.0.0.0")
	assert.Contains(t, envVars, "CP_ENABLE_DEBUG=false")
	assert.Contains(t, envVars, "CP_KUBE_CONFIG_LOCATION=/tmp/apate/config")
	assert.Contains(t, envVars, "CP_DOCKER_POLICY=pull-if-not-local")
	assert.Contains(t, envVars, "CP_KIND_CLUSTER_NAME=apate")
	assert.Contains(t, envVars, "CP_PROMETHEUS_CONFIG_LOCATION=config/prometheus.yml")
	assert.Contains(t, envVars, "CP_NODE_CRD_LOCATION=config/crd/apate.opendc.org_nodeconfigurations.yaml")
	assert.Contains(t, envVars, "CP_LISTEN_PORT=8085")
	assert.Contains(t, envVars, "CP_MANAGER_CONFIG_LOCATION=config/kind.yml")
	assert.Contains(t, envVars, "CP_KUBE_CONFIG=")
	assert.Contains(t, envVars, "CP_APATELET_RUN_TYPE=ROUTINES")
	assert.Contains(t, envVars, "CP_PROMETHEUS=true")
	assert.Contains(t, envVars, "CP_PROMETHEUS_NAMESPACE=apate-prometheus")
	assert.Contains(t, envVars, "CP_POD_CRD_LOCATION=config/crd/apate.opendc.org_podconfigurations.yaml")

	assert.Equal(t, "/app/controlplane", inspect.Config.Cmd[0])

	cancel()
}

func hasContainer(ctx context.Context, t *testing.T, cli *client.Client) (bool, *types.Container) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	assert.NoError(t, err)

	for _, c := range containers {
		if c.Image == env.ControlPlaneFullImage {
			return true, &c
		}
	}

	return false, nil
}
