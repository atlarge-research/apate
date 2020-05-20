// Package run runs the apatelets based on the run policy determined in the environment
// Warning: This package directly imports apatelet, this should be taken into considerations
// when making changes to this package
package run

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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
		return errors.Wrap(useDocker(ctx, amountOfNodes, environment), "failed to use docker to start Apatelets")
	case env.Routine:
		return errors.Wrap(useRoutines(amountOfNodes, environment), "failed to use goroutines to start Apatelets")
	default:
		return fmt.Errorf("unknown run type: %v", runType)
	}
}

func useDocker(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	// Retrieve pull policy
	pullPolicy := env.RetrieveFromEnvironment(env.ControlPlaneDockerPolicy, env.ControlPlaneDockerPolicyDefault)
	fmt.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Spawn the apatelets
	return errors.Wrap(container.SpawnApateletContainers(ctx, amountOfNodes, pullPolicy, environment), "failed to spawn Apatelet docker containers")
}

func useRoutines(amountOfNodes int, environment env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return errors.Wrap(err, "failed to set certificates")
	}

	readyCh := make(chan bool)

	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment
		const attemts = 3
		ports, err := freeport.GetFreePorts(attemts)

		if err != nil {
			return errors.Wrapf(err, "failed to get a free port after %v attempts", attemts)
		}

		apateletEnv.ListenPort = ports[0]

		go func() {
			// TODO: Add retry logic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Apatelet failed to start: %+v\n", r)
				}
			}()
			err := apateRun.StartApatelet(apateletEnv, ports[1], ports[2], &readyCh)
			if err != nil {
				log.Printf("Apatelet failed to start: %+v\n", err)
			}
		}()

		<-readyCh
	}
	return nil
}
