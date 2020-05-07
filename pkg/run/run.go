package run

import (
	"context"
	"fmt"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"log"
)

type ApateletStartFunction = func(env.ApateletEnvironment, int, int) error

func StartApatelets(ctx context.Context, amountOfNodes int, environment env.ApateletEnvironment, fun ApateletStartFunction) error {
	docker := false // TODO: Get from env

	if docker {
		return useDocker(ctx, amountOfNodes, environment)
	} else {
		return useRoutines(ctx, amountOfNodes, environment, fun)
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
func useRoutines(_ context.Context, amountOfNodes int, environment env.ApateletEnvironment, fun ApateletStartFunction) error {
	for i := 0; i < amountOfNodes; i++ {
		apateletEnv := environment.Copy()
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
			err := fun(apateletEnv, kubernetesPort, metricsPort)
			panic(err)
		}()
	}
	return nil
}
