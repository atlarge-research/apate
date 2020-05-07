package container

import (
	"context"
	"strconv"

	"github.com/docker/go-connections/nat"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnApateletContainers spawns multiple Apatelet Docker containers
func SpawnApateletContainers(ctx context.Context, amountOfNodes int, pullPolicy string, cpEnv env.ApateletEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Get docker port for apatelet
	port, err := nat.NewPort("tcp", strconv.Itoa(cpEnv.ListenPort))

	if err != nil {
		return err
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, env.ApateletFullImage, env.ApateletContainerPrefix, amountOfNodes, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: env.ApateletImageName,
			Env: []string{
				env.ControlPlaneAddress + "=" + cpEnv.ControlPlaneAddress,
				env.ControlPlanePort + "=" + strconv.Itoa(cpEnv.ControlPlanePort),
				env.ApateletListenAddress + "=" + cpEnv.ListenAddress,
				env.ApateletListenPort + "=" + strconv.Itoa(cpEnv.ListenPort),
			},
			// TODO: Make variable
			ExposedPorts: nat.PortSet{
				port:        struct{}{},
				"10250/tcp": struct{}{},
				"10255/tcp": struct{}{},
			},
		}, nil, nil, env.ApateletContainerPrefix+strconv.Itoa(i))

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn apatelets
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
