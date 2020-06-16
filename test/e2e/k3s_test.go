package e2e

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

// K3s DEPLOYMENT
// This test will only work when the K3S_KUBE_CONFIG is set
func TestK3sRoutine(t *testing.T) {
	testK3s(t, env.Routine)
}

func TestK3sDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testK3s(t, env.Docker)
}

func testK3s(t *testing.T, rt env.RunType) {
	k3sKubeConfigEnv, ok := os.LookupEnv("K3S_KUBE_CONFIG")
	if !ok {
		t.Skip("WARNING: skipping k3s test due to missing env variable `K3S_KUBE_CONFIG`")
	}

	setup(t, "TestK3s"+string(rt), rt)

	bytes, err := ioutil.ReadFile(filepath.Clean(k3sKubeConfigEnv))
	if err != nil {
		t.Skip(errors.Wrap(err, "reading k3s config failed"))
	}

	cpEnv := env.ControlPlaneEnv()
	cpEnv.KubeConfig = string(bytes)
	env.SetEnv(cpEnv)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.StartControlPlane(ctx, runner.New())

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	// Test simple deployment
	simpleNodeDeployment(t, kcfg)

	cancel()

	teardown(t)
}
