package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	apateRun "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	cp "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

// NODE DEPLOYMENT
func TestSimpleNodeDeploymentRoutine(t *testing.T) {
	testSimpleNodeDeployment(t, env.Routine)
}

func TestSimpleNodeDeploymentDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testSimpleNodeDeployment(t, env.Docker)
}

const SimpleNodeDeployment = `
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

// To run this, make sure ./config/kind.yml is put in the right directory (/tmp/apate/manager)
// or the env var CP_MANAGER_CONFIG_LOCATION point to it
func testSimpleNodeDeployment(t *testing.T, rt env.RunType) {
	setup(t, "TestSimpleNodeDeployment"+string(rt), rt)

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
	rc := SimpleNodeDeployment

	err := kubectl.Create([]byte(rc), kcfg)
	assert.NoError(t, err)
	log.Println("Waiting before querying k8s")
	time.Sleep(waitTimeout)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	log.Println("Getting number of nodes from k8s")
	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 3, nodes)
}

// NODE FAILURE
func TestNodeFailureDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testNodeFailure(t, env.Docker)
}

func TestNodeFailureRoutine(t *testing.T) {
	testNodeFailure(t, env.Routine)
}

func testNodeFailure(t *testing.T, rt env.RunType) {
	setup(t, "TestNodeFailure"+string(rt), rt)

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

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)

	assert.NoError(t, err)

	// Check if everything is ready
	ready, apatelets := getApateletWaitForCondition(t, cluster, num, createApateletConditionFunction(t, num, corev1.ConditionTrue))

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

	// Ideally we would wait 10s but as these tests are quite flaky we use 8 for safety
	time.Sleep(time.Second * 8)
	nodes, err := cluster.GetNodes()
	assert.NoError(t, err)

	assert.True(t, createApateletConditionFunction(t, num, corev1.ConditionTrue)(getApatelets(nodes)))

	stopped, _ := getApateletWaitForCondition(t, cluster, num, createApateletConditionFunction(t, num, corev1.ConditionUnknown))
	assert.True(t, stopped)
}

// SHUTDOWN APATELET
func TestShutdownApateletRoutine(t *testing.T) {
	testShutdownApatelet(t, env.Routine)
}

func TestShutdownApateletDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testShutdownApatelet(t, env.Docker)
}

func testShutdownApatelet(t *testing.T, rt env.RunType) {
	setup(t, "TestShutdownApatelet", rt)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.StartControlPlane(ctx, runner.New())

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

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

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	log.Println("Getting number of nodes from k8s")
	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 3, nodes)

	cancel()

	time.Sleep(time.Second * 30)

	teardown(t)
}

// Test nops the spawning of apatelets
const Test env.RunType = "TEST"

// SHUTDOWN APATELETS
func TestShutdownApateletApateletSide(t *testing.T) {
	setup(t, "TestShutdownApateletApateletSide", Test)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	registry := runner.New()

	var testRoutineRunner runner.ApateletRunner = &TestRoutineRunner{}
	registry.RegisterRunner(Test, &testRoutineRunner)

	go cp.StartControlPlane(ctx, registry)

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	rc := `
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: e2e-deployment
spec:
    replicas: 1
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

	environment, err := env.ApateletEnv()
	assert.NoError(t, err)
	environment.KubeConfigLocation = env.ControlPlaneEnv().KubeConfigLocation

	apctx, apcancel := context.WithCancel(context.Background())

	apateletEnv := environment

	// Apatelets should figure out their own ports when running in go routines
	apateletEnv.KubernetesPort = 0
	apateletEnv.MetricsPort = 0
	apateletEnv.ListenPort = 0

	readyCh := make(chan struct{}, 1)
	stop := make(chan os.Signal, 1)

	readyCh <- struct{}{}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Apatelet failed to start: %v\n", r)
				readyCh <- struct{}{} // Just continue to next one. Don't retry, as the resources may have been removed from the queue already
			}
		}()
		err1 := apateRun.StartApateletInternal(apctx, apateletEnv, readyCh, stop)
		if err1 != nil {
			log.Printf("Apatelet failed to start: %v\n", err1)
			readyCh <- struct{}{}
		}
	}()

	<-readyCh

	time.Sleep(time.Second * 30)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	log.Println("Getting number of nodes from k8s")
	nodes, err := cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 2, nodes)

	stop <- syscall.SIGTERM
	time.Sleep(time.Second * 30)

	log.Println("Getting number of nodes from k8s")
	nodes, err = cluster.GetNumberOfNodes()
	assert.NoError(t, err)
	assert.Equal(t, 1, nodes)

	apcancel()
	cancel()

	teardown(t)
}

// UP-DOWN SCALE
func TestUpDownScaleRoutine(t *testing.T) {
	testUpDownScale(t, env.Routine)
}

func TestUpDownScaleDocker(t *testing.T) {
	if detectCI() {
		t.Skip()
	}
	testUpDownScale(t, env.Docker)
}

func testUpDownScale(t *testing.T, rt env.RunType) {
	setup(t, "TestScale5000", rt)

	ctx, cancel := context.WithCancel(context.Background())

	// Start CP
	go cp.StartControlPlane(ctx, runner.New())

	// Wait
	waitForCP(t)

	kcfg := getKubeConfig(t)
	time.Sleep(time.Second * 5)

	rc1 := SimpleNodeDeployment

	rc2 := `
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: e2e-deployment
spec:
    replicas: 100
    resources:
        memory: 5G
        cpu: 1000
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
`

	err := kubectl.Apply([]byte(rc1), kcfg)
	assert.NoError(t, err)

	cmh := kubernetes.NewClusterManagerHandler()
	cluster, err := cmh.NewClusterFromKubeConfig(kcfg)
	assert.NoError(t, err)

	checkNodes(t, cluster, 3)

	err = kubectl.Apply([]byte(rc2), kcfg)
	assert.NoError(t, err)
	log.Println("Waiting before querying k8s")

	println("UPSCALING")

	checkNodes(t, cluster, 101)

	println("DOWNSCALING")

	err = kubectl.Apply([]byte(rc1), kcfg)
	assert.NoError(t, err)
	log.Println("Waiting before querying k8s")

	checkNodes(t, cluster, 3)

	println("DELETING")

	err = kubectl.Delete([]byte(rc1), kcfg)
	assert.NoError(t, err)
	log.Println("Waiting before querying k8s")

	checkNodes(t, cluster, 1)

	cancel()

	teardown(t)
}
