package deserialize

import (
	"errors"
	"log"

	"github.com/buger/jsonparser"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	apiEvents "github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	anyMarshal "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
)

type protoEventMap map[int32]*anypb.Any

type customFlagParser struct {
	scenario *controlplane.PublicScenario
}

func (cfp *customFlagParser) parse(json []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	tasksBytes, _, _, err := jsonparser.Get(json, "tasks")
	if err != nil {
		log.Panic(err)
	}

	taskIndex := 0
	_, err = jsonparser.ArrayEach(tasksBytes, func(taskBytes []byte, _ jsonparser.ValueType, _ int, _ error) {
		currentTask := cfp.scenario.Tasks[taskIndex]

		taskCustomFlags := make(protoEventMap)
		cfp.parseCustomFlags(taskBytes, &taskCustomFlags)

		if hasCustomFlags(taskBytes) {
			currentTask.NodeEvent = &controlplane.Task_CustomFlags{CustomFlags: &apiEvents.CustomFlags{CustomFlags: taskCustomFlags}}
		}

		// Explicitly ignore this error as it's also returned when pod_configs is not set, which we check by the next if
		podConfigsBytes, podConfigsDataType, _, _ := jsonparser.Get(taskBytes, "pod_configs")
		if podConfigsDataType == jsonparser.Array {
			podConfigIndex := 0
			_, err = jsonparser.ArrayEach(podConfigsBytes, func(podConfigBytes []byte, _ jsonparser.ValueType, _ int, _ error) {
				podCustomFlags := make(protoEventMap)
				cfp.parseCustomFlags(podConfigBytes, &podCustomFlags)

				if hasCustomFlags(podConfigBytes) {
					currentTask.PodConfigs[podConfigIndex].PodEvent = &controlplane.PodConfig_CustomFlags{CustomFlags: &apiEvents.CustomFlags{CustomFlags: podCustomFlags}}
				}

				podConfigIndex++
			})

			if err != nil {
				log.Panic(err)
			}
		}

		taskIndex++
	})

	return err
}

func (cfp *customFlagParser) parseCustomFlags(flagBytes []byte, customFlags *protoEventMap) {
	// Explicitly ignore this error as it's also returned when pod_configs is not set, which we check by the next if
	flagBytes, flagDataType, _, _ := jsonparser.Get(flagBytes, "custom_flags")
	if flagDataType == jsonparser.Object {
		err := jsonparser.ObjectEach(flagBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			ef, anyValue := cfp.parseKey(string(key), value)
			(*customFlags)[ef] = anyMarshal.MarshalOrDie(anyValue)
			return nil
		})

		if err != nil {
			log.Panic(err)
		}
	}
}

