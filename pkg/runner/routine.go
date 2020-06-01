package runner

import (
	"context"
	"log"

	"github.com/phayes/freeport"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	apateRun "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

// RoutineRunner spawns apatelets using go routines
type RoutineRunner struct{}

// SpawnApatelets spawns apatelets using go routines
func (d *RoutineRunner) SpawnApatelets(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return errors.Wrap(err, "failed to set certificates")
	}

	readyCh := make(chan struct{})

	environment.KubeConfigLocation = env.ControlPlaneEnv().KubeConfigLocation

	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment

		const numports = 3
		ports, err := freeport.GetFreePorts(numports)
		if err != nil {
			return errors.Wrapf(err, "failed to get %v free ports", numports)
		}

		apateletEnv.ListenPort = ports[0]
		apateletEnv.MetricsPort = ports[1]
		apateletEnv.KubernetesPort = ports[2]

		go func() {
			// TODO: Add retry logic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Apatelet failed to start: %v\n", r)
					readyCh <- struct{}{}
					i-- // retry
				}
			}()
			err := apateRun.StartApatelet(ctx, apateletEnv, readyCh)
			if err != nil {
				log.Printf("Apatelet failed to start: %v\n", err)
				readyCh <- struct{}{}
				i-- // retry
			}
		}()

		<-readyCh
	}
	return nil
}
