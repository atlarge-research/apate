// Package translate contains utilities to translate events of tasks between public API and internal API formats
package translate

import (
	"errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	ef "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// EventTranslator is a utility to translate events between events sent through the public api to events understood by the Apatelets
type EventTranslator struct {
	originalTask *controlplane.Task
	newTask      *apatelet.Task
}

// NewEventTranslator constructs a new EventTranslator
func NewEventTranslator(originalTask *controlplane.Task, newTask *apatelet.Task) *EventTranslator {
	return &EventTranslator{
		originalTask: originalTask,
		newTask:      newTask,
	}
}

// TranslateEvent translates events sent through the public api to events understood by the Apatelets
func (et *EventTranslator) TranslateEvent() (err error) {
	err = et.translateNodeEventFlags()
	if err != nil {
		return err
	}

	err = et.translatePodEventFlags()
	if err != nil {
		return err
	}

	return nil
}

func (et *EventTranslator) translatePodEventFlags() error {
	if et.originalTask.PodConfigs == nil {
		return nil
	}

	for _, podConfig := range et.originalTask.PodConfigs {
		pef := newEventFlags()

		switch pe := podConfig.PodEvent.(type) {
		case *controlplane.PodConfig_PodResponseState:
			state := pe.PodResponseState

			if state.Percentage < 0 || state.Percentage > 100 {
				return errors.New("percentage should be between 0 and 100")
			}

			switch state.Type {
			case events.RequestType_CREATE_POD:
				pef.flag(state.Response, ef.PodCreatePodResponse)
				pef.flag(state.Percentage, ef.PodCreatePodResponsePercentage)

			case events.RequestType_UPDATE_POD:
				pef.flag(state.Response, ef.PodUpdatePodResponse)
				pef.flag(state.Percentage, ef.PodUpdatePodResponsePercentage)

			case events.RequestType_DELETE_POD:
				pef.flag(state.Response, ef.PodDeletePodResponse)
				pef.flag(state.Percentage, ef.PodDeletePodResponsePercentage)

			case events.RequestType_GET_POD:
				pef.flag(state.Response, ef.PodGetPodResponse)
				pef.flag(state.Percentage, ef.PodGetPodResponsePercentage)

			case events.RequestType_GET_POD_STATUS:
				pef.flag(state.Response, ef.PodGetPodStatusResponse)
				pef.flag(state.Percentage, ef.PodGetPodStatusResponsePercentage)

			default:
				return errors.New("can't alter the GetPods / Ping response on pod level")
			}

		case *controlplane.PodConfig_PodStatusUpdate:
			if pe.PodStatusUpdate.Percentage < 0 || pe.PodStatusUpdate.Percentage > 100 {
				return errors.New("percentage should be between 0 and 100")
			}

			pef.flag(pe.PodStatusUpdate.NewStatus, ef.PodUpdatePodStatus)
			pef.flag(pe.PodStatusUpdate.Percentage, ef.PodUpdatePodStatusPercentage)
		}

		podConfig := &apatelet.PodConfig{
			MetadataName: podConfig.MetadataName,
			EventFlags:   pef,
		}
		et.newTask.PodConfigs = append(et.newTask.PodConfigs, podConfig)
	}
	return nil
}

func (et *EventTranslator) translateNodeEventFlags() error {
	nef := newEventFlags()

	if et.originalTask.NodeEvent == nil {
		return nil
	}

	// et.originalTask.Event can be one of many types (see generated protobuf code)
	// ne will be the cast version of this event to the corresponding event, depending on the case
	switch ne := et.originalTask.NodeEvent.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		nef.flags(scenario.Response_TIMEOUT, nodeEventFlags)
		nef.flags(100, nodeEventPercentageFlags)

	case *controlplane.Task_NetworkLatency:
		if ne.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		nef.flag(true, ef.NodeAddedLatencyEnabled)
		nef.flag(ne.NetworkLatency.LatencyMsec, ef.NodeAddedLatencyMsec)

	case *controlplane.Task_TimeoutKeepHeartbeat:
		nef.flags(scenario.Response_TIMEOUT, nodeEventFlags)
		nef.flags(100, nodeEventPercentageFlags)

		// Reset ping
		nef.flag(scenario.Response_NORMAL, ef.NodePingResponse)
		nef.flag(0, ef.NodePingResponsePercentage)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		nef.flag(scenario.Response_TIMEOUT, ef.NodePingResponse)
		nef.flag(100, ef.NodePingResponsePercentage)

	case *controlplane.Task_NodeResponseState:
		state := ne.NodeResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.RequestType_CREATE_POD:
			nef.flag(state.Response, ef.NodeCreatePodResponse)
			nef.flag(state.Percentage, ef.NodeCreatePodResponsePercentage)

		case events.RequestType_UPDATE_POD:
			nef.flag(state.Response, ef.NodeUpdatePodResponse)
			nef.flag(state.Percentage, ef.NodeUpdatePodResponsePercentage)

		case events.RequestType_DELETE_POD:
			nef.flag(state.Response, ef.NodeDeletePodResponse)
			nef.flag(state.Percentage, ef.NodeDeletePodResponsePercentage)

		case events.RequestType_GET_POD:
			nef.flag(state.Response, ef.NodeGetPodResponse)
			nef.flag(state.Percentage, ef.NodeGetPodResponsePercentage)

		case events.RequestType_GET_POD_STATUS:
			nef.flag(state.Response, ef.NodeGetPodStatusResponse)
			nef.flag(state.Percentage, ef.NodeGetPodStatusResponsePercentage)

		case events.RequestType_GET_PODS:
			nef.flag(state.Response, ef.NodeGetPodsResponse)
			nef.flag(state.Percentage, ef.NodeGetPodsResponsePercentage)

		case events.RequestType_PING:
			nef.flag(state.Response, ef.NodePingResponse)
			nef.flag(state.Percentage, ef.NodePingResponsePercentage)
		}

	case *controlplane.Task_ResourcePressure:
		nef.flag(true, ef.NodeEnableResourceAlteration)

		rp := ne.ResourcePressure

		if rp.CpuUsage < 0 {
			return errors.New("CPU usage should be at least 0")
		}
		nef.flag(rp.CpuUsage, ef.NodeCPUUsage)

		memory, err := GetInBytes(rp.MemoryUsage, "memory")
		if err != nil {
			return err
		}
		nef.flag(memory, ef.NodeMemoryUsage)

		storage, err := GetInBytes(rp.StorageUsage, "storage")
		if err != nil {
			return err
		}
		nef.flag(storage, ef.NodeStorageUsage)

		ephStorage, err := GetInBytes(rp.EphemeralStorageUsage, "ephemeral storage")
		if err != nil {
			return err
		}
		nef.flag(ephStorage, ef.NodeEphemeralStorageUsage)
	default:
		return errors.New("unknown node event")
	}

	et.newTask.NodeEventFlags = nef

	return nil
}
