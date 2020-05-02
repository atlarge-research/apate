// Package kind contains code to manage a KinD cluster and its configuration file
package kind

import (
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/kind/cmd/kind/app"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind"
)

// KinD is a struct implementing Manager for KinD clusters.
type KinD struct{}

// CreateCluster creates a new cluster with a given name.
func (KinD) CreateCluster(name string, kubeConfigLocation string, managerConfigPath string) error {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"create",
		"cluster",
	}

	args = append(args, "--name", name)
	args = append(args, "--kubeconfig", kubeConfigLocation)
	args = append(args, "--config", managerConfigPath)

	// Set up a cluster
	// Can't simply call Run as is done in Delete since we want to get error messages back.
	// Run just logs the error and returns.
	c := kind.NewCommand(logger, cmd.StandardIOStreams())
	c.SetArgs(args)
	if err := c.Execute(); err != nil {
		return err
	}

	// TODO: Move to cli command to get kubeconfig
	cfg, _ := ioutil.ReadFile("/tmp/apate/config")
	fmt.Print(string(cfg))

	// Update kube config to use internal
	err := useInternalKubeConfig(name, kubeConfigLocation)

	if err != nil {
		return err
	}

	// Only gets here after the cluster is running
	return nil
}

func useInternalKubeConfig(name string, kubeConfigLocation string) error {
	logger := cmd.NewLogger()

	args := []string{
		"get",
		"kubeconfig",
	}

	args = append(args, "--name", name)
	args = append(args, "--internal")

	cfg, err := os.Create(kubeConfigLocation)

	if err != nil {
		return err
	}

	c := kind.NewCommand(logger, cmd.IOStreams{
		In:     os.Stdin,
		Out:    cfg,
		ErrOut: os.Stderr,
	})

	c.SetArgs(args)
	if err := c.Execute(); err != nil {
		return err
	}

	return nil
}

// DeleteCluster deletes a cluster with a given name.
// This function never errors, even if the cluster didn't exist yet.
// Therefore it can be used to ensure no cluster with a certain name exists.
func (*KinD) DeleteCluster(name string) error {
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
	return app.Run(logger, cmd.StandardIOStreams(), args)
}
