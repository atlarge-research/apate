package kubectl

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"strings"

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
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to run helm %v", strings.Join(args[:], " "))
	}

	args = []string{
		"repo",
		"update",
	}

	// #nosec
	cmd = exec.Command("helm", args...)

	return errors.Wrapf(cmd.Run(), "failed to run helm %v", strings.Join(args[:], " "))
}

func installPrometheus(kubecfg *kubeconfig.KubeConfig) error {
	if err := prepareHelm(); err != nil {
		return errors.Wrap(err, "failed to prepare Helm")
	}

	args := []string{
		"install",
		"prometheus-operator",
		"google/prometheus-operator",
	}

	// Basic args
	args = append(args, "--namespace", prometheusNamespace)
	args = append(args, "--kubeconfig", kubecfg.Path)

	// Values args
	args = append(args, "--set", "nodeExporter.enabled=false")

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("helm", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return errors.Wrapf(cmd.Run(), "failed to run helm %v", strings.Join(args[:], " "))

}

// CreatePrometheusStack attempts to create the prometheus operator in the kubernetes cluster
func CreatePrometheusStack(kubecfg *kubeconfig.KubeConfig) {
	log.Println("enabling prometheus stack")
	if err := createNameSpace(prometheusNamespace, kubecfg); err != nil {
		log.Printf("error while creating prometheus namespace: %+v", err)
		return
	}

	err := installPrometheus(kubecfg)
	if err != nil {
		log.Printf("error while creating prometheus cluster: %+v, prometheus stack not installed on the cluster\n", err)
		return
	}

	log.Println("enabled prometheus stack")
}
