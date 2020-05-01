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
	if et.originalTask.Event == nil {
		return errors.New("you must pass an event to be executed")
	}

	f := newEventFlags()

	// et.originalTask.Event can be one of many types (see generated protobuf code)
	// taskEvent will be the cast version of this event to the corresponding event, depending on the case
	switch taskEvent := et.originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		f.flags(scenario.Response_TIMEOUT, nodeEventFlags)
		f.flags(100, nodeEventPercentageFlags)

	case *controlplane.Task_NetworkLatency:
		if taskEvent.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		f.flag(true, ef.NodeAddedLatencyEnabled)
		f.flag(taskEvent.NetworkLatency.LatencyMsec, ef.NodeAddedLatencyMsec)

	case *controlplane.Task_TimeoutKeepHeartbeat:
		f.flags(scenario.Response_TIMEOUT, nodeEventFlags)
		f.flags(100, nodeEventPercentageFlags)

		// Reset ping
		f.flag(scenario.Response_NORMAL, ef.NodePingResponse)
		f.flag(0, ef.NodePingResponsePercentage)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		f.flag(scenario.Response_TIMEOUT, ef.NodePingResponse)
		f.flag(100, ef.NodePingResponsePercentage)

	case *controlplane.Task_NodeResponseState:
		state := taskEvent.NodeResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.RequestType_CREATE_POD:
			f.flag(state.Response, ef.NodeCreatePodResponse)
			f.flag(state.Percentage, ef.NodeCreatePodResponsePercentage)

		case events.RequestType_UPDATE_POD:
			f.flag(state.Response, ef.NodeUpdatePodResponse)
			f.flag(state.Percentage, ef.NodeUpdatePodResponsePercentage)

		case events.RequestType_DELETE_POD:
			f.flag(state.Response, ef.NodeDeletePodResponse)
			f.flag(state.Percentage, ef.NodeDeletePodResponsePercentage)

		case events.RequestType_GET_POD:
			f.flag(state.Response, ef.NodeGetPodResponse)
			f.flag(state.Percentage, ef.NodeGetPodResponsePercentage)

		case events.RequestType_GET_POD_STATUS:
			f.flag(state.Response, ef.NodeGetPodStatusResponse)
			f.flag(state.Percentage, ef.NodeGetPodStatusResponsePercentage)

		case events.RequestType_GET_PODS:
			f.flag(state.Response, ef.NodeGetPodsResponse)
			f.flag(state.Percentage, ef.NodeGetPodsResponsePercentage)

		case events.RequestType_PING:
			f.flag(state.Response, ef.NodePingResponse)
			f.flag(state.Percentage, ef.NodePingResponsePercentage)
		}

	case *controlplane.Task_ResourcePressure:
		f.flag(true, ef.NodeEnableResourceAlteration)

		rp := taskEvent.ResourcePressure

		if rp.CpuUsage < 0 {
			return errors.New("CPU usage should be at least 0")
		}
		f.flag(rp.CpuUsage, ef.NodeCPUUsage)

		memory, err := GetInBytes(rp.MemoryUsage, "memory")
		if err != nil {
			return err
		}
		f.flag(memory, ef.NodeMemoryUsage)

		storage, err := GetInBytes(rp.StorageUsage, "storage")
		if err != nil {
			return err
		}
		f.flag(storage, ef.NodeStorageUsage)

		ephStorage, err := GetInBytes(rp.EphemeralStorageUsage, "ephemeral storage")
		if err != nil {
			return err
		}
		f.flag(ephStorage, ef.NodeEphemeralStorageUsage)

	// Pod events
	case *controlplane.Task_PodResponseState:
		state := taskEvent.PodResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.RequestType_CREATE_POD:
			f.flag(state.Response, ef.PodCreatePodResponse)
			f.flag(state.Percentage, ef.PodCreatePodResponsePercentage)

		case events.RequestType_UPDATE_POD:
			f.flag(state.Response, ef.PodUpdatePodResponse)
			f.flag(state.Percentage, ef.PodUpdatePodResponsePercentage)

		case events.RequestType_DELETE_POD:
			f.flag(state.Response, ef.PodDeletePodResponse)
			f.flag(state.Percentage, ef.PodDeletePodResponsePercentage)

		case events.RequestType_GET_POD:
			f.flag(state.Response, ef.PodGetPodResponse)
			f.flag(state.Percentage, ef.PodGetPodResponsePercentage)

		case events.RequestType_GET_POD_STATUS:
			f.flag(state.Response, ef.PodGetPodStatusResponse)
			f.flag(state.Percentage, ef.PodGetPodStatusResponsePercentage)

		default:
			return errors.New("can't alter the GetPods / Ping response on pod level")
		}

	case *controlplane.Task_PodStatusUpdate:
		if taskEvent.PodStatusUpdate.Percentage < 0 || taskEvent.PodStatusUpdate.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		f.flag(taskEvent.PodStatusUpdate.NewStatus, ef.PodUpdatePodStatus)
		f.flag(taskEvent.PodStatusUpdate.Percentage, ef.PodUpdatePodStatusPercentage)
	}

	et.newTask.EventFlags = f

	return nil
}
