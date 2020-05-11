package kubectl

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

const (
	prometheusNamespace = "apate-prometheus"
)

func prepareHelm() error {
	// helm repo add google https://kubernetes-charts.storage.googleapis.com/
	// helm repo update
	args := []string{
		"repo",
		"add",
		"google",
		"https://kubernetes-charts.storage.googleapis.com/",
	}

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("helm", args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	args = []string{
		"repo",
		"update",
	}

	// #nosec
	cmd = exec.Command("helm", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installPrometheus(kubecfg *kubeconfig.KubeConfig) error {
	if err := prepareHelm(); err != nil {
		return err
	}

	args := []string{
		"install",
		"prometheus-operator",
		"google/prometheus-operator",
	}

	args = append(args, "--namespace", prometheusNamespace)
	args = append(args, "--kubeconfig", kubecfg.Path)

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("helm", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreatePrometheusStack attempts to create the prometheus operator in the kubernetes cluster
// if the file cannot be found it will simply log an error.
func CreatePrometheusStack(kubecfg *kubeconfig.KubeConfig) {
	if err := createNameSpace(prometheusNamespace, kubecfg); err != nil {
		log.Printf("error while creating prometheus namespace: %v", err)
		return
	}

	time.Sleep(time.Second)
	err := installPrometheus(kubecfg)
	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
		return
	}
}
