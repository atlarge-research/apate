// Package events contains utilities to translate events of tasks between public API and internal API formats
package events

import (
	"errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
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

// createNodeEvent creates a new node event and hooks it up to the current task
func (et *EventTranslator) createNodeEvent() *apatelet.NodeEvent {
	ne := &apatelet.NodeEvent{
		NodeState: &apatelet.NodeState{
			NodeResponseState: &apatelet.NodeState_NodeResponseState{
				ResponseState: &apatelet.ResponseState{},
			},
			ResourceState:     &apatelet.NodeState_ResourceState{},
			AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
		},
	}
	et.newTask.Event = &apatelet.Task_NodeEvent{NodeEvent: ne}
	return ne
}

// createPodEvent creates a new pod event and hooks it up to the current task
func (et *EventTranslator) createPodEvent() *apatelet.PodEvent {
	pe := &apatelet.PodEvent{
		PodState: &apatelet.PodState{
			PodResponseState: &apatelet.PodState_PodResponseState{
				ResponseState: &apatelet.ResponseState{},
			},
		},
	}
	et.newTask.Event = &apatelet.Task_PodEvent{PodEvent: pe}
	return pe
}

// TranslateEvent translates events sent through the public api to events understood by the Apatelets
func (et *EventTranslator) TranslateEvent() error {
	if et.originalTask.Event == nil {
		return errors.New("you must pass an event to be executed")
	}

	// et.originalTask.Event can be one of many types (see generated protobuf code)
	// taskEvent will be the cast version of this event to the corresponding event, depending on the case
	switch taskEvent := et.originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		applyNodeResponse(et.createNodeEvent(), scenario.Response_TIMEOUT, 100)

	case *controlplane.Task_NetworkLatency:
		latencyState := et.createNodeEvent().GetNodeState().GetAddedLatencyState()
		latencyState.AddedLatencyEnabled = true

		if taskEvent.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		latencyState.AddedLatencyMsec = taskEvent.NetworkLatency.LatencyMsec

	case *controlplane.Task_TimeoutKeepHeartbeat:
		ne := et.createNodeEvent()
		applyNodeResponse(ne, scenario.Response_TIMEOUT, 100)
		setPingResponse(getNodeResponseState(ne), scenario.Response_NORMAL, 0)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		ne := et.createNodeEvent()
		setPingResponse(getNodeResponseState(ne), scenario.Response_TIMEOUT, 100)

	case *controlplane.Task_NodeResponseState:
		ne := et.createNodeEvent()
		nodeResponseState := getNodeResponseState(ne)
		ResponseState := nodeResponseState.GetResponseState()

		state := taskEvent.NodeResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.RequestType_CREATE_POD:
			setCreatePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_UPDATE_POD:
			setUpdatePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_DELETE_POD:
			setDeletePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_GET_POD:
			setGetPodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_GET_POD_STATUS:
			setGetPodStatusResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_GET_PODS:
			setGetPodsResponse(nodeResponseState, state.Response, state.Percentage)

		case events.RequestType_PING:
			setPingResponse(nodeResponseState, state.Response, state.Percentage)
		}

	case *controlplane.Task_ResourcePressure:
		resourceState := et.createNodeEvent().GetNodeState().GetResourceState()
		resourceState.EnableResourceAlteration = true

		rp := taskEvent.ResourcePressure

		if rp.CpuUsage < 0 {
			return errors.New("CPU usage should be at least 0")
		}
		resourceState.CpuUsage = rp.CpuUsage

		memory, err := GetInBytes(rp.MemoryUsage, "memory")
		if err != nil {
			return err
		}
		resourceState.MemoryUsage = memory

		storage, err := GetInBytes(rp.StorageUsage, "storage")
		if err != nil {
			return err
		}
		resourceState.StorageUsage = storage

		ephStorage, err := GetInBytes(rp.EphemeralStorageUsage, "ephemeral storage")
		if err != nil {
			return err
		}
		resourceState.EphemeralStorageUsage = ephStorage

	// Pod events
	case *controlplane.Task_PodResponseState:
		pe := et.createPodEvent()
		ResponseState := getPodResponseState(pe).GetResponseState()

		state := taskEvent.PodResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.RequestType_CREATE_POD:
			setCreatePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_UPDATE_POD:
			setUpdatePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_DELETE_POD:
			setDeletePodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_GET_POD:
			setGetPodResponse(ResponseState, state.Response, state.Percentage)

		case events.RequestType_GET_POD_STATUS:
			setGetPodStatusResponse(ResponseState, state.Response, state.Percentage)

		default:
			return errors.New("can't alter the GetPods / Ping response on pod level")
		}

	case *controlplane.Task_PodStatusUpdate:
		event := et.createPodEvent()

		if taskEvent.PodStatusUpdate.Percentage < 0 || taskEvent.PodStatusUpdate.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		event.GetPodState().PodStatus = taskEvent.PodStatusUpdate.NewStatus
		event.GetPodState().PodStatusPercentage = taskEvent.PodStatusUpdate.Percentage

	case *controlplane.Task_PodStartTimeUpdate:
		et.createPodEvent().GetPodState().PodStartTime = taskEvent.PodStartTimeUpdate.NewStartTime
	}

	return nil
}
