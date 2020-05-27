package e2e

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/app"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	cmd "github.com/atlarge-research/opendc-emulate-kubernetes/cmd/apate/app"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

func TestSimplePodDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E")
	}

	err := os.Setenv("KIND_CLUSTER_NAME", "TestSimplePodDeployment")
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go cp.Main(ctx)

	waitForCP(t)

	kcfg := getKubeConfig(t)

	// Setup some nodes
	simpleNodeDeployment(t, kcfg)
	time.Sleep(time.Second * 5)

	// Test pods
	simpleReplicaSet(t, kcfg)

	cancel()

	// #nosec
	_ = exec.Command("docker", "kill", "apate-cp").Run()
	time.Sleep(time.Second * 5)
}

// To run this, make sure ./config/kind is put in the right directory (/tmp/apate/manager)
// or the env var CP_K8S_CONFIG point to it
func TestSimpleNodeDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E")
	}

	err := os.Setenv("KIND_CLUSTER_NAME", "TestSimpleNodeDeployment")
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.Main(ctx)

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	// Test simple deployment
	simpleNodeDeployment(t, kcfg)

	cancel()

	// #nosec
	_ = exec.Command("docker", "kill", "apate-cp").Run()
	time.Sleep(time.Second * 5)
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
	cmd.Main(args)
	cfg := c.stop()
	println(cfg)

	kcfg, err := kubeconfig.FromBytes([]byte(cfg), os.TempDir()+"/apate-e2e-kubeconfig-"+uuid.New().String())
	assert.NoError(t, err)

	return kcfg
}

func simpleNodeDeployment(t *testing.T, kcfg *kubeconfig.KubeConfig) {
	rc := `
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: e2e-deployment
spec:
    replicas: 2
    resources:
        memory: 5G
        cpu: 1000
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
`

	err := kubectl.Create([]byte(rc), kcfg)
	assert.NoError(t, err)
	time.Sleep(time.Second)

	cluster, err := kubernetes.ClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 3, nodes)
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
