package container

import (
	"context"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ControlPlaneEnvironment represents the environment variables of the control plane
type ControlPlaneEnvironment struct {
	Address, Port, ManagerConfig, ExternalIP, DockerPolicy string
}

// DefaultControlPlaneEnvironment returns the default control plane environment
func DefaultControlPlaneEnvironment() ControlPlaneEnvironment {
	return ControlPlaneEnvironment{
		Address:       ControlPlaneListenAddressDefault,
		Port:          ControlPlaneListenPortDefault,
		ManagerConfig: ManagedClusterConfigDefault,
		ExternalIP:    ControlPlaneExternalIPDefault,
		DockerPolicy:  ControlPlaneDockerPolicyDefault,
	}
}

// SpawnControlPlane spawns a single control plane container
func SpawnControlPlane(ctx context.Context, pullPolicy string, env ControlPlaneEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Get docker port for control plane
	port, err := nat.NewPort("tcp", env.Port)

	if err != nil {
		return err
	}

	// Set spawn information
	spawnInfo := NewSpawnInformation(pullPolicy, controlPlaneFullImage, controlPlaneContainerName, 1, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: controlPlaneImageName,
			Env: []string{
				ControlPlaneListenAddress + "=" + env.Address,
				ControlPlaneListenPort + "=" + env.Port,
				ManagedClusterConfig + "=" + env.ManagerConfig,
				ControlPlaneExternalIP + "=" + env.ExternalIP,
				ControlPlaneDockerPolicy + "=" + env.DockerPolicy,
			},
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: env.Port,
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
		}, &network.NetworkingConfig{}, controlPlaneContainerName)

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn control plane
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
