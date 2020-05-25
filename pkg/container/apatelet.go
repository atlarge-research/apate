package container

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/docker/go-connections/nat"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnApateletContainers spawns multiple Apatelet Docker containers
func SpawnApateletContainers(ctx context.Context, amountOfNodes int, pullPolicy env.PullPolicy, apateletEnv env.ApateletEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create a new Docker client")
	}

	// Create docker ports
	ports, err := createPortSet(apateletEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create ports")
	}

	// Dump environment as string array
	envArray, err := env.DumpAsKeyValue(apateletEnv)
	if err != nil {
		return errors.Wrap(err, "failed to dump apatelet environment to strings")
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, env.ApateletFullImage, env.ApateletContainerPrefix, amountOfNodes, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        env.ApateletImageName,
			Env:          envArray,
			ExposedPorts: ports,
		}, nil, nil, env.ApateletContainerPrefix+strconv.Itoa(i))

		if err != nil {
			return errors.Wrapf(err, "failed to create container %v", env.ApateletContainerPrefix+strconv.Itoa(i))
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn apatelets
	return errors.Wrap(HandleSpawnContainers(ctx, cli, spawnInfo), "failed to spawn Apatelet containers")
}

func createPortSet(apateletEnv env.ApateletEnvironment) (nat.PortSet, error) {
	// Get docker port for apatelet
	lp, err := nat.NewPort("tcp", strconv.Itoa(apateletEnv.ListenPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker listen port")
	}

	// Get docker port for metrics
	mp, err := nat.NewPort("tcp", strconv.Itoa(apateletEnv.MetricsPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker metric port")
	}

	// Get docker port for kubernetes
	kp, err := nat.NewPort("tcp", strconv.Itoa(apateletEnv.KubernetesPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker kubernetes port")
	}

	return nat.PortSet{
		lp: struct{}{},
		mp: struct{}{},
		kp: struct{}{},
	}, nil
}
