package events

import (
	"errors"
	"fmt"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/docker/go-units"
)

type EventTranslator struct {
	originalTask *controlplane.Task
	newTask      *apatelet.Task
}

func NewEventTranslator(originalTask *controlplane.Task, newTask *apatelet.Task) *EventTranslator {
	return &EventTranslator{
		originalTask: originalTask,
		newTask:      newTask,
	}
}

func (th *EventTranslator) createNodeEvent() *apatelet.NodeEvent {
	ne := &apatelet.NodeEvent{
		NodeState: &apatelet.NodeState{
			NodeResponseState: &apatelet.NodeState_NodeResponseState{
				ResponseState: &apatelet.ResponseState{},
			},
			ResourceState:     &apatelet.NodeState_ResourceState{},
			AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
		},
	}
	th.newTask.Event = &apatelet.Task_NodeEvent{NodeEvent: ne}
	return ne
}

func (th *EventTranslator) createPodEvent() *apatelet.PodEvent {
	pe := &apatelet.PodEvent{
		PodState: &apatelet.PodState{
			PodResponseState: &apatelet.PodState_PodResponseState{
				ResponseState: &apatelet.ResponseState{},
			},
		},
	}
	th.newTask.Event = &apatelet.Task_PodEvent{PodEvent: pe}
	return pe
}

func (th *EventTranslator) TranslateEvent() error {
	if th.originalTask.Event == nil {
		return errors.New("you must pass an event to be executed")
	}

	switch x := th.originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		applyNodeResponse(th.createNodeEvent(), scenario.Response_TIMEOUT, 100)

	case *controlplane.Task_NetworkLatency:
		latencyState := th.createNodeEvent().GetNodeState().GetAddedLatencyState()
		latencyState.AddedLatencyEnabled = true

		if x.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		latencyState.AddedLatencyMsec = x.NetworkLatency.LatencyMsec

	case *controlplane.Task_TimeoutKeepHeartbeat:
		ne := th.createNodeEvent()
		applyNodeResponse(ne, scenario.Response_TIMEOUT, 100)
		setPingResponse(getNodeResponseState(ne), scenario.Response_NORMAL, 0)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		ne := th.createNodeEvent()
		setPingResponse(getNodeResponseState(ne), scenario.Response_TIMEOUT, 100)

	case *controlplane.Task_NodeResponseState:
		ne := th.createNodeEvent()
		nodeResponseState := getNodeResponseState(ne)
		ResponseState := nodeResponseState.GetResponseState()

		state := x.NodeResponseState

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
		resourceState := th.createNodeEvent().GetNodeState().GetResourceState()
		resourceState.EnableResourceAlteration = true

		rp := x.ResourcePressure

		if rp.CpuUsage < 0 {
			return errors.New("CPU usage should be at least 0")
		}
		resourceState.CpuUsage = rp.CpuUsage

		memory, err := getInBytes(rp.MemoryUsage, "memory")
		if err != nil {
			return err
		}
		resourceState.MemoryUsage = memory

		storage, err := getInBytes(rp.StorageUsage, "storage")
		if err != nil {
			return err
		}
		resourceState.StorageUsage = storage

		ephStorage, err := getInBytes(rp.EphemeralStorageUsage, "ephemeral storage")
		if err != nil {
			return err
		}
		resourceState.EphemeralStorageUsage = ephStorage

	// Pod events
	case *controlplane.Task_PodResponseState:
		pe := th.createPodEvent()
		ResponseState := getPodResponseState(pe).GetResponseState()

		state := x.PodResponseState

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
		event := th.createPodEvent()

		if x.PodStatusUpdate.Percentage < 0 || x.PodStatusUpdate.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		event.GetPodState().PodStatus = x.PodStatusUpdate.NewStatus
		event.GetPodState().PodStatusPercentage = x.PodStatusUpdate.Percentage

	case *controlplane.Task_PodStartTimeUpdate:
		th.createPodEvent().GetPodState().PodStartTime = x.PodStartTimeUpdate.NewStartTime
	}

	return nil
}

func getInBytes(unit string, unitName string) (int64, error) {
	unitInt, err := units.RAMInBytes(unit)
	if err != nil {
		return 0, err
	} else if unitInt < 0 {
		return 0, fmt.Errorf("%s usage should be at least 0", unitName)
	}
	return unitInt, nil
}
