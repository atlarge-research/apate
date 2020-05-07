// Package kubectl provides functions to interact with the kubectl binary
package kubectl

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
)

const resourceConfigPathSuffix = "/apate/resourceConfig.yml"

// Create calls `kubectl create` with the saved resourceconfig in [SaveResourceConfig]
func Create(cfg kubeconfig.KubeConfig) error {
	args := []string{
		"create",
		"-f",
		os.TempDir() + resourceConfigPathSuffix,
	}

	// specify config
	args = append(args, "--kubeconfig", cfg.Path)

	// #nosec as the arguments are controlled this is not a security problem
	cmd := exec.Command("kubectl", args...)

	return cmd.Run()
}

// SaveResourceConfig saves the bytes of a resource.yml to an (internal) file.
func SaveResourceConfig(bytes []byte) error {
	return ioutil.WriteFile(os.TempDir()+resourceConfigPathSuffix, bytes, 0o600)
}
