package kubectl

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

const (
	prometheusNamespace = "apate-prometheus"
)

func prepareHelm() error {
	// helm repo add google https://kubernetes-charts.storage.googleapis.com/
	// helm repo update
	const repo = "https://kubernetes-charts.storage.googleapis.com/"

	args := []string{
		"repo",
		"add",
		"google",
		repo,
	}

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("helm", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to add repository %v with helm %v", repo, strings.Join(args, " "))
	}

	args = []string{
		"repo",
		"update",
	}

	// #nosec
	cmd = exec.Command("helm", args...)

	return errors.Wrapf(cmd.Run(), "failed to update repositories with helm %v", strings.Join(args, " "))
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

	//// Values args
	//args = append(args, "--set", "nodeExporter.enabled=false")

	// Add settings
	args = append(args, "-f", env.ControlPlaneEnv().PrometheusConfigLocation)

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("helm", args...)

	if env.ControlPlaneEnv().DebugEnabled {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		println(cmd.String())
	}

	return errors.Wrapf(cmd.Run(), "failed to install Prometheus with helm %v", strings.Join(args, " "))
}

// CreatePrometheusStack attempts to create the prometheus operator in the kubernetes cluster
func CreatePrometheusStack(kubecfg *kubeconfig.KubeConfig) {
	log.Println("enabling prometheus stack")
	if err := CreateNameSpace(prometheusNamespace, kubecfg); err != nil {
		log.Printf("error while creating prometheus namespace: %v", err)
		return
	}

	err := installPrometheus(kubecfg)
	if err != nil {
		log.Printf("error while creating prometheus cluster: %v, prometheus stack not installed on the cluster\n", err)
		return
	}

	log.Println("enabled prometheus stack")
}
