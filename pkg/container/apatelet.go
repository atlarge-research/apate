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

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	Address, Port string
}

// DefaultApateEnvironment returns the default apate environment
func DefaultApateEnvironment() ApateletEnvironment {
	return ApateletEnvironment{
		Address: ec.ApateletListenAddressDefault,
		Port:    ec.ApateletListenPortDefault,
	}
}

// SpawnApatelets spawns multiple Apatelet Docker containers
func SpawnApatelets(ctx context.Context, amountOfNodes int, info *service.ConnectionInfo, pullPolicy string, env ApateletEnvironment) error {
	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Get docker port for apatelet
	port, err := nat.NewPort("tcp", env.Port)

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
					ec.ControlPlaneAddress + "=" + info.Address,
					ec.ControlPlanePort + "=" + strconv.Itoa(info.Port),
					ec.ApateletListenAddress + "=" + env.Address,
					ec.ApateletListenPort + "=" + env.Port,
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
			port1 := 7000 + i
			if port1 == 8085 {
				port1 = 6999
			}

			port2 := 12000 + i
			port3 := 17000 + i

			go run.StartApatelet(info.Address, info.Port, env.Address, port1, port2, port3)
			return nil
		}
	})

	// Spawn apatelets
	return HandleSpawnContainers(ctx, cli, spawnInfo)
}
