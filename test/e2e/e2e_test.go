package e2e

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	cmd "github.com/atlarge-research/opendc-emulate-kubernetes/cmd/apate/run"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

const enableDockerApatelets = false

// Please set the `$CI_PROJECT_DIR` to the root of the project
func setup(t *testing.T, kindClusterName string, runType env.RunType) {
	if testing.Short() {
		t.Skip("Skipping E2E")
	}

	if !enableDockerApatelets {
		log.Println("WARNING: Docker tests disabled!")
	}

	os.Args = []string{"apate-cp"}

	dir := os.Getenv("CI_PROJECT_DIR")
	if len(dir) == 0 {
		dir = "../../"
	}

	initEnv := env.ControlPlaneEnv()
	initEnv.PodCRDLocation = dir + "/config/crd/apate.opendc.org_podconfigurations.yaml"
	initEnv.NodeCRDLocation = dir + "/config/crd/apate.opendc.org_nodeconfigurations.yaml"
	initEnv.ManagerConfigLocation = dir + "/config/gitlab-kind.yml"
	initEnv.KinDClusterName = kindClusterName
	initEnv.ApateletRunType = runType
	initEnv.PrometheusStackEnabled = false
	env.SetEnv(initEnv)
}

func teardown(t *testing.T) {
	// #nosec
	_ = exec.Command("sh", "-c", "docker ps --filter name=apate --format \"{{.ID}}\" | xargs docker kill").Run()
	// #nosec
	// _ = exec.Command("sh", "-c", "docker ps -a --filter name=apate --format \"{{.ID}}\" | xargs docker rm").Run()

	// #nosec
	_ = exec.Command("docker", "kill", "apate-cp").Run()
	time.Sleep(time.Second * 5)

	err := os.Remove(env.ControlPlaneEnv().KubeConfigLocation)
	assert.NoError(t, err)
}

func waitForCP(t *testing.T) {
	cpEnv := env.DefaultControlPlaneEnvironment()
	statusClient, _ := controlplane.GetStatusClient(service.NewConnectionInfo(cpEnv.ListenAddress, cpEnv.ListenPort, false))
	ctx := context.Background()
	err := statusClient.WaitForControlPlane(ctx, time.Duration(5)*time.Minute)
	assert.NoError(t, err)
}

func getKubeConfig(t *testing.T) *kubeconfig.KubeConfig {
	args := []string{
		"apatectl",
		"kubeconfig",
	}

	// read kubeconfig
	c := capture()
	cmd.StartCmd(args)
	cfg := c.stop()
	println(cfg)

	kcfg, err := kubeconfig.FromBytes([]byte(cfg), os.TempDir()+"/apate-e2e-kubeconfig-"+uuid.New().String(), true)
	assert.NoError(t, err)

	return kcfg
}

func runScenario(t *testing.T) {
	args := []string{
		"apatectl",
		"run",
	}

	r, w, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdin = r
	go cmd.StartCmd(args)
	w.Write([]byte("\n"))
}
