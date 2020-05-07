// Package run runs the apatelets based on the run policy determined in the environment
// Warning: This package directly imports apatelet, this should be taken into considerations
// when making changes to this package
package run

import (
	"context"
	"fmt"
	"log"

	"github.com/phayes/freeport"

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
		return useRoutines(amountOfNodes, environment)
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

func useRoutines(amountOfNodes int, environment env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return err
	}

	readyCh := make(chan bool)

	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment
		ports, err := freeport.GetFreePorts(3)

		if err != nil {
			return err
		}

		apateletEnv.ListenPort = ports[0]

		go func() {
			// TODO: Add retry logic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Apatelet failed to start: %v\n", r)
				}
			}()
			err := apateRun.StartApatelet(apateletEnv, ports[1], ports[2], &readyCh)
			if err != nil {
				log.Printf("Apatelet failed to start: %v\n", err)
			}
		}()

		<-readyCh
	}
	return nil
}
