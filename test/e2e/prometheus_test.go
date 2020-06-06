package e2e

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

func TestRunPrometheusRoutines(t *testing.T) {
	rt := env.Routine
	setup(t, "TestRunPrometheus_"+string(rt), rt)

	testRunPrometheus(t)

	teardown(t)
}

func TestRunPrometheusDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}

	rt := env.Docker
	setup(t, "TestRunPrometheus_"+string(rt), rt)

	testRunPrometheus(t)

	teardown(t)
}

func testRunPrometheus(t *testing.T) {
	// By default, setup disables prometheus, but in this test it's enabled
	e := env.ControlPlaneEnv()
	e.PrometheusStackEnabled = true
	env.SetEnv(e)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.StartControlPlane(ctx, runner.New())

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	// Wait for prometheus to have fully started
	allrunning := false
	for i := 0; i < 10; i++ {
		pods, err := cluster.GetPods("apate-prometheus")
		assert.NoError(t, err)

		// Prometheus should spawn 5 pods
		numpods := 5
		log.Printf("Testing if all %v pods are running", numpods)

		numrunning := 0
		for _, pod := range pods.Items {
			assert.True(t, strings.Contains(pod.Name, "prometheus"))
			if corev1.PodRunning == pod.Status.Phase {
				numrunning++
			}
		}

		if numpods == numrunning {
			allrunning = true
			println("all pods are running")
			break
		}

		println("retrying in 30 seconds")
		time.Sleep(30 * time.Second)
	}

	assert.True(t, allrunning)

	cancel()
}
