// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

func createNameSpace(namespace string, kubeConfig *kubeconfig.KubeConfig) error {
	args := []string{
		"create",
		"namespace",
		namespace,
	}

	// specify config
	args = append(args, "--kubeconfig", kubeConfig.Path)

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("kubectl", args...)
	return cmd.Run()
}

// CreateWithNameSpace calls `kubectl create` with the given resourceConfig in the given namespace
// When this config is empty, it will not be called
func CreateWithNameSpace(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig, namespace string) error {
	if len(resourceConfig) > 0 {
		cfgFile, err := ioutil.TempFile("", "apate-")
		if err != nil {
			return err
		}
		_, err = cfgFile.Write(resourceConfig)
		if err != nil {
			return err
		}
		defer func() {
			err = os.Remove(cfgFile.Name())
			if err != nil {
				// Unable to remove temp file, doesn't matter that much but logging anyway
				log.Printf("unable to delete temporary file: %v\n", err)
			}
		}()

		args := []string{
			"create",
		}

		// specify config
		args = append(args, "-f", cfgFile.Name())
		args = append(args, "--kubeconfig", kubeConfig.Path)

		// If namespace is non-null
		if len(namespace) > 0 {
			args = append(args, "--namespace", namespace)
		}

		// #nosec as the arguments are controlled this is not a security problem
		cmd := exec.Command("kubectl", args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

// Create calls `kubectl create` with the given resourceConfig
// When this config is empty, it will not be called
func Create(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return CreateWithNameSpace(resourceConfig, kubeConfig, "")
}
