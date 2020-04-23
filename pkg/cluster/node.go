package cluster

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnNode spawns a Virtual-Kubelet Docker container
func SpawnNode(nodeIndex int) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	ctx := context.Background()
	c, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "virtual_kubelet:latest",
		Env: []string{
			"CP_HOSTNAME=" + hostname,
		},
	}, nil, nil, "virtual_kubelet-"+string(nodeIndex))

	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}
