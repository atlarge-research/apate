// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	helmTemplatePathSuffix = "/prometheus.yml"
	prometheusNamespace    = "apate-prom"
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
	resourceConfigPath := os.TempDir() + "/apate/res" //TODO: Fancy tmp dir
	//test, err := ioutil.TempFile("/tmp", "apate")

	_ = ioutil.WriteFile(resourceConfigPath, resourceConfig, 0o600)

	if len(resourceConfigPath) > 0 {
		args := []string{
			"create",
		}

		// specify config
		args = append(args, "-f", resourceConfigPath)
		args = append(args, "--kubeconfig", kubeConfig.Path)

		// If namespace is non-null
		if len(namespace) > 0 {
			args = append(args, "--namespace", namespace)
		}

		// #nosec as the arguments are controlled this is not a security problem
		cmd := exec.Command("kubectl", args...)
		//pipe, err := cmd.StdinPipe()
		//if err != nil {
		//	return err
		//}
		//
		//_, err = pipe.Write(resourceConfig)
		//if err != nil {
		//	return err
		//}
		//
		//if err := pipe.Close(); err != nil {
		//	return err
		//}

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

// CreatePrometheusStack attempts to create the prometheus operator in the kubernetes cluster
// if the file cannot be found it will simply log an error.
func CreatePrometheusStack(kubeConfig *kubeconfig.KubeConfig) {
	bytes, err := ioutil.ReadFile(os.TempDir() + helmTemplatePathSuffix)
	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
		return
	}

	if err := createNameSpace(prometheusNamespace, kubeConfig); err != nil {
		log.Printf("error while creating prometheus namespace: %v", err)
		return
	}

	time.Sleep(time.Second)
	err = Create(bytes, kubeConfig)

	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
		return
	}
}
