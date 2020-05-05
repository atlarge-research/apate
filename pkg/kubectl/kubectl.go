package kubectl

import (
	_ "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"io/ioutil"
	"os"
	"os/exec"
)

const resourceConfigPathSuffix = "/apate/resourceConfig.yml"

func Create(cfg kubeconfig.KubeConfig) error {
	args := []string{
		"create",
		"-f",
		os.TempDir() + resourceConfigPathSuffix,
	}

	// specify config
	args = append(args, "--kubeconfig", cfg.Path)

	cmd := exec.Command("kubectl", args...)

	return cmd.Run()

}

func SaveResourceConfig(bytes []byte) error {
	return ioutil.WriteFile(os.TempDir() + resourceConfigPathSuffix, bytes, 0o600)
}