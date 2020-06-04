package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	err = kubectl.CreateWithNameSpace([]byte(pods), kcfg, namespace)
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)

	cluster, err := kubernetes.ClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	numpods, err := cluster.GetNumberOfPods(namespace)
	assert.NoError(t, err)
	assert.Equal(t, 3, numpods)
}

// func podFailure(t *testing.T, kcfg *kubeconfig.KubeConfig) {
// 	ncfg := `
// apiVersion: apate.opendc.org/v1
// kind: NodeConfiguration
// metadata:
// 	name: test-node1
// spec:
// 	replicas: 3
// 	resources:
// 		memory: 5G
// 		cpu: 1000
// 		storage: 5T
// 		ephemeral_storage: 120G
// 		max_pods: 150
// `

// 	pcfg := `
// apiVersion: apate.opendc.org/v1
// kind: PodConfiguration
// metadata:
// 	name: test-pod1
// spec:
// 	inline:
// 		pod_status: FAILED
// `

// 	err := kubectl.Create([]byte(ncfg), kcfg)
// 	assert.NoError(t, err)
// 	time.Sleep(time.Second * 60)

// 	cluster, err := kubernetes.ClusterFromKubeConfig(kcfg)
// 	assert.NoError(t, err)

// 	// Check if everything is ready
// 	ready, _ := getApateletWaitForCondition(t, cluster, func(apatelets []*corev1.Node) bool {
// 		assert.Equal(t, 1, len(apatelets))
// 		apatelet := apatelets[0]

// 		for _, c := range apatelet.Status.Conditions {
// 			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
// 				return true
// 			}
// 		}

// 		return false
// 	})

// 	assert.True(t, ready)

// 	// Deploy pods

// 	// run scenario

// 	// assert state
// }
