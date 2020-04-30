package cluster

import (
	"context"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// SpawnNodes spawns multiple Apatelet Docker containers
func SpawnNodes(amountOfNodes int) error {
	var err error

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	if err = pullImage(ctx, cli); err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// TODO async
	for i := 0; i < amountOfNodes; i++ {
		if err := spawnNode(ctx, cli, hostname, i); err != nil {
			return err
		}
	}

	return nil
}

func pullImage(ctx context.Context, cli *client.Client) error {
	// TODO: Get from env/config
	imageName := "apatekubernetes/apatelet:latest"
	if _, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{}); err != nil {
		return err
	}

	return nil
}

func spawnNode(ctx context.Context, cli *client.Client, hostname string, nodeIndex int) error {
	// TODO: check if exists
	c, err := cli.ContainerCreate(ctx, &container.Config{
		// TODO: Should match the imageName from pullImage
		Image: "apatelet:latest",
		Env: []string{
			"CP_HOSTNAME=" + hostname,
		},
	}, nil, nil, "apatelet-"+strconv.Itoa(nodeIndex))

	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}
