package clustercreationinterface

import (
	"sigs.k8s.io/kind/cmd/kind/app"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind"
)

type Kind struct {}

// Creates a new cluster with a given name.
func (Kind) CreateCluster(name string, kubeconfiglocation string) error {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"create",
		"cluster",
	}

	args = append(args, "--name", name)
	args = append(args, "--kubeconfig", kubeconfiglocation)

	// TODO: Avoid kind overwriting the global config file.
	// Set up a cluster
	// Can't simply call Run as is done in Delete since we want to get error messages back.
	// Run just logs the error and returns.
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
func (*Kind) DeleteCluster(name string) {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"delete",
		"cluster",
	}

	args = append(args, "--name", name)

	// Deletes the cluster
	// As far as I could test this call never errors (it just doesn't do anything
	// when the cluster doesn't exist) so I don't think the system used in CreateCluster is necessary.
	app.Run(logger, cmd.StandardIOStreams(), args)
	// Only gets here after the cluster is deleted
}

// Returns the name of a context for kubernetes to use for a given cluster name.
func (*Kind) ClusterContext(name string) string {
	return "kind-" + name
}
