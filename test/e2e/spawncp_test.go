package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

func TestSpawnControlPlaneDocker(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E")
	}

	ctx, cancel := context.WithCancel(context.Background())

	err := container.SpawnControlPlaneContainer(ctx, env.PullIfNotLocal, env.ControlPlaneEnv())
	assert.NoError(t, err)

	waitForCP(t)

	cancel()
}