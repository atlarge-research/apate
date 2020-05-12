// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"os"
	"os/exec"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

// Create calls `kubectl create` with the given resourceConfig
// When this config is empty, it will not be called
func Create(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return call("create", resourceConfig, kubeConfig)
}

func call(command string, resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	if len(resourceConfig) > 0 {
		args := []string{
			command,
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

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := pipe.Close(); err != nil {
			return err
		}

		return cmd.Run()
	}

	return nil
}
