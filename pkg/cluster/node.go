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
func SpawnNodes(ctx context.Context, amountOfNodes int) error {
	var err error

	imageName := "apatekubernetes/apatelet:latest"
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	localAvailable, err := checkLocal(ctx, cli, imageName)

	if err != nil {
		return err
	}

	if !localAvailable {
		if err = pullImage(ctx, cli, imageName); err != nil {
			return err
		}
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

//TODO: Optimise
func checkLocal(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
	imgs, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return false, err
	}

	for _, img := range imgs {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}

	return false, nil
}

func pullImage(ctx context.Context, cli *client.Client, imageName string) error {
	if _, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{}); err != nil {
		return err
	}

	return nil
}

func spawnNode(ctx context.Context, cli *client.Client, hostname string, nodeIndex int) error {
	// TODO check if exists
	c, err := cli.ContainerCreate(ctx, &container.Config{
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
