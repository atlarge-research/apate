// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// CreateNameSpace creates a namespace on the cluster
func CreateNameSpace(namespace string, kubeConfig *kubeconfig.KubeConfig) error {
	args := []string{
		"create",
		"namespace",
		namespace,
	}

	// specify config
	args = append(args, "--kubeconfig", kubeConfig.Path)

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("kubectl", args...)

	return errors.Wrapf(cmd.Run(), "failed to create namespace with kubectl %v", strings.Join(args, " "))
}

// ExecuteWithNamespace calls `kubectl <command>` with the given resourceConfig in the given namespace
// When this config is empty, it will not be called
func ExecuteWithNamespace(command string, resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig, namespace string) error {
	if len(resourceConfig) > 0 {
		cfgFile, err := ioutil.TempFile("", "apate-kubectl-")
		if err != nil {
			return errors.Wrap(err, "failed to create tempfile for Kubeconfig")
		}
		_, err = cfgFile.Write(resourceConfig)
		if err != nil {
			return errors.Wrap(err, "failed to write Kubeconfig to tempfile")
		}
		defer func() {
			err = os.Remove(cfgFile.Name())
			if err != nil {
				// Unable to remove temp file, doesn't matter that much but logging anyway
				log.Printf("unable to delete temporary file: %v\n", err)
			}
		}()

		args := []string{
			command,
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
		return errors.Wrapf(cmd.Run(), "failed to create resource %v in namespace %v with kubectl %v", cfgFile.Name(), namespace, strings.Join(args, " "))
	}

	return nil
}

// Create calls `kubectl create` with the given resourceConfig
// When this config is empty, it will not be called
func Create(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return errors.Wrapf(ExecuteWithNamespace("create", resourceConfig, kubeConfig, ""), "failed to create resource in default namespace")
}

// Apply calls `kubectl apply` with the given resourceConfig
// When this config is empty, it will not be called
func Apply(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return errors.Wrapf(ExecuteWithNamespace("apply", resourceConfig, kubeConfig, ""), "failed to create resource in default namespace")
}

// Delete calls `kubectl delete` with the given resourceConfig
// When this config is empty, it will not be called
func Delete(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return errors.Wrapf(ExecuteWithNamespace("delete", resourceConfig, kubeConfig, ""), "failed to create resource in default namespace")
}
