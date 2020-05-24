package container

import (
	"context"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnControlPlaneContainer spawns a single control plane container
func SpawnControlPlaneContainer(ctx context.Context, pullPolicy env.PullPolicy, cpEnv env.ControlPlaneEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "getting docker cli for spawning control plane container failed")
	}

	// Get docker port for control plane
	port, err := nat.NewPort("tcp", cpEnv.Port)

	if err != nil {
		return errors.Wrap(err, "failed to create docker port for Control plane")
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, env.ControlPlaneFullImage, env.ControlPlaneContainerName, 1, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: env.ControlPlaneImageName,
			Env: []string{
				env.ControlPlaneListenAddress + "=" + cpEnv.Address,
				env.ControlPlaneListenPort + "=" + cpEnv.Port,
				env.ManagedClusterConfig + "=" + cpEnv.ManagerConfig,
				env.ControlPlaneExternalIP + "=" + cpEnv.ExternalIP,
				env.ControlPlaneDockerPolicy + "=" + cpEnv.DockerPolicy,
				env.PrometheusStackEnabled + "=" + cpEnv.PrometheusStackEnabled,
			},
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: cpEnv.Port,
					},
				},
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
			},
		}, &network.NetworkingConfig{}, "apate-cp")

		if err != nil {
			return errors.Wrap(err, "failed to create Docker container for control plane")
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn control plane
	return errors.Wrap(HandleSpawnContainers(ctx, cli, spawnInfo), "failed to spawn containers for control plane")
}
