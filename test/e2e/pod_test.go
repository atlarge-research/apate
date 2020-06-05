package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

func TestSimplePodDeployment(t *testing.T) {
	if !enableDockerApatelets {
		t.Skip()
	}
	testSimplePodDeployment(t, env.Docker)
}

func TestSimplePodDeploymentRoutine(t *testing.T) {
	testSimplePodDeployment(t, env.Routine)
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

	err = kubectl.CreateWithNamespace([]byte(pods), kcfg, namespace)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	// TODO: Is this correct?
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	numpods, err := cluster.GetNumberOfPods(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 3, numpods)
}

func TestPodFailureDocker(t *testing.T) {
	if !enableDockerApatelets {
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
	simpleNodeDeployment(t, kcfg)
	time.Sleep(time.Second * 5)

	podFailure(t, kcfg)

	cancel()

	teardown(t)
}

func createConditionFunction(t *testing.T, numapatelets int, status corev1.ConditionStatus) func([]*corev1.Node) bool {
	return func(apatelets []*corev1.Node) bool {
		assert.Equal(t, numapatelets, len(apatelets))

		for _, apatelet := range apatelets {
			for _, c := range apatelet.Status.Conditions {
				if c.Type == corev1.NodeReady && c.Status == status {
					return true
				}
			}
		}

		return false
	}
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
	time.Sleep(time.Second * 60)

	// Get cluster object
	// TODO: Is this correct?
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	// Check if everything is ready
	ready, _ := getApateletWaitForCondition(t, cluster, 2, createConditionFunction(t, 2, corev1.ConditionTrue))

	assert.True(t, ready)

	// Deploy pods
	deployment := `
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
        containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
          - containerPort: 80
`
	err = kubectl.Create([]byte(deployment), kcfg)
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
