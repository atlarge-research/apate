package runner

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	apateRun "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

// RoutineRunner spawns apatelets using go routines
type RoutineRunner struct{}

// SpawnApatelets spawns apatelets using go routines
func (d RoutineRunner) SpawnApatelets(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return errors.Wrap(err, "failed to set certificates")
	}

	environment.KubeConfigLocation = env.ControlPlaneEnv().KubeConfigLocation

	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment
		readyCh := make(chan struct{}, 1)

		go func() {
			// TODO: Add retry logic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Apatelet failed to start: %v\n", r)
					readyCh <- struct{}{}
				}
			}()
			err := apateRun.StartApatelet(ctx, apateletEnv, readyCh)
			if err != nil {
				log.Printf("Apatelet failed to start: %v\n", err)
				readyCh <- struct{}{}
			}
		}()

		<-readyCh
	}
	return nil
}
