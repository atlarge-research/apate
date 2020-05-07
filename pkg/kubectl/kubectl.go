// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"os/exec"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

// Create calls `kubectl create` with the given resourceConfig
// When this config is empty, it will not be called
func Create(resourceConfig []byte, kubeConfig kubeconfig.KubeConfig) error {
	if len(resourceConfig) > 0 {
		args := []string{
			"create",
			"-f",
			"-",
		}

		// specify config
		args = append(args, "--kubeconfig", kubeConfig.Path)

		// #nosec as the arguments are controlled this is not a security problem
		cmd := exec.Command("kubectl", args...)
		pipe, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		_, err = pipe.Write(resourceConfig)
		if err != nil {
			return err
		}

		if err := pipe.Close(); err != nil {
			return err
		}

		return cmd.Run()
	}

	return nil
}
