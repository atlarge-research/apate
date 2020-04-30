package events

import (
	"bytes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/golang/protobuf/jsonpb"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ghodss/yaml"
)

// Node events
func TestNodeFailure(t *testing.T) {
	newTask := getApateletTask(t, `
node_failure: {}
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
			LifecycleState: &apatelet.LifecycleState{
				CreatePodAction:              scenario.LifecycleAction_TIMEOUT,
				CreatePodActionPercentage:    100,
				UpdatePodAction:              scenario.LifecycleAction_TIMEOUT,
				UpdatePodActionPercentage:    100,
				DeletePodAction:              scenario.LifecycleAction_TIMEOUT,
				DeletePodActionPercentage:    100,
				GetPodAction:                 scenario.LifecycleAction_TIMEOUT,
				GetPodActionPercentage:       100,
				GetPodStatusAction:           scenario.LifecycleAction_TIMEOUT,
				GetPodStatusActionPercentage: 100,
			},
			GetPodsAction:           scenario.LifecycleAction_TIMEOUT,
			GetPodsActionPercentage: 100,
			PingAction:              scenario.LifecycleAction_TIMEOUT,
			PingActionPercentage:    100,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNetworkLatency(t *testing.T) {
	newTask := getApateletTask(t, `
network_latency:
  latency_msec: 100
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{}},
		ResourceState:      &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{
			AddedLatencyEnabled: true,
			AddedLatencyMsec:    100,
		},
	}), newTask.Event)
}

func TestNegativeNetworkLatency(t *testing.T) {
	getApateletErroredTask(t, `
network_latency:
  latency_msec: -100
`, "latency should be at least 0")
}

