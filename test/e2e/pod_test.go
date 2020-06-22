package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-units"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/apate/internal/kubectl"
	"github.com/atlarge-research/apate/pkg/env"
	"github.com/atlarge-research/apate/pkg/kubernetes"
	"github.com/atlarge-research/apate/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/apate/pkg/runner"
	cp "github.com/atlarge-research/apate/services/controlplane/run"
)

// POD DEPLOYMENT
func TestSimplePodDeployment(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testSimplePodDeployment(t, env.Docker)
}

func TestSimplePodDeploymentRoutine(t *testing.T) {
	testSimplePodDeployment(t, env.Routine)
}

func testSimplePodDeployment(t *testing.T, rt env.RunType) {
	setup(t, strings.ToLower("TestSimplePodDeployment"+string(rt)), rt)

	ctx, cancel := context.WithCancel(context.Background())

	go cp.StartControlPlane(ctx, runner.New())

	waitForCP(t)

	kcfg := getKubeConfig(t)

	// Setup some nodes
	simpleNodeDeployment(t, kcfg, 2)
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
      nodeSelector:
        emulated: "yes"
      tolerations:
        -   key: emulated
            operator: Exists
            effect: NoSchedule
      containers:
      - name: php-redis
        image: gcr.io/google_samples/gb-frontend:v3
`

	namespace := "simple-replica-set"

	err := kubectl.CreateNameSpace(namespace, kcfg)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	err = kubectl.CreateWithNamespace([]byte(pods), kcfg, namespace)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	numpods, err := cluster.GetNumberOfPods(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 3, numpods)

	running := false
	for i := 0; i < 10; i++ {
		podlist, err := cluster.GetPods(namespace)
		assert.NoError(t, err)
		if running = arePodsAreRunning(podlist); running {
			break
		}

		time.Sleep(time.Second * 30)
	}

	assert.True(t, running)
}

var nginx = `
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx-deployment3
    labels:
      app: nginx
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: nginx
    template:
      metadata:
        labels:
          app: nginx
          apate: test-pod1
      spec:
        nodeSelector:
          emulated: "yes"
        tolerations:
          -   key: emulated
              operator: Exists
              effect: NoSchedule
        containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
          - containerPort: 80
`

// POD FAILURE
func TestPodFailureDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testPodFailure(t, env.Docker)
}

func TestPodFailureRoutine(t *testing.T) {
	testPodFailure(t, env.Routine)
}

func testPodFailure(t *testing.T, rt env.RunType) {
	setup(t, strings.ToLower("testPodFailure"+string(rt)), rt)

	ctx, cancel := context.WithCancel(context.Background())

	go cp.StartControlPlane(ctx, runner.New())

	waitForCP(t)

	kcfg := getKubeConfig(t)

	// test
	simpleNodeDeployment(t, kcfg, 2)
	time.Sleep(time.Second * 5)

	podFailure(t, kcfg)

	cancel()

	teardown(t)
}

func podFailure(t *testing.T, kcfg *kubeconfig.KubeConfig) {
	pcfg := `
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
  name: test-pod1
spec:
  pod_status: FAILED
`

	// Create pod CRDs
	err := kubectl.Create([]byte(pcfg), kcfg)
	assert.NoError(t, err)
	time.Sleep(longTimeout)

	// Get cluster object
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	// Check if everything is ready
	ready, _ := getApateletWaitForCondition(t, cluster, 2, createApateletConditionFunction(t, 2, corev1.ConditionTrue))

	assert.True(t, ready)

	// Deploy pods
	err = kubectl.Create([]byte(nginx), kcfg)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	// assert state
	pods, err := cluster.GetPods("default")
	assert.NoError(t, err)

	numFailed := 0

	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, "nginx-deployment3") {
			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionFalse {
					numFailed++
				}
			}
		}
	}

	assert.GreaterOrEqual(t, 1, numFailed)
}

// POD RESOURCE UPDATE
func TestPodResourceDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testPodResource(t, env.Docker)
}

func TestPodResourceRoutine(t *testing.T) {
	testPodResource(t, env.Routine)
}

func testPodResource(t *testing.T, rt env.RunType) {
	setup(t, strings.ToLower("testPodResourceUpdate"+string(rt)), rt)

	ctx, cancel := context.WithCancel(context.Background())

	go cp.StartControlPlane(ctx, runner.New())

	waitForCP(t)

	kcfg := getKubeConfig(t)

	// test
	simpleNodeDeployment(t, kcfg, 1)
	time.Sleep(time.Second * 5)

	podResource(t, kcfg)

	cancel()

	teardown(t)
}

func podResource(t *testing.T, kcfg *kubeconfig.KubeConfig) {
	pcfg := `
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
  name: test-pod1
spec:
  pod_resources:
    memory: 2G
    cpu: 10
`

	// Create pod CRDs
	err := kubectl.Create([]byte(pcfg), kcfg)
	assert.NoError(t, err)
	time.Sleep(longTimeout)

	// Get cluster object
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	// Check if everything is ready
	ready, _ := getApateletWaitForCondition(t, cluster, 1, createApateletConditionFunction(t, 1, corev1.ConditionTrue))

	assert.True(t, ready)

	// Deploy pods
	err = kubectl.Create([]byte(nginx), kcfg)
	assert.NoError(t, err)
	time.Sleep(longTimeout)

	time.Sleep(30 * time.Second) // Once every 30 seconds an update in status is scheduled

	// assert state
	nodes, err := cluster.GetNodes()
	assert.NoError(t, err)

	for _, node := range nodes.Items {
		if strings.HasPrefix(node.Name, "apatelet-") {
			allocatable := node.Status.Allocatable
			assert.Equal(t, int64(990), allocatable.Cpu().Value())
			assert.Equal(t, int64(5*units.GiB-2*units.GiB), allocatable.Memory().Value())
			assert.Equal(t, int64(149), allocatable.Pods().Value())
			assert.Equal(t, int64(120*units.GiB), allocatable.StorageEphemeral().Value())

			storage := allocatable[corev1.ResourceStorage]
			assert.Equal(t, int64(5*units.TiB), storage.Value())
		}
	}
}
