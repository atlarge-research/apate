package runner

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/atlarge-research/apate/pkg/container"
	"github.com/atlarge-research/apate/pkg/env"
)

// DockerRunner runs the apatelets using docker containers
type DockerRunner struct{}

// SpawnApatelets spawns the apatelets using docker containers
func (d *DockerRunner) SpawnApatelets(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	// Retrieve pull policy
	pullPolicy := env.ControlPlaneEnv().DockerPolicy
	fmt.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Spawn the apatelets
	return errors.Wrap(container.SpawnApateletContainers(ctx, amountOfNodes, pullPolicy, environment), "failed to spawn Apatelet docker containers")
}
