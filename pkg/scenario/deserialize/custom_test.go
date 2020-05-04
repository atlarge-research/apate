package deserialize

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"

	"github.com/docker/go-units"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

func TestCustomEventsNormal(t *testing.T) {
	ps := getPublicScenario(t, `
custom_flags:
    node_create_pod_response: ERROR
    node_create_pod_response_percentage: 5

    node_update_pod_response: TIMEOUT
    node_update_pod_response_percentage: 10

    node_delete_pod_response: ERROR
    node_delete_pod_response_percentage: 15

    node_get_pod_response: NORMAL
    node_get_pod_response_percentage: 20

    node_get_pod_status_response: ERROR
    node_get_pod_status_response_percentage: 25

    node_get_pods_response: TIMEOUT
    node_get_pods_response_percentage: 50

    node_ping_response: TIMEOUT
    node_ping_response_percentage: 75

    node_enable_resource_alteration: true
    node_memory_usage: 10M
    node_cpu_usage: 200
    node_storage_usage: 500G
    node_ephemeral_storage_usage: 10T

    node_added_latency_enabled: true
    node_added_latency_msec: 500

pod_configs:
    -
        metadata_name: a
        custom_flags:
            pod_create_pod_response: ERROR
            pod_create_pod_response_percentage: 10

            pod_update_pod_response: NORMAL
            pod_update_pod_response_percentage: 20

            pod_delete_pod_response: NORMAL
            pod_delete_pod_response_percentage: 30

            pod_get_pod_response: TIMEOUT
            pod_get_pod_response_percentage: 40
    -
        metadata_name: b
        custom_flags:
            pod_get_pod_status_response: NORMAL
            pod_get_pod_status_response_percentage: 44

            pod_update_pod_status: POD_PENDING
            pod_update_pod_status_percentage: 55
`)
	// Node custom flags
	ncf := ps.Tasks[0].GetCustomFlags().CustomFlags

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_ERROR), ncf[events.NodeCreatePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(5), ncf[events.NodeCreatePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_TIMEOUT), ncf[events.NodeUpdatePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(10), ncf[events.NodeUpdatePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_ERROR), ncf[events.NodeDeletePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(15), ncf[events.NodeDeletePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_NORMAL), ncf[events.NodeGetPodResponse])
	assert.EqualValues(t, any.MarshalOrDie(20), ncf[events.NodeGetPodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_ERROR), ncf[events.NodeGetPodStatusResponse])
	assert.EqualValues(t, any.MarshalOrDie(25), ncf[events.NodeGetPodStatusResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_TIMEOUT), ncf[events.NodeGetPodsResponse])
	assert.EqualValues(t, any.MarshalOrDie(50), ncf[events.NodeGetPodsResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_TIMEOUT), ncf[events.NodePingResponse])
	assert.EqualValues(t, any.MarshalOrDie(75), ncf[events.NodePingResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(true), ncf[events.NodeEnableResourceAlteration])
	assert.EqualValues(t, any.MarshalOrDie(10*units.MiB), ncf[events.NodeMemoryUsage])
	assert.EqualValues(t, any.MarshalOrDie(200), ncf[events.NodeCPUUsage])
	assert.EqualValues(t, any.MarshalOrDie(500*units.GiB), ncf[events.NodeStorageUsage])
	assert.EqualValues(t, any.MarshalOrDie(10*units.TiB), ncf[events.NodeEphemeralStorageUsage])

	assert.EqualValues(t, any.MarshalOrDie(true), ncf[events.NodeAddedLatencyEnabled])
	assert.EqualValues(t, any.MarshalOrDie(500), ncf[events.NodeAddedLatencyMsec])

	// Pod config 1
	pc1 := ps.Tasks[0].GetPodConfigs()[0]
	assert.EqualValues(t, "a", pc1.GetMetadataName())

	pcf1 := pc1.GetCustomFlags().CustomFlags

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_ERROR), pcf1[events.PodCreatePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(10), pcf1[events.PodCreatePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_NORMAL), pcf1[events.PodUpdatePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(20), pcf1[events.PodUpdatePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_NORMAL), pcf1[events.PodDeletePodResponse])
	assert.EqualValues(t, any.MarshalOrDie(30), pcf1[events.PodDeletePodResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_TIMEOUT), pcf1[events.PodGetPodResponse])
	assert.EqualValues(t, any.MarshalOrDie(40), pcf1[events.PodGetPodResponsePercentage])

	// Pod config 2
	pc2 := ps.Tasks[0].GetPodConfigs()[1]
	assert.EqualValues(t, "b", pc2.GetMetadataName())

	pcf2 := pc2.GetCustomFlags().CustomFlags

	assert.EqualValues(t, any.MarshalOrDie(scenario.Response_NORMAL), pcf2[events.PodGetPodStatusResponse])
	assert.EqualValues(t, any.MarshalOrDie(44), pcf2[events.PodGetPodStatusResponsePercentage])

	assert.EqualValues(t, any.MarshalOrDie(scenario.PodStatus_POD_PENDING), pcf2[events.PodUpdatePodStatus])
	assert.EqualValues(t, any.MarshalOrDie(55), pcf2[events.PodUpdatePodStatusPercentage])
}

func TestCustomInvalidKey(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    a: ERROR
`, "invalid custom flag key 'a'")
}

func TestCustomInvalidResponse(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    pod_get_pod_status_response: ERROR2
`, "invalid response 'ERROR2'")
}

func TestCustomInvalidPodStatus(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    pod_update_pod_status: F
`, "invalid pod status 'F'")
}

func TestCustomInvalidPercentLow(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    pod_update_pod_status_percentage: -10
`, "percentage should be between 0 and 100")
}

