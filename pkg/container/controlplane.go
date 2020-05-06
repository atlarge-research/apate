package container

import (
	"context"
	ec "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

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
		Address:       ec.ControlPlaneListenAddressDefault,
		Port:          ec.ControlPlaneListenPortDefault,
		ManagerConfig: ec.ManagedClusterConfigDefault,
		ExternalIP:    ec.ControlPlaneExternalIPDefault,
		DockerPolicy:  ec.ControlPlaneDockerPolicyDefault,
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
	spawnInfo := NewSpawnInformation(pullPolicy, ec.ControlPlaneFullImage, ec.ControlPlaneContainerName, 1, func(i int, ctx context.Context) error {
		c, err := cli.ContainerCreate(ctx, &container.Config{
			Image: ec.ControlPlaneImageName,
			Env: []string{
				ec.ControlPlaneListenAddress + "=" + env.Address,
				ec.ControlPlaneListenPort + "=" + env.Port,
				ec.ManagedClusterConfig + "=" + env.ManagerConfig,
				ec.ControlPlaneExternalIP + "=" + env.ExternalIP,
				ec.ControlPlaneDockerPolicy + "=" + env.DockerPolicy,
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
		}, &network.NetworkingConfig{}, ec.ControlPlaneContainerName)

		if err != nil {
			return err
		}

		return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	})

	// Spawn control plane
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