func (cfp *customFlagParser) parseKey(key string, value []byte) (events.EventFlag, interface{}) {
	switch key {
	case "node_create_pod_response":
		return events.NodeCreatePodResponse, getResponse(value)
	case "node_create_pod_response_percentage":
		return events.NodeCreatePodResponsePercentage, getPercent(value)

	case "node_update_pod_response":
		return events.NodeUpdatePodResponse, getResponse(value)
	case "node_update_pod_response_percentage":
		return events.NodeUpdatePodResponsePercentage, getPercent(value)

	case "node_delete_pod_response":
		return events.NodeDeletePodResponse, getResponse(value)
	case "node_delete_pod_response_percentage":
		return events.NodeDeletePodResponsePercentage, getPercent(value)

	case "node_get_pod_response":
		return events.NodeGetPodResponse, getResponse(value)
	case "node_get_pod_response_percentage":
		return events.NodeGetPodResponsePercentage, getPercent(value)

	case "node_get_pod_status_response":
		return events.NodeGetPodStatusResponse, getResponse(value)
	case "node_get_pod_status_response_percentage":
		return events.NodeGetPodStatusResponsePercentage, getPercent(value)

	case "node_get_pods_response":
		return events.NodeGetPodsResponse, getResponse(value)
	case "node_get_pods_response_percentage":
		return events.NodeGetPodsResponsePercentage, getPercent(value)

	case "node_ping_response":
		return events.NodePingResponse, getResponse(value)
	case "node_ping_response_percentage":
		return events.NodePingResponsePercentage, getPercent(value)

	case "node_enable_resource_alteration":
		return events.NodeEnableResourceAlteration, getBool(value)
	case "node_memory_usage":
		return events.NodeMemoryUsage, getSize(value, "memory")
	case "node_cpu_usage":
		return events.NodeCPUUsage, getIntMinZero(value)
	case "node_storage_usage":
		return events.NodeStorageUsage, getSize(value, "storage")
	case "node_ephemeral_storage_usage":
		return events.NodeEphemeralStorageUsage, getSize(value, "ephemeral storage")

	case "node_added_latency_enabled":
		return events.NodeAddedLatencyEnabled, getBool(value)
	case "node_added_latency_msec":
		return events.NodeAddedLatencyMsec, getIntMinZero(value)

	case "pod_create_pod_response":
		return events.PodCreatePodResponse, getResponse(value)
	case "pod_create_pod_response_percentage":
		return events.PodCreatePodResponsePercentage, getPercent(value)

	case "pod_update_pod_response":
		return events.PodUpdatePodResponse, getResponse(value)
	case "pod_update_pod_response_percentage":
		return events.PodUpdatePodResponsePercentage, getPercent(value)

	case "pod_delete_pod_response":
		return events.PodDeletePodResponse, getResponse(value)
	case "pod_delete_pod_response_percentage":
		return events.PodDeletePodResponsePercentage, getPercent(value)

	case "pod_get_pod_response":
		return events.PodGetPodResponse, getResponse(value)
	case "pod_get_pod_response_percentage":
		return events.PodGetPodResponsePercentage, getPercent(value)

	case "pod_get_pod_status_response":
		return events.PodGetPodStatusResponse, getResponse(value)
	case "pod_get_pod_status_response_percentage":
		return events.PodGetPodStatusResponsePercentage, getPercent(value)

	case "pod_update_pod_status":
		return events.PodUpdatePodStatus, getPodStatus(value)
	case "pod_update_pod_status_percentage":
		return events.PodUpdatePodStatusPercentage, getPercent(value)

	default:
		log.Panicf("invalid custom flag key '%s'", key)
		return 0, nil
	}
}

func hasCustomFlags(value []byte) bool {
	flagBytes, flagDataType, _, _ := jsonparser.Get(value, "custom_flags")
	return flagDataType == jsonparser.Object && len(flagBytes) > 0
}

func getResponse(value []byte) scenario.Response {
	if response, ok := scenario.Response_value[getString(value)]; ok {
		return scenario.Response(response)
	}
	log.Panicf("invalid response '%v'", getString(value))
	return 0
}

func getPodStatus(value []byte) scenario.PodStatus {
	if podStatus, ok := scenario.PodStatus_value[getString(value)]; ok {
		return scenario.PodStatus(podStatus)
	}
	log.Panicf("invalid pod status '%v'", getString(value))
	return 0
}

func getSize(value []byte, unitName string) int64 {
	inBytes, err := translate.GetInBytes(getString(value), unitName)
	if err != nil {
		log.Panic(err)
	}
	return inBytes
}

func getPercent(value []byte) int64 {
	percent := getInt(value)
	if percent < 0 || percent > 100 {
		log.Panic("percentage should be between 0 and 100")
	}
	return percent
}

func getIntMinZero(value []byte) int64 {
	valueInt := getInt(value)
	if valueInt < 0 {
		log.Panic("value should be at least 0")
	}
	return valueInt
}

func getInt(value []byte) int64 {
	valueInt, err := jsonparser.GetInt(value)
	if err != nil {
		log.Panic(err)
	}
	return valueInt
}

func getBool(value []byte) bool {
	valueBool, err := jsonparser.GetBoolean(value)
	if err != nil {
		log.Panic(err)
	}
	return valueBool
}

func getString(value []byte) string {
	return string(value)
}
