package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/api/types/filters"

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

	// ControlPlaneDockerPolicy specifies the docker pull policy for apatelet images
	ControlPlaneDockerPolicy = "CP_DOCKER_POLICY"
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

	// Remove old containers
	err = removeOldContainers(ctx, cli)
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

func checkLocalImage(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
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

func removeOldContainers(ctx context.Context, cli *client.Client) error {
	// Retrieve all old apatelet containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true,
		Filters: filters.NewArgs(filters.Arg("status", "exited"), filters.Arg("name", "apatelet"))})

	if err != nil {
		return nil
	}

	// Remove old apatelet containers
	for _, cnt := range containers {
		err := cli.ContainerRemove(ctx, cnt.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true, RemoveLinks: false})

		if err != nil {
			return err
		}
	}

	return nil
}

func prepareImage(ctx context.Context, cli *client.Client) error {
	imageName := "apatekubernetes/apatelet:latest"

	// Retrieve pull policy
	var policy string
	if val, ok := os.LookupEnv(ControlPlaneDockerPolicy); ok {
		policy = val
	} else {
		policy = "pull-not-local"
	}

	switch policy {
	case "pull-always":
		return alwaysPull(ctx, cli, imageName)
	case "pull-not-local":
		return pullIfNotLocal(ctx, cli, imageName)
	case "cache-always":
		return alwaysCache(ctx, cli, imageName)
	default:
		return fmt.Errorf("unknown docker pull policy: %s", policy)
	}
}

func alwaysPull(ctx context.Context, cli *client.Client, imageName string) error {
	return pullImage(ctx, cli, imageName)
}

func pullIfNotLocal(ctx context.Context, cli *client.Client, imageName string) error {
	// Get version from docker hub
	localAvailable, err := checkLocalImage(ctx, cli, imageName)

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

func alwaysCache(ctx context.Context, cli *client.Client, imageName string) error {
	localAvailable, err := checkLocalImage(ctx, cli, imageName)

	if err != nil {
		return err
	}

	if !localAvailable {
		return errors.New("image %s not available ")
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
