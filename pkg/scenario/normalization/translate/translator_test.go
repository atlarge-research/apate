package translate

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"

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
		events.NodeCreatePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeUpdatePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeDeletePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodResponse:       any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodStatusResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodsResponse:      any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodePingResponse:         any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
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
		events.NodeCreatePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeUpdatePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeDeletePodResponse:    any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodResponse:       any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodStatusResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodeGetPodsResponse:      any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
		events.NodePingResponse:         any.MarshalOrDie(scenario.Response_RESPONSE_NORMAL),
	}, newTask.NodeEventFlags)
}

func TestNoTimeoutNoHeartbeat(t *testing.T) {
	newTask := getApateletTask(t, `
no_timeout_no_heartbeat: {}
`)
	assert.EqualValues(t, EventFlags{
		events.NodePingResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
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
		events.NodeCreatePodResponse: any.MarshalOrDie(scenario.Response_RESPONSE_ERROR),
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
		events.NodeUpdatePodResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
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
		events.NodeDeletePodResponse: any.MarshalOrDie(scenario.Response_RESPONSE_ERROR),
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
		events.NodeGetPodResponse: any.MarshalOrDie(scenario.Response_RESPONSE_ERROR),
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
		events.NodeGetPodStatusResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
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
		events.NodeGetPodsResponse: any.MarshalOrDie(scenario.Response_RESPONSE_TIMEOUT),
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
		events.NodePingResponse: any.MarshalOrDie(scenario.Response_RESPONSE_ERROR),
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
