// Package container provides methods to create containers for the required Apate components
package container

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

const (
	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress = "CP_LISTEN_ADDRESS"

	// ControlPlanePort is the port of the control plane
	ControlPlanePort = "CP_LISTEN_PORT"

	// ControlPlaneDockerPolicy specifies the docker pull policy for apatelet images
	ControlPlaneDockerPolicy = "CP_DOCKER_POLICY"

	// General apate docker constant
	apateDocker = "apatekubernetes"

	// Apatelet docker constants
	apateletContainerPrefix = "apatelet-"
	apateletImageName       = "apatelet:latest"
	apateletFullImage       = apateDocker + "/" + apateletImageName

	// Docker docker constants
	controlPlaneContainerName = "apate-cp"
	controlPlaneImageName     = "controlplane:latest"
	controlPlaneFullImage     = apateDocker + "/" + controlPlaneImageName
)

// SpawnApatelets spawns multiple Apatelet Docker containers
func SpawnApatelets(ctx context.Context, amountOfNodes int, info *service.ConnectionInfo, pullPolicy string) error {
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
			},
		}, nil, nil, apateletContainerPrefix+strconv.Itoa(i))

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn apatelets
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}

// SpawnControlPlane spawns a single control plane container
func SpawnControlPlane(ctx context.Context, pullPolicy string) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, controlPlaneFullImage, controlPlaneContainerName, 1, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: controlPlaneImageName,
		}, nil, nil, controlPlaneContainerName)

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn control plane
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