func TestTimeoutKeepHeartbeat(t *testing.T) {
	newTask := getApateletTask(t, `
timeout_keep_heartbeat: {}
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
			LifecycleState: &apatelet.LifecycleState{
				CreatePodAction:              scenario.LifecycleAction_TIMEOUT,
				CreatePodActionPercentage:    100,
				UpdatePodAction:              scenario.LifecycleAction_TIMEOUT,
				UpdatePodActionPercentage:    100,
				DeletePodAction:              scenario.LifecycleAction_TIMEOUT,
				DeletePodActionPercentage:    100,
				GetPodAction:                 scenario.LifecycleAction_TIMEOUT,
				GetPodActionPercentage:       100,
				GetPodStatusAction:           scenario.LifecycleAction_TIMEOUT,
				GetPodStatusActionPercentage: 100,
			},
			GetPodsAction:           scenario.LifecycleAction_TIMEOUT,
			GetPodsActionPercentage: 100,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNoTimeoutNoHeartbeat(t *testing.T) {
	newTask := getApateletTask(t, `
no_timeout_no_heartbeat: {}
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
			LifecycleState:       &apatelet.LifecycleState{},
			PingAction:           scenario.LifecycleAction_TIMEOUT,
			PingActionPercentage: 100,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateCreatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: CREATE_POD
  action: ERROR
  percentage: 42
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{
			CreatePodAction:           scenario.LifecycleAction_ERROR,
			CreatePodActionPercentage: 42,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: UPDATE_POD
  action: TIMEOUT
  percentage: 15
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{
			UpdatePodAction:           scenario.LifecycleAction_TIMEOUT,
			UpdatePodActionPercentage: 15,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: DELETE_POD
  action: ERROR
  percentage: 100
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{
			DeletePodAction:           scenario.LifecycleAction_ERROR,
			DeletePodActionPercentage: 100,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: GET_POD
  action: ERROR
  percentage: 14
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{
			GetPodAction:           scenario.LifecycleAction_ERROR,
			GetPodActionPercentage: 14,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: GET_POD_STATUS
  action: TIMEOUT
  percentage: 42
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{
			GetPodStatusAction:           scenario.LifecycleAction_TIMEOUT,
			GetPodStatusActionPercentage: 42,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateGetPods(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: GET_PODS
  action: TIMEOUT
  percentage: 65
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
			LifecycleState:          &apatelet.LifecycleState{},
			GetPodsAction:           scenario.LifecycleAction_TIMEOUT,
			GetPodsActionPercentage: 65,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStatePing(t *testing.T) {
	newTask := getApateletTask(t, `
node_lifecycle_state:
  type: PING
  action: ERROR
  percentage: 50
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
			LifecycleState:       &apatelet.LifecycleState{},
			PingAction:           scenario.LifecycleAction_ERROR,
			PingActionPercentage: 50,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeLifecycleStateLessThan0(t *testing.T) {
	getApateletErroredTask(t, `
node_lifecycle_state:
  type: PING
  action: ERROR
  percentage: -50
`, "percentage should be between 0 and 100")
}

func TestNodeLifecycleStateMoreThan100(t *testing.T) {
	getApateletErroredTask(t, `
node_lifecycle_state:
  type: PING
  action: ERROR
  percentage: 420
`, "percentage should be between 0 and 100")
}

func TestResourcePressure(t *testing.T) {
	newTask := getApateletTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21
  storage_usage: 84
  ephemeral_storage_usage: 105
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{LifecycleState: &apatelet.LifecycleState{}},
		ResourceState: &apatelet.NodeState_ResourceState{
			EnableResourceAlteration: true,
			CpuUsage:                 42,
			MemoryUsage:              21,
			StorageUsage:             84,
			EphemeralStorageUsage:    105,
		},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestResourcePressureCpuBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: -42
  memory_usage: 21
  storage_usage: 84
  ephemeral_storage_usage: 105
`, "CPU usage should be at least 0")
}

func TestResourcePressureMemoryBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: -21
  storage_usage: 84
  ephemeral_storage_usage: 105
`, "memoy usage should be at least 0")
}

func TestResourcePressureStorageBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21
  storage_usage: -84
  ephemeral_storage_usage: 105
`, "storage usage should be at least 0")
}

func TestResourcePressureEphemeralStorageBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21
  storage_usage: 84
  ephemeral_storage_usage: -105
`, "ephemeral storage usage should be at least 0")
}

// Pod events
func TestPodLifecycleStateCreatePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_lifecycle_state:
  type: CREATE_POD
  action: ERROR
  percentage: 42
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{
			CreatePodAction:           scenario.LifecycleAction_ERROR,
			CreatePodActionPercentage: 42,
		}},
	}), newTask.Event)
}

func TestPodLifecycleStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_lifecycle_state:
  type: UPDATE_POD
  action: TIMEOUT
  percentage: 15
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{
			UpdatePodAction:           scenario.LifecycleAction_TIMEOUT,
			UpdatePodActionPercentage: 15,
		}},
	}), newTask.Event)
}

func TestPodLifecycleStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_lifecycle_state:
  type: DELETE_POD
  action: ERROR
  percentage: 100
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{
			DeletePodAction:           scenario.LifecycleAction_ERROR,
			DeletePodActionPercentage: 100,
		}},
	}), newTask.Event)
}

func TestPodLifecycleStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_lifecycle_state:
  type: GET_POD
  action: ERROR
  percentage: 14
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{
			GetPodAction:           scenario.LifecycleAction_ERROR,
			GetPodActionPercentage: 14,
		}},
	}), newTask.Event)
}

func TestPodLifecycleStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
pod_lifecycle_state:
  type: GET_POD_STATUS
  action: TIMEOUT
  percentage: 42
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{
			GetPodStatusAction:           scenario.LifecycleAction_TIMEOUT,
			GetPodStatusActionPercentage: 42,
		}},
	}), newTask.Event)
}

func TestPodLifecycleStateGetPods(t *testing.T) {
	getApateletErroredTask(t, `
pod_lifecycle_state:
  type: GET_PODS
  action: TIMEOUT
  percentage: 65
`, "can't alter the GetPods / Ping response on pod level")
}

func TestPodLifecycleStatePing(t *testing.T) {
	getApateletErroredTask(t, `
pod_lifecycle_state:
  type: PING
  action: ERROR
  percentage: 50
`, "can't alter the GetPods / Ping response on pod level")
}

func TestPodLifecycleStateLessThan0(t *testing.T) {
	getApateletErroredTask(t, `
pod_lifecycle_state:
  type: PING
  action: ERROR
  percentage: -50
`, "percentage should be between 0 and 100")
}

func TestPodLifecycleStateMoreThan100(t *testing.T) {
	getApateletErroredTask(t, `
pod_lifecycle_state:
  type: PING
  action: ERROR
  percentage: 420
`, "percentage should be between 0 and 100")
}

func TestPodStatusUpdate(t *testing.T) {
	newTask := getApateletTask(t, `
pod_status_update:
  new_status: POD_FAILED
  percentage: 15
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState:   &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{}},
		PodStatus:           scenario.PodStatus_POD_FAILED,
		PodStatusPercentage: 15,
	}), newTask.Event)
}

func TestPodStatusUpdateLessThan0(t *testing.T) {
	getApateletErroredTask(t, `
pod_status_update:
  new_status: POD_FAILED
  percentage: -15
`, "percentage should be between 0 and 100")
}

func TestPodStatusUpdateMoreThan100(t *testing.T) {
	getApateletErroredTask(t, `
pod_status_update:
  new_status: POD_FAILED
  percentage: 150
`, "percentage should be between 0 and 100")
}

func TestPodStartTimeUpdate(t *testing.T) {
	newTask := getApateletTask(t, `
pod_start_time_update:
  new_start_time: "2020-04-30T11:32:05+0000"
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodLifecycleState: &apatelet.PodState_PodLifecycleState{LifecycleState: &apatelet.LifecycleState{}},
		PodStartTime:      "2020-04-30T11:32:05+0000",
	}), newTask.Event)
}

// Utils
func getApateletTask(t *testing.T, task string) *apatelet.Task {
	origTask := translateYaml(t, []byte(task))
	newTask := &apatelet.Task{}
	err := NewEventTranslator(origTask, newTask).TranslateEvent()
	assert.NoError(t, err)
	return newTask
}

func getApateletErroredTask(t *testing.T, task string, msg string) {
	origTask := translateYaml(t, []byte(task))
	newTask := &apatelet.Task{}
	err := NewEventTranslator(origTask, newTask).TranslateEvent()
	assert.Error(t, err, msg)
}

func createNodeEvent(nodeState *apatelet.NodeState) *apatelet.Task_NodeEvent {
	return &apatelet.Task_NodeEvent{NodeEvent: &apatelet.NodeEvent{NodeState: nodeState}}
}

func createPodEvent(podState *apatelet.PodState) *apatelet.Task_PodEvent {
	return &apatelet.Task_PodEvent{PodEvent: &apatelet.PodEvent{PodState: podState}}
}

func translateYaml(t *testing.T, data []byte) *controlplane.Task {
	json, err := yaml.YAMLToJSON(data)
	assert.NoError(t, err)

	var task controlplane.Task
	err = jsonpb.Unmarshal(bytes.NewReader(json), &task)
	assert.NoError(t, err)

	return &task
}
