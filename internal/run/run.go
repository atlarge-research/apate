// Package run runs the apatelets based on the run policy determined in the environment
// Warning: This package directly imports apatelet, this should be taken into considerations
// when making changes to this package
package run

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	apateRun "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

// StartApatelets starts the apatelets based on the given parameters using the run type defined in the environment
func StartApatelets(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	runType := env.RetrieveFromEnvironment(env.ControlPlaneApateletRunType, env.ControlPlaneApateletRunTypeDefault)

	switch runType {
	case env.Docker:
		return useDocker(ctx, amountOfNodes, environment)
	case env.Routine:
		return useRoutines(ctx, amountOfNodes, environment)
	default:
		return fmt.Errorf("unknown run type: %v", runType)
	}
}

func useDocker(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	// Retrieve pull policy
	pullPolicy := env.RetrieveFromEnvironment(env.ControlPlaneDockerPolicy, env.ControlPlaneDockerPolicyDefault)
	fmt.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Spawn the apatelets
	return container.SpawnApateletContainers(ctx, amountOfNodes, pullPolicy, environment)
}

// TODO: Use ctx
func useRoutines(_ context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return err
	}

	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment.Copy()

		// TODO fix these ports
		apateletEnv.ListenPort = 7000 + i
		if apateletEnv.ListenPort == 8085 {
			apateletEnv.ListenPort = 6999
		}

		kubernetesPort := 12000 + i
		metricsPort := 17000 + i

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Apatelet failed to start: %v\n", r)
				}
			}()
			err := apateRun.StartApatelet(apateletEnv, kubernetesPort, metricsPort)
			panic(err)
		}()

		time.Sleep(time.Millisecond * 50)
	}
	return nil
}
