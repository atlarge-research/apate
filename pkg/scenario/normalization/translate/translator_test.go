package translate

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"

	"github.com/docker/go-units"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

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
	flags := EventFlags{
		events.NodeCreatePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeCreatePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeUpdatePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeUpdatePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeDeletePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeDeletePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeGetPodResponse:                 any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodResponsePercentage:       any.MarshalOrDie(100),
		events.NodeGetPodStatusResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodStatusResponsePercentage: any.MarshalOrDie(100),
		events.NodeGetPodsResponse:                any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodsResponsePercentage:      any.MarshalOrDie(100),
		events.NodePingResponse:                   any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodePingResponsePercentage:         any.MarshalOrDie(100),
	}
	assert.EqualValues(t, flags, newTask.NodeEventFlags)
}

func TestNetworkLatency(t *testing.T) {
	newTask := getApateletTask(t, `
network_latency:
 latency_msec: 100
`)
	flags := EventFlags{
		events.NodeAddedLatencyEnabled: any.MarshalOrDie(true),
		events.NodeAddedLatencyMsec:    any.MarshalOrDie(int32(100)),
	}
	assert.EqualValues(t, flags, newTask.NodeEventFlags)
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
	assert.EqualValues(t, EventFlags{
		events.NodeCreatePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeCreatePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeUpdatePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeUpdatePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeDeletePodResponse:              any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeDeletePodResponsePercentage:    any.MarshalOrDie(100),
		events.NodeGetPodResponse:                 any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodResponsePercentage:       any.MarshalOrDie(100),
		events.NodeGetPodStatusResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodStatusResponsePercentage: any.MarshalOrDie(100),
		events.NodeGetPodsResponse:                any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodsResponsePercentage:      any.MarshalOrDie(100),
		events.NodePingResponse:                   any.MarshalOrDie(scenario.Response_NORMAL),
		events.NodePingResponsePercentage:         any.MarshalOrDie(0),
	}, newTask.NodeEventFlags)
}

func TestNoTimeoutNoHeartbeat(t *testing.T) {
	newTask := getApateletTask(t, `
no_timeout_no_heartbeat: {}
`)
	assert.EqualValues(t, EventFlags{
		events.NodePingResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodePingResponsePercentage: any.MarshalOrDie(100),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStateCreatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: CREATE_POD
 response: ERROR
 percentage: 42
`)
	assert.EqualValues(t, EventFlags{
		events.NodeCreatePodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.NodeCreatePodResponsePercentage: any.MarshalOrDie(42),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: UPDATE_POD
 response: TIMEOUT
 percentage: 15
`)
	assert.EqualValues(t, EventFlags{
		events.NodeUpdatePodResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeUpdatePodResponsePercentage: any.MarshalOrDie(15),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: DELETE_POD
 response: ERROR
 percentage: 100
`)
	assert.EqualValues(t, EventFlags{
		events.NodeDeletePodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.NodeDeletePodResponsePercentage: any.MarshalOrDie(100),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: GET_POD
 response: ERROR
 percentage: 14
`)
	assert.EqualValues(t, EventFlags{
		events.NodeGetPodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.NodeGetPodResponsePercentage: any.MarshalOrDie(14),
	}, newTask.NodeEventFlags)
}
func TestNodeResponseStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: GET_POD_STATUS
 response: TIMEOUT
 percentage: 42
`)
	assert.EqualValues(t, EventFlags{
		events.NodeGetPodStatusResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodStatusResponsePercentage: any.MarshalOrDie(42),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStateGetPods(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: GET_PODS
 response: TIMEOUT
 percentage: 65
`)
	assert.EqualValues(t, EventFlags{
		events.NodeGetPodsResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeGetPodsResponsePercentage: any.MarshalOrDie(65),
	}, newTask.NodeEventFlags)
}

func TestNodeResponseStatePing(t *testing.T) {
	newTask := getApateletTask(t, `
node_response_state:
 type: PING
 response: ERROR
 percentage: 50
`)
	assert.EqualValues(t, EventFlags{
		events.NodePingResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.NodePingResponsePercentage: any.MarshalOrDie(50),
	}, newTask.NodeEventFlags)
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
	assert.EqualValues(t, EventFlags{
		events.NodeEnableResourceAlteration: any.MarshalOrDie(true),
		events.NodeCPUUsage:                 any.MarshalOrDie(42),
		events.NodeMemoryUsage:              any.MarshalOrDie(21 * units.GiB),
		events.NodeStorageUsage:             any.MarshalOrDie(84 * units.MiB),
		events.NodeEphemeralStorageUsage:    any.MarshalOrDie(105 * units.KiB),
	}, newTask.NodeEventFlags)
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
	assert.EqualValues(t, EventFlags{
		events.PodCreatePodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.PodCreatePodResponsePercentage: any.MarshalOrDie(42),
	}, newTask.PodConfigs[0].EventFlags)
}

func TestPodResponseStateUpdatePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
 type: UPDATE_POD
 response: TIMEOUT
 percentage: 15
`)
	assert.EqualValues(t, EventFlags{
		events.PodUpdatePodResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.PodUpdatePodResponsePercentage: any.MarshalOrDie(15),
	}, newTask.PodConfigs[0].EventFlags)
}

func TestPodResponseStateDeletePod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
 type: DELETE_POD
 response: ERROR
 percentage: 100
`)
	assert.EqualValues(t, EventFlags{
		events.PodDeletePodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.PodDeletePodResponsePercentage: any.MarshalOrDie(100),
	}, newTask.PodConfigs[0].EventFlags)
}

func TestPodResponseStateGetPod(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
 type: GET_POD
 response: ERROR
 percentage: 14
`)
	assert.EqualValues(t, EventFlags{
		events.PodGetPodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.PodGetPodResponsePercentage: any.MarshalOrDie(14),
	}, newTask.PodConfigs[0].EventFlags)
}

func TestPodResponseStateGetPodStatus(t *testing.T) {
	newTask := getApateletTask(t, `
pod_response_state:
 type: GET_POD_STATUS
 response: TIMEOUT
 percentage: 42
`)
	assert.EqualValues(t, EventFlags{
		events.PodGetPodStatusResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.PodGetPodStatusResponsePercentage: any.MarshalOrDie(42),
	}, newTask.NodeEventFlags)
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
	assert.EqualValues(t, EventFlags{
		events.PodUpdatePodStatus:           any.MarshalOrDie(scenario.PodStatus_POD_FAILED),
		events.PodUpdatePodStatusPercentage: any.MarshalOrDie(15),
	}, newTask.PodConfigs[0].EventFlags)
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

func translateYaml(t *testing.T, data []byte) *controlplane.Task {
	json, err := yaml.YAMLToJSON(data)
	assert.NoError(t, err)

	var task controlplane.Task
	err = protojson.Unmarshal(json, &task)
	assert.NoError(t, err)

	return &task
}
