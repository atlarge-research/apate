package cluster

import (
	"sigs.k8s.io/kind/cmd/kind/app"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind"
)

// Creates a new cluster with a given name.
func CreateCluster(name string) error {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"create",
		"cluster",
	}

	args = append(args, "--name", name)

	// TODO: Avoid kind overwriting the global config file.
	// Set up a cluster

	c := kind.NewCommand(logger, cmd.StandardIOStreams())
	c.SetArgs(args)
	if err := c.Execute(); err != nil {
		return err
	}

	// Only gets here after the cluster is running
	return nil
}

// Deletes a cluster with a given name.
// This function never errors, even if the cluster didn't exist yet.
// Therefore it can be used to ensure no cluster with a certain name exists.
func DeleteCluster(name string) {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"delete",
		"cluster",
	}

	args = append(args, "--name", name)


	// Deletes the cluster
	app.Run(logger, cmd.StandardIOStreams(), args)
	// Only gets here after the cluster is deleted
}