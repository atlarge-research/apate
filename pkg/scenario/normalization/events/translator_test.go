package events

import (
	"testing"

	"github.com/docker/go-units"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/stretchr/testify/assert"

	"github.com/ghodss/yaml"
)

// Node events
func TestNodeFailure(t *testing.T) {
	newTask := getApateletTask(t, `
node_failure: {}
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{
			ResponseState: &apatelet.ResponseState{
				CreatePodResponse:              scenario.Response_TIMEOUT,
				CreatePodResponsePercentage:    100,
				UpdatePodResponse:              scenario.Response_TIMEOUT,
				UpdatePodResponsePercentage:    100,
				DeletePodResponse:              scenario.Response_TIMEOUT,
				DeletePodResponsePercentage:    100,
				GetPodResponse:                 scenario.Response_TIMEOUT,
				GetPodResponsePercentage:       100,
				GetPodStatusResponse:           scenario.Response_TIMEOUT,
				GetPodStatusResponsePercentage: 100,
			},
			GetPodsResponse:           scenario.Response_TIMEOUT,
			GetPodsResponsePercentage: 100,
			PingResponse:              scenario.Response_TIMEOUT,
			PingResponsePercentage:    100,
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
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
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
		NodeResponseState: &apatelet.NodeState_NodeResponseState{
			ResponseState: &apatelet.ResponseState{
				CreatePodResponse:              scenario.Response_TIMEOUT,
				CreatePodResponsePercentage:    100,
				UpdatePodResponse:              scenario.Response_TIMEOUT,
				UpdatePodResponsePercentage:    100,
				DeletePodResponse:              scenario.Response_TIMEOUT,
				DeletePodResponsePercentage:    100,
				GetPodResponse:                 scenario.Response_TIMEOUT,
				GetPodResponsePercentage:       100,
				GetPodStatusResponse:           scenario.Response_TIMEOUT,
				GetPodStatusResponsePercentage: 100,
			},
			GetPodsResponse:           scenario.Response_TIMEOUT,
			GetPodsResponsePercentage: 100,
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
		NodeResponseState: &apatelet.NodeState_NodeResponseState{
			ResponseState:          &apatelet.ResponseState{},
			PingResponse:           scenario.Response_TIMEOUT,
			PingResponsePercentage: 100,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateCreatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: CREATE_POD
  response: ERROR
  percentage: 42
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{
			CreatePodResponse:           scenario.Response_ERROR,
			CreatePodResponsePercentage: 42,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: UPDATE_POD
  response: TIMEOUT
  percentage: 15
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{
			UpdatePodResponse:           scenario.Response_TIMEOUT,
			UpdatePodResponsePercentage: 15,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: DELETE_POD
  response: ERROR
  percentage: 100
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{
			DeletePodResponse:           scenario.Response_ERROR,
			DeletePodResponsePercentage: 100,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: GET_POD
  response: ERROR
  percentage: 14
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{
			GetPodResponse:           scenario.Response_ERROR,
			GetPodResponsePercentage: 14,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: GET_POD_STATUS
  response: TIMEOUT
  percentage: 42
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{
			GetPodStatusResponse:           scenario.Response_TIMEOUT,
			GetPodStatusResponsePercentage: 42,
		}},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateGetPods(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: GET_PODS
  response: TIMEOUT
  percentage: 65
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{
			ResponseState:             &apatelet.ResponseState{},
			GetPodsResponse:           scenario.Response_TIMEOUT,
			GetPodsResponsePercentage: 65,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStatePing(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
  type: PING
  response: ERROR
  percentage: 50
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{
			ResponseState:          &apatelet.ResponseState{},
			PingResponse:           scenario.Response_ERROR,
			PingResponsePercentage: 50,
		},
		ResourceState:     &apatelet.NodeState_ResourceState{},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestNodeResponseStateLessThan0(t *testing.T) {
	getApateletErroredTask(t, `
node_response_state:
  type: PING
  response: ERROR
  percentage: -50
`, "percentage should be between 0 and 100")
}

func TestNodeResponseStateMoreThan100(t *testing.T) {
	getApateletErroredTask(t, `
node_response_state:
  type: PING
  response: ERROR
  percentage: 420
`, "percentage should be between 0 and 100")
}

func TestResourcePressure(t *testing.T) {
	newTask := getApateletTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21GB
  storage_usage: 84MB
  ephemeral_storage_usage: 105KB
`)
	assert.Equal(t, createNodeEvent(&apatelet.NodeState{
		NodeResponseState: &apatelet.NodeState_NodeResponseState{ResponseState: &apatelet.ResponseState{}},
		ResourceState: &apatelet.NodeState_ResourceState{
			EnableResourceAlteration: true,
			CpuUsage:                 42,
			MemoryUsage:              21 * units.GiB,
			StorageUsage:             84 * units.MiB,
			EphemeralStorageUsage:    105 * units.KiB,
		},
		AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
	}), newTask.Event)
}

func TestResourcePressureCpuBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: -42
  memory_usage: 21G
  storage_usage: 84M
  ephemeral_storage_usage: 105K
`, "CPU usage should be at least 0")
}

func TestResourcePressureMemoryBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: -21G
  storage_usage: 84M
  ephemeral_storage_usage: 105K
`, "memoy usage should be at least 0")
}

func TestResourcePressureStorageBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21G
  storage_usage: -84M
  ephemeral_storage_usage: 105K
`, "storage usage should be at least 0")
}

func TestResourcePressureEphemeralStorageBelow0(t *testing.T) {
	getApateletErroredTask(t, `
resource_pressure:
  cpu_usage: 42
  memory_usage: 21G
  storage_usage: 84M
  ephemeral_storage_usage: -105K
`, "ephemeral storage usage should be at least 0")
}

// Pod events
func TestPodResponseStateCreatePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
  type: CREATE_POD
  response: ERROR
  percentage: 42
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{
			CreatePodResponse:           scenario.Response_ERROR,
			CreatePodResponsePercentage: 42,
		}},
	}), newTask.Event)
}

func TestPodResponseStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
  type: UPDATE_POD
  response: TIMEOUT
  percentage: 15
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{
			UpdatePodResponse:           scenario.Response_TIMEOUT,
			UpdatePodResponsePercentage: 15,
		}},
	}), newTask.Event)
}

func TestPodResponseStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
  type: DELETE_POD
  response: ERROR
  percentage: 100
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{
			DeletePodResponse:           scenario.Response_ERROR,
			DeletePodResponsePercentage: 100,
		}},
	}), newTask.Event)
}

func TestPodResponseStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
  type: GET_POD
  response: ERROR
  percentage: 14
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{
			GetPodResponse:           scenario.Response_ERROR,
			GetPodResponsePercentage: 14,
		}},
	}), newTask.Event)
}

func TestPodResponseStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
  type: GET_POD_STATUS
  response: TIMEOUT
  percentage: 42
`)
	assert.Equal(t, createPodEvent(&apatelet.PodState{
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{
			GetPodStatusResponse:           scenario.Response_TIMEOUT,
			GetPodStatusResponsePercentage: 42,
		}},
	}), newTask.Event)
}

func TestPodResponseStateGetPods(t *testing.T) {
	getApateletErroredTask(t, `
pod_response_state:
  type: GET_PODS
  response: TIMEOUT
  percentage: 65
`, "can't alter the GetPods / Ping response on pod level")
}

func TestPodResponseStatePing(t *testing.T) {
	getApateletErroredTask(t, `
pod_response_state:
  type: PING
  response: ERROR
  percentage: 50
`, "can't alter the GetPods / Ping response on pod level")
}

func TestPodResponseStateLessThan0(t *testing.T) {
	getApateletErroredTask(t, `
pod_response_state:
  type: PING
  response: ERROR
  percentage: -50
`, "percentage should be between 0 and 100")
}

func TestPodResponseStateMoreThan100(t *testing.T) {
	getApateletErroredTask(t, `
pod_response_state:
  type: PING
  response: ERROR
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
		PodResponseState:    &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{}},
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
		PodResponseState: &apatelet.PodState_PodResponseState{ResponseState: &apatelet.ResponseState{}},
		PodStartTime:     "2020-04-30T11:32:05+0000",
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
	err = protojson.Unmarshal(json, &task)
	assert.NoError(t, err)

	return &task
}
