package kubernetes

import (
	"os"
	"os/exec"
	"strings"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/pkg/errors"

	"sigs.k8s.io/kind/cmd/kind/app"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind"
)

// KinDClusterManager is a struct which implements ClusterManager by creating a kind cluster.
type KinDClusterManager struct{}

// GetKubeConfig validates the cluster name, creates the Kind cluster and returns its kubeconfig.
func (k *KinDClusterManager) GetKubeConfig() (*kubeconfig.KubeConfig, error) {
	if env.ControlPlaneEnv().KinDClusterName == "" {
		return nil, errors.New("trying to create a KinDClusterManager cluster with an empty name")
	}

	err := k.deleteCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failure to delete kind cluster to prepare for creating a new one")
	}

	if err = k.prepareCluster(); err != nil {
		return nil, errors.Wrapf(err, "failed to start kind cluster. Possible cleanup error: %v", k.deleteCluster())
	}

	err = k.writeKubeConfig()
	if err != nil {
		err = errors.Wrap(err, "failed to write kubeconfig")
		return nil, errors.Wrapf(k.deleteCluster(), "failed to delete kind cluster to clean up earlier failure (%v)", err)
	}

	kubeConfig, err := kubeconfig.FromPath(env.ControlPlaneEnv().KubeConfigLocation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read kubeconfig")
	}

	return kubeConfig, nil
}

// Shutdown deletes the KinD cluster.
func (k *KinDClusterManager) Shutdown(*Cluster) error {
	return errors.Wrap(k.deleteCluster(), "failed to delete kind cluster")
}

func (k *KinDClusterManager) writeKubeConfig() error {
	kubeConfigLocation := env.ControlPlaneEnv().KubeConfigLocation
	if env.ControlPlaneEnv().UseDockerHostname {
		// Replace any https address by the "docker" hostname.
		// This is used in CI, where the control plane had a 172.17.0.0/16 address, and the KinDClusterManager cluster a 172.18.0.0/16 address
		// which was only reachable using "docker"s as hostname.
		// #nosec
		cmdSed := exec.Command("sed", "-i", "-r", "s/https:\\/\\/(.+):/https:\\/\\/docker:/g", kubeConfigLocation)
		return errors.Wrap(cmdSed.Run(), "failed to apply sed to the kube config")
	}

	// Update kube config to use internal
	return errors.Wrapf(k.useInternalKubeConfig(env.ControlPlaneEnv().KinDClusterName, kubeConfigLocation), "failed to use internal Kubeconfig")
}

func (k *KinDClusterManager) prepareCluster() error {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"create",
		"cluster",
	}

	args = append(args, "--name", strings.ToLower(env.ControlPlaneEnv().KinDClusterName))
	args = append(args, "--kubeconfig", env.ControlPlaneEnv().KubeConfigLocation)
	args = append(args, "--config", env.ControlPlaneEnv().ManagerConfigLocation)

	if env.ControlPlaneEnv().DebugEnabled {
		// This number is also used in KinD internally to denote trace.
		args = append(args, "--verbosity=2147483647")
	}

	// Set up a cluster
	// Can't simply call Run as is done in Delete since we want to get error messages back.
	// Run just logs the error and returns.
	c := kind.NewCommand(logger, cmd.StandardIOStreams())
	c.SetArgs(args)
	return errors.Wrapf(c.Execute(), "failed to create kind cluster with kind %v", strings.Join(args, " "))
}

func (k *KinDClusterManager) useInternalKubeConfig(name string, kubeConfigLocation string) error {
	logger := cmd.NewLogger()

	args := []string{
		"get",
		"kubeconfig",
	}

	args = append(args, "--name", strings.ToLower(name))
	args = append(args, "--internal")

	cfg, err := os.Create(kubeConfigLocation)
	if err != nil {
		return errors.Wrapf(err, "failed to create file for Kubeconfig at %v", kubeConfigLocation)
	}

	c := kind.NewCommand(logger, cmd.IOStreams{
		In:     os.Stdin,
		Out:    cfg,
		ErrOut: os.Stderr,
	})

	c.SetArgs(args)
	return errors.Wrapf(c.Execute(), "failed run kind %v to retrieve internal Kubeconfig", strings.Join(args, " "))
}

func (k *KinDClusterManager) deleteCluster() error {
	// TODO: use our own/a global logger?
	logger := cmd.NewLogger()

	args := []string{
		"delete",
		"cluster",
	}

	args = append(args, "--name", strings.ToLower(env.ControlPlaneEnv().KinDClusterName))

	// Deletes the cluster
	// As far as I could test this call never errors (it just doesn't do anything
	// when the cluster doesn't exist) so I don't think the system used in CreateCluster is necessary.
	return errors.Wrapf(app.Run(logger, cmd.StandardIOStreams(), args), "failed to run kind %v to delete cluster", strings.Join(args, " "))
}
