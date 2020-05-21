// Package container deals with the creation of docker containers
package container

import (
	"context"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

// SpawnCall is function which takes in the number of the container, then spawns the container
// and finally returns any error that might have occurred in the process
type SpawnCall func(int, context.Context) error

// SpawnInformation represents all required information to prepare the docker environment for spawning the containers
type SpawnInformation struct {
	pullPolicy           env.PullPolicy
	image, containerName string
	amount               int
	callback             SpawnCall
}

// NewSpawnInformation creates a new SpawnInformation struct using the given information
func NewSpawnInformation(pullPolicy env.PullPolicy, image, containerName string, amount int, callback SpawnCall) SpawnInformation {
	return SpawnInformation{
		pullPolicy:    pullPolicy,
		image:         image,
		containerName: containerName,
		amount:        amount,
		callback:      callback,
	}
}

// HandleSpawnContainers handles the preparation of docker images, removing of old containers and calling the given
// spawn call async
func HandleSpawnContainers(ctx context.Context, cli *client.Client, info SpawnInformation) error {
	// Prepare image
	err := prepareImage(ctx, cli, info.image, info.pullPolicy)
	if err != nil {
		return errors.Wrap(err, "failed to prepare image")
	}

	// Remove old containers
	err = removeOldContainers(ctx, cli, info.containerName)
	if err != nil {
		return errors.Wrap(err, "failed to remove old containers")
	}

	// Create error group to handle async spawning
	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < info.amount; i++ {
		i := i

		// Spawn container
		group.Go(func() error {
			return errors.Wrap(info.callback(i, ctx), "failed to run container callback")
		})
	}

	return errors.Wrap(group.Wait(), "error spawning containers")
}

func checkLocalImage(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return false, errors.Wrap(err, "failed to list images")
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}

	return false, nil
}

func removeOldContainers(ctx context.Context, cli *client.Client, name string) error {
	// Retrieve all old apatelet containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true,
		Filters: filters.NewArgs(filters.Arg("status", "exited"), filters.Arg("name", name))})

	if err != nil {
		return errors.Wrap(err, "failed to list exited containers")
	}

	// Remove old apatelet containers
	for _, cnt := range containers {
		err := cli.ContainerRemove(ctx, cnt.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true, RemoveLinks: false})

		if err != nil {
			return errors.Wrapf(err, "failed to remove old container %v", name)
		}
	}

	return nil
}

func prepareImage(ctx context.Context, cli *client.Client, imageName string, pullPolicy env.PullPolicy) error {
	switch pullPolicy {
	case env.AlwaysPull:
		return errors.Wrap(alwaysPull(ctx, cli, imageName), "failed to run alwaysPull to prepare image")
	case env.PullIfNotLocal:
		return errors.Wrap(pullIfNotLocal(ctx, cli, imageName), "failed to run pullIfNotLocal to prepare image")
	case env.AlwaysLocal:
		return errors.Wrap(alwaysCache(ctx, cli, imageName), "failed to run alwaysCache to prepare image")
	default:
		return errors.Errorf("unknown docker pull policy: %s", pullPolicy)
	}
}

func alwaysPull(ctx context.Context, cli *client.Client, imageName string) error {
	return errors.Wrapf(pullImage(ctx, cli, imageName), "failed to pull image %v", imageName)
}

func pullIfNotLocal(ctx context.Context, cli *client.Client, imageName string) error {
	// Check if the image is locally available
	localAvailable, err := checkLocalImage(ctx, cli, imageName)

	if err != nil {
		return errors.Wrap(err, "failed to check local image")
	}

	// If not, pull the image
	if !localAvailable {
		if err = pullImage(ctx, cli, imageName); err != nil {
			return errors.Wrap(err, "failed to pull image")
		}
	}

	return nil
}

func alwaysCache(ctx context.Context, cli *client.Client, imageName string) error {
	localAvailable, err := checkLocalImage(ctx, cli, imageName)

	if err != nil {
		return errors.Wrapf(err, "failed to check local image %v", imageName)
	}

	if !localAvailable {
		return errors.Errorf("image %v not available ", imageName)
	}

	return nil
}

func pullImage(ctx context.Context, cli *client.Client, imageName string) error {
	readCloser, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return errors.Wrapf(err, "failed to pull image %v", imageName)
	}

	return errors.Wrap(readCloser.Close(), "failed to close image pull reader")
}
