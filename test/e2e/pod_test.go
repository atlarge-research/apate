package e2e

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSimplePodDeployment(t *testing.T) {
	testSimplePodDeployment(t, env.Routine)
	testSimplePodDeployment(t, env.Docker)
}

func testSimplePodDeployment(t *testing.T, rt env.RunType) {
	setup(t, strings.ToLower("testSimplePodDeployment_"+string(rt)), rt)

	ctx, cancel := context.WithCancel(context.Background())

	go cp.StartControlPlane(ctx, runner.New())

	waitForCP(t)

	kcfg := getKubeConfig(t)

	// Setup some nodes
	simpleNodeDeployment(t, kcfg)
	time.Sleep(time.Second * 5)

	// Test pods
	simpleReplicaSet(t, kcfg)

	cancel()

	teardown(t)
}


func simpleReplicaSet(t *testing.T, kcfg *kubeconfig.KubeConfig) {
	pods := `
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # modify replicas according to your case
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: gcr.io/google_samples/gb-frontend:v3
`

	namespace := "simple-replica-set"

	err := kubectl.CreateNameSpace(namespace, kcfg)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	err = kubectl.CreateWithNameSpace([]byte(pods), kcfg, namespace)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	cluster, err := kubernetes.ClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	numpods, err := cluster.GetNumberOfPods(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 3, numpods)
}