package cluster

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

const (
	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress = "CP_LISTEN_ADDRESS"

	// ControlPlanePort is the port of the control plane
	ControlPlanePort = "CP_LISTEN_PORT"
)

// SpawnNodes spawns multiple Apatelet Docker containers
func SpawnNodes(ctx context.Context, amountOfNodes int, info *service.ConnectionInfo) error {
	var err error

	// Get docker cli
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Prepare image
	err = prepareImage(ctx, cli)
	if err != nil {
		return err
	}

	// Create error group to handle async spawning
	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < amountOfNodes; i++ {
		i := i

		// Spawn container
		group.Go(func() error {
			if err := spawnNode(ctx, cli, info, i); err != nil {
				return err
			}

			return nil
		})
	}

	return group.Wait()
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

func prepareImage(ctx context.Context, cli *client.Client) error {
	imageName := "apatekubernetes/apatelet:latest"

	// Check if the image is already available
	localAvailable, err := checkLocal(ctx, cli, imageName)

	if err != nil {
		return err
	}

	// If not, pull the image
	if !localAvailable {
		if err = pullImage(ctx, cli, imageName); err != nil {
			return err
		}
	}

	return nil
}

func pullImage(ctx context.Context, cli *client.Client, imageName string) error {
	if _, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{}); err != nil {
		return err
	}

	return nil
}

func spawnNode(ctx context.Context, cli *client.Client, info *service.ConnectionInfo, nodeIndex int) error {
	c, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "apatelet:latest",
		Env: []string{
			ControlPlaneAddress + "=" + info.Address,
			ControlPlanePort + "=" + strconv.Itoa(info.Port),
		},
	}, nil, nil, "apatelet-"+strconv.Itoa(nodeIndex))

	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}
