// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

const (
	helmTemplatePathSuffix = "/prometheus.yml"
)

func createNameSpace(namespace string) error {
	args := []string{
		"create",
		"namespace",
		namespace,
	}

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("kubectl", args...)

	return cmd.Run()
}

// CreateWithNameSpace calls `kubectl create` with the given resourceConfig in the given namespace
// When this config is empty, it will not be called
func CreateWithNameSpace(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig, namespace string) error {
	if len(resourceConfig) > 0 {
		args := []string{
			"create",
			"-f",
			"-",
		}

		// specify config
		args = append(args, "--kubeconfig", kubeConfig.Path)

		// If namespace is non-null
		if namespace != "" {
			args = append(args, "--namespace", namespace)
		}

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

// Create calls `kubectl create` with the given resourceConfig
// When this config is empty, it will not be called
func Create(resourceConfig []byte, kubeConfig *kubeconfig.KubeConfig) error {
	return CreateWithNameSpace(resourceConfig, kubeConfig, "")
}

// CreatePrometheusStack attempts to create the prometheus operator in the kubernetes cluster
// if the file cannot be found it will simply log an error.
func CreatePrometheusStack(namespace string, kubeConfig *kubeconfig.KubeConfig) {
	bytes, err := ioutil.ReadFile(os.TempDir() + helmTemplatePathSuffix)
	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
	}

	if err := createNameSpace(namespace); err != nil {
		log.Printf("error while creating prometheus namespace: %v", err)
	}

	err = CreateWithNameSpace(bytes, kubeConfig, namespace)

	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
	}
}