func TestCustomInvalidPercentHigh(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    pod_update_pod_status_percentage: 110
`, "percentage should be between 0 and 100")
}

func TestCustomInvalidLowInt(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    node_added_latency_msec: -110
`, "value should be at least 0")
}

func TestCustomInvalidBool(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    node_added_latency_enabled: fffffff
`)
}

func TestCustomInvalidInt(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    node_added_latency_msec: fffffff
`)
}

func TestCustomInvalidSize(t *testing.T) {
	getErrorPublicScenario(t, `
custom_flags:
    node_storage_usage: 100FFF
`)
}

// Utils
func getPublicScenario(t *testing.T, events string) *controlplane.PublicScenario {
	jsonBytes, ps := getJSONBytes(t, events)
	cfp := customFlagParser{scenario: ps}
	err := cfp.parse(jsonBytes)
	assert.NoError(t, err)
	return ps
}

func getErrorPublicScenario(t *testing.T, events string, msg ...string) {
	jsonBytes, ps := getJSONBytes(t, events)
	cfp := customFlagParser{scenario: ps}
	err := cfp.parse(jsonBytes)
	assert.Error(t, err, msg)
}

func getJSONBytes(t *testing.T, events string) ([]byte, *controlplane.PublicScenario) {
	scanner := bufio.NewScanner(strings.NewReader(events))
	indentedText := ""
	for scanner.Scan() {
		indentedLine := "        " + scanner.Text()
		indentedText += indentedLine + "\n"
	}

	jsonBytes, err := yaml.YAMLToJSON([]byte(`
nodes:
    -
        node_type: testnode
        memory: 2G
        cpu: 42
        storage: 2G
        ephemeral_storage: 2M
        max_pods: 42
node_groups:
    -
        group_name: testgroup1
        node_type: testnode
        amount: 42
tasks:
    -
        name: testtask2
        time: 10s
        node_groups:
            - all
` + indentedText))

	assert.NoError(t, err)

	var ps controlplane.PublicScenario
	if err := json.Unmarshal(jsonBytes, &ps); err != nil {
		return nil, nil
	}

	return jsonBytes, &ps
}
