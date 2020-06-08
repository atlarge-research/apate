package e2e

import (
	"context"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

func TestShutdownCPRoutines(t *testing.T) {
	rt := env.Routine
	setup(t, "TestShutdownCP_"+string(rt), rt)

	testShutdownCp(t)

	teardown(t)
}

func TestShutdownCPDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}

	rt := env.Docker
	setup(t, "TestShutdownCP_"+string(rt), rt)

	testShutdownCp(t)

	teardown(t)
}

func testShutdownCp(t *testing.T) {
	// By default, setup disables prometheus, but in this test it's enabled
	ctx := context.Background()
	stop := make(chan os.Signal, 1)

	// Start CP
	go cp.StartControlPlaneWithStopChannel(ctx, runner.New(), stop)

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	nodelst, err := cluster.GetNodes()

	// Assert the control plane is running
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nodelst.Items))
	assert.True(t, strings.Contains(nodelst.Items[0].Name, "control-plane"))

	stop <- syscall.SIGTERM

	// Wait for the controlplane to stop itself
	done := false
	for i := 0; i < 10; i++ {
		nodelst, err = cluster.GetNodes()

		if err != nil || len(nodelst.Items) == 0 {
			done = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	assert.True(t, done)
}
