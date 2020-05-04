package deserialize

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	apiEvents "github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	anyMarshal "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
	"github.com/buger/jsonparser"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

type protoEventMap map[int32]*anypb.Any

type customFlagParser struct {
	scenario *controlplane.PublicScenario
}

func (cfp *customFlagParser) Get(json []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
		}
	}()

	tasksBytes, _, _, err := jsonparser.Get(json, "tasks")
	if err != nil {
		log.Fatal(err)
	}

	taskIndex := 0
	_, err = jsonparser.ArrayEach(tasksBytes, func(taskBytes []byte, _ jsonparser.ValueType, _ int, _ error) {

		taskCustomFlags := make(protoEventMap)
		cfp.parseCustomFlags(taskBytes, &taskCustomFlags)
		cfp.scenario.Tasks[taskIndex].NodeEvent = &controlplane.Task_CustomFlags{CustomFlags: &apiEvents.CustomFlags{CustomFlags: taskCustomFlags}}

		podConfigsBytes, podConfigsDataType, _, err := jsonparser.Get(taskBytes, "pod_configs")
		if podConfigsDataType == jsonparser.Array {
			if err != nil {
				log.Fatal(err)
			}

			podConfigIndex := 0
			_, err = jsonparser.ArrayEach(podConfigsBytes, func(podConfigBytes []byte, _ jsonparser.ValueType, _ int, _ error) {
				podCustomFlags := make(protoEventMap)
				cfp.parseCustomFlags(podConfigBytes, &podCustomFlags)
				cfp.scenario.Tasks[taskIndex].PodConfigs[podConfigIndex].PodEvent = &controlplane.PodConfig_CustomFlags{CustomFlags: &apiEvents.CustomFlags{CustomFlags: taskCustomFlags}}
				podConfigIndex++
			})
		}

		taskIndex++
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (cfp *customFlagParser) parseCustomFlags(flagBytes []byte, customFlags *protoEventMap) {
	flagBytes, flagDataType, _, err := jsonparser.Get(flagBytes, "custom_flags")

	if flagDataType == jsonparser.Object {
		if err != nil {
			log.Fatal(err)
		}

		err = jsonparser.ObjectEach(flagBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			ef, anyValue := cfp.parseKey(string(key), value)
			(*customFlags)[ef] = anyValue
			return nil
		})
	}
}

func (cfp *customFlagParser) parseKey(key string, value []byte) (events.EventFlag, *anypb.Any) {
	switch key {
	case "node_create_pod_response":
		return events.NodeCreatePodResponse, marshalResponse(value)
	case "node_create_pod_response_percentage":
		return events.NodeCreatePodResponsePercentage, marshalPercent(value)

	case "node_update_pod_response":
		return events.NodeUpdatePodResponse, marshalResponse(value)
	case "node_update_pod_response_percentage":
		return events.NodeUpdatePodResponsePercentage, marshalPercent(value)

	case "node_delete_pod_response":
		return events.NodeDeletePodResponse, marshalResponse(value)
	case "node_delete_pod_response_percentage":
		return events.NodeDeletePodResponsePercentage, marshalPercent(value)

	case "node_get_pod_response":
		return events.NodeGetPodResponse, marshalResponse(value)
	case "node_get_pod_response_percentage":
		return events.NodeGetPodResponsePercentage, marshalPercent(value)

	case "node_get_pod_status_response":
		return events.NodeGetPodStatusResponse, marshalResponse(value)
	case "node_get_pod_status_response_percentage":
		return events.NodeGetPodStatusResponsePercentage, marshalPercent(value)

	case "node_get_pods_response":
		return events.NodeGetPodsResponse, marshalResponse(value)
	case "node_get_pods_response_percentage":
		return events.NodeGetPodsResponsePercentage, marshalPercent(value)

	case "node_ping_response":
		return events.NodePingResponse, marshalResponse(value)
	case "node_ping_response_percentage":
		return events.NodePingResponsePercentage, marshalPercent(value)

	case "node_enable_resource_alteration":
		return events.NodeEnableResourceAlteration, marshalBool(value)
	case "node_memory_usage":
		return events.NodeMemoryUsage, marshalBytes(value, "memory")
	case "node_cpu_usage":
		return events.NodeCPUUsage, marshalInt(value)
	case "node_storage_usage":
		return events.NodeStorageUsage, marshalBytes(value, "storage")
	case "node_ephemeral_storage_usage":
		return events.NodeEphemeralStorageUsage, marshalBytes(value, "ephemeral storage")

	case "node_added_latency_enabled":
		return events.NodeAddedLatencyEnabled, marshalBool(value)
	case "node_added_latency_msec":
		return events.NodeAddedLatencyMsec, marshalInt(value)

	case "pod_create_pod_response":
		return events.PodCreatePodResponse, marshalResponse(value)
	case "pod_create_pod_response_percentage":
		return events.PodCreatePodResponsePercentage, marshalPercent(value)

	case "pod_update_pod_response":
		return events.PodCreatePodResponsePercentage, marshalResponse(value)
	case "pod_update_pod_response_percentage":
		return events.PodUpdatePodResponsePercentage, marshalPercent(value)

	case "pod_delete_pod_response":
		return events.PodDeletePodResponse, marshalResponse(value)
	case "pod_delete_pod_response_percentage":
		return events.PodDeletePodResponsePercentage, marshalPercent(value)

	case "pod_get_pod_response":
		return events.PodGetPodResponse, marshalResponse(value)
	case "pod_get_pod_response_percentage":
		return events.PodGetPodResponsePercentage, marshalPercent(value)

	case "pod_get_pod_status_response":
		return events.PodGetPodStatusResponse, marshalResponse(value)
	case "pod_get_pod_status_response_percentage":
		return events.PodGetPodStatusResponsePercentage, marshalPercent(value)

	case "pod_update_pod_status":
		return events.PodUpdatePodStatus, marshalResponse(value)
	case "pod_update_pod_status_percentage":
		return events.PodUpdatePodStatusPercentage, marshalPodStatus(value)

	default:
		log.Fatalf("invalid key %s", key)
		return 0, nil
	}
}

func marshalResponse(value []byte) *anypb.Any {
	responseStr, err := jsonparser.GetString(value)
	if err != nil {
		log.Fatal(err)
	}

	response := scenario.Response(scenario.Response_value[responseStr])
	return anyMarshal.MarshalOrDie(response)
}

func marshalPodStatus(value []byte) *anypb.Any {
	podStatusStr, err := jsonparser.GetString(value)
	if err != nil {
		log.Fatal(err)
	}

	podStatus := scenario.PodStatus(scenario.PodStatus_value[podStatusStr])
	return anyMarshal.MarshalOrDie(podStatus)
}

func marshalBytes(value []byte, unitName string) *anypb.Any {
	bytes, err := jsonparser.GetString(value)
	if err != nil {
		log.Fatal(err)
	}

	inBytes, err := translate.GetInBytes(bytes, unitName)
	if err != nil {
		log.Fatal(err)
	}

	return anyMarshal.MarshalOrDie(inBytes)
}

func marshalPercent(value []byte) *anypb.Any {
	percent, err := jsonparser.GetInt(value)
	if err != nil {
		log.Fatal(err)
	}

	if percent < 0 || percent > 100 {
		log.Fatal("percentage should be between 0 and 100")
	}

	return anyMarshal.MarshalOrDie(percent)
}

func marshalInt(value []byte) *anypb.Any {
	valueInt, err := jsonparser.GetInt(value)
	if err != nil {
		log.Fatal(err)
	}

	if valueInt < 0 {
		log.Fatal("value should be above 0")
	}

	return anyMarshal.MarshalOrDie(valueInt)
}

func marshalBool(value []byte) *anypb.Any {
	valueBool, err := jsonparser.GetBoolean(value)
	if err != nil {
		log.Fatal(err)
	}

	return anyMarshal.MarshalOrDie(valueBool)
}
