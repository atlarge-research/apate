package e2e

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	cmd "github.com/atlarge-research/opendc-emulate-kubernetes/cmd/apate/run"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// Please set the `$CI_PROJECT_DIR` to the root of the project
func setup(t *testing.T, kindClusterName string, runType env.RunType) {
	if testing.Short() {
		t.Skip("Skipping E2E")
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
	env.SetEnv(initEnv)
}

func teardown(t *testing.T) {
	// #nosec
	_ = exec.Command("docker", "kill", "apate-cp").Run()
	time.Sleep(time.Second * 5)

	err := os.Remove(env.ControlPlaneEnv().KubeConfigLocation)
	assert.NoError(t, err)
}

func TestSimplePodDeployment(t *testing.T) {
	testSimplePodDeployment(t, env.Routine)
	//testSimplePodDeployment(t, env.Docker)
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

func TestSimpleNodeDeployment(t *testing.T) {
	testSimpleNodeDeployment(t, env.Routine)
	//testSimpleNodeDeployment(t, env.Docker)
}

// To run this, make sure ./config/kind.yml is put in the right directory (/tmp/apate/manager)
// or the env var CP_MANAGER_CONFIG_LOCATION point to it
func testSimpleNodeDeployment(t *testing.T, rt env.RunType) {
	setup(t, "TestSimpleNodeDeployment_"+string(rt), rt)

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
	time.Sleep(time.Second * 15)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
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

	err = kubectl.ExecuteWithNamespace("create", []byte(pods), kcfg, namespace)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	numpods, err := cluster.GetNumberOfPods(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 3, numpods)
}
