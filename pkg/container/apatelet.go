package container

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	Address, Port string
}

// DefaultApateEnvironment returns the default apate environment
func DefaultApateEnvironment() ApateletEnvironment {
	return ApateletEnvironment{
		Address: ApateletListenAddressDefault,
		Port:    ApateletListenPortDefault,
	}
}

// SpawnApatelets spawns multiple Apatelet Docker containers
func SpawnApatelets(ctx context.Context, amountOfNodes int, info *service.ConnectionInfo, pullPolicy string, env ApateletEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, apateletFullImage, apateletContainerPrefix, amountOfNodes, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: apateletImageName,
			Env: []string{
				ControlPlaneAddress + "=" + info.Address,
				ControlPlanePort + "=" + strconv.Itoa(info.Port),
				ApateletListenAddress + "=" + env.Address,
				ApateletListenPort + "=" + env.Port,
			},
		}, nil, nil, apateletContainerPrefix+strconv.Itoa(i))

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	err = HandleSpawnContainers(ctx, cli, spawnInfo)

	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		// TODO: Add retry functionality
	}

	return nil
}
