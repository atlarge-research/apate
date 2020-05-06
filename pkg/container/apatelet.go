package container

import (
	"context"
	ec "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
	"github.com/docker/go-connections/nat"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnApateletContainers spawns multiple Apatelet Docker containers
func SpawnApateletContainers(ctx context.Context, amountOfNodes int, pullPolicy string, env ec.ApateletEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Get docker port for apatelet
	port, err := nat.NewPort("tcp", strconv.Itoa(env.ListenPort))

	if err != nil {
		return err
	}

	useDocker := false

	run.SetCerts()

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, ec.ApateletFullImage, ec.ApateletContainerPrefix, amountOfNodes, func(i int, ctx context.Context) error {
		if useDocker {
			c, err := cli.ContainerCreate(ctx, &container.Config{
				Image: ec.ApateletImageName,
				Env: []string{
					ec.ControlPlaneAddress + "=" + env.ControlPlaneAddress,
					ec.ControlPlanePort + "=" + strconv.Itoa(env.ControlPlanePort),
					ec.ApateletListenAddress + "=" + env.ListenAddress,
					ec.ApateletListenPort + "=" + strconv.Itoa(env.ListenPort),
				},
				ExposedPorts: nat.PortSet{
					port:        struct{}{},
					"10250/tcp": struct{}{},
				},
			}, nil, nil, ec.ApateletContainerPrefix+strconv.Itoa(i))

			if err != nil {
				return err
			}

			return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
		} else {
			env.ListenPort = 7000 + i
			if env.ListenPort == 8085 {
				env.ListenPort = 6999
			}

			kubernetesPort := 12000 + i
			metricsPort := 17000 + i

			go run.StartApatelet(env, kubernetesPort, metricsPort)
			return nil
		}
	})

	// Spawn apatelets
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
