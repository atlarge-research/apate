package e2e

import (
	"context"
	"fmt"
	"log"
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

func TestSimpleNodeDeploymentDocker(t *testing.T) {
	if !enableDockerApatelets {
		t.Skip()
	}
	testSimpleNodeDeployment(t, env.Docker)
}

func TestSimpleNodeDeploymentRoutine(t *testing.T) {
	testSimpleNodeDeployment(t, env.Routine)
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
	log.Println("Waiting before querying k8s")
	time.Sleep(time.Second * 60)

	// TODO: Is this correct?
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	log.Println("Getting number of nodes from k8s")
	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 3, nodes)
}

func TestNodeFailureDocker(t *testing.T) {
	if !enableDockerApatelets {
		t.Skip()
	}
	testNodeFailure(t, env.Docker)
}

func TestNodeFailureRoutine(t *testing.T) {
	testNodeFailure(t, env.Routine)
}

func testNodeFailure(t *testing.T, rt env.RunType) {
	setup(t, "testNodeFailure_"+string(rt), rt)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.StartControlPlane(ctx, runner.New())

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	// Test node failure
	nodeFailure(t, kcfg)

	cancel()

	teardown(t)
}

func nodeFailure(t *testing.T, kcfg *kubeconfig.KubeConfig) {
	num := 3

	ncfg := fmt.Sprintf(`
    apiVersion: apate.opendc.org/v1
    kind: NodeConfiguration
    metadata:
        name: test-deployment1
    spec:
        replicas: %d
        resources:
            memory: 5G
            cpu: 1000
            storage: 5T
            ephemeral_storage: 120G
            max_pods: 150
        tasks:
            - timestamp: 10s
              state:
                  node_failed: true
`, num)

	err := kubectl.Create([]byte(ncfg), kcfg)
	assert.NoError(t, err)
	time.Sleep(time.Second * 60)

	// TODO: Is this correct?
	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	// Check if everything is ready
	ready, apatelets := getApateletWaitForCondition(t, cluster, num, createConditionFunction(t, num, corev1.ConditionTrue))

	assert.True(t, ready)

	node := apatelets[0]

	cpu := node.Status.Capacity[corev1.ResourceCPU]
	v, _ := cpu.AsInt64()
	assert.Equal(t, int64(1000), v)

	eph := node.Status.Capacity[corev1.ResourceEphemeralStorage]
	v, _ = eph.AsInt64()
	assert.Equal(t, int64(120*1024*1024*1024), v)

	memory := node.Status.Capacity[corev1.ResourceMemory]
	v, _ = memory.AsInt64()
	assert.Equal(t, int64(5*1024*1024*1024), v)

	storage := node.Status.Capacity[corev1.ResourceStorage]
	v, _ = storage.AsInt64()
	assert.Equal(t, int64(5*1024*1024*1024*1024), v)

	// Check if they stopped
	runScenario(t)

	stopped, _ := getApateletWaitForCondition(t, cluster, num, createConditionFunction(t, num, corev1.ConditionUnknown))
	assert.True(t, stopped)
}

func getApateletWaitForCondition(t *testing.T, cluster *kubernetes.Cluster, numApatelets int, check func([]*corev1.Node) bool) (bool, []*corev1.Node) {
	for i := 0; i <= 10; i++ {
		// get nodes and check that there are 2
		nodes, err := cluster.GetNodes()
		assert.NoError(t, err)
		assert.Equal(t, numApatelets+1, len(nodes.Items))

		apatelets := getApatelets(nodes)
		assert.Equal(t, numApatelets, len(apatelets))

		if check(apatelets) {
			return true, apatelets
		}

		time.Sleep(time.Second * 10)
	}

	return false, nil
}

func getApatelets(nodes *corev1.NodeList) (node []*corev1.Node) {
	for _, v := range nodes.Items {
		v := v
		if strings.HasPrefix(v.Name, "apatelet-") {
			node = append(node, &v)
		}
	}
	return
}
