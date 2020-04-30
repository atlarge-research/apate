package events

import (
	"errors"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
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
			NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{},
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
			PodLifecycleState: &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{},
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
		applyNodeAction(th.createNodeEvent(), scenario.LifecycleAction_TIMEOUT, 100)

	case *controlplane.Task_NetworkLatency:
		latencyState := th.createNodeEvent().GetNodeState().GetAddedLatencyState()
		latencyState.AddedLatencyEnabled = true

		if x.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		latencyState.AddedLatencyMsec = x.NetworkLatency.LatencyMsec

	case *controlplane.Task_TimeoutKeepHeartbeat:
		ne := th.createNodeEvent()
		applyNodeAction(ne, scenario.LifecycleAction_TIMEOUT, 100)
		setPingAction(getNodeLifecycleState(ne), scenario.LifecycleAction_NORMAL, 0)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		ne := th.createNodeEvent()
		setPingAction(getNodeLifecycleState(ne), scenario.LifecycleAction_TIMEOUT, 100)

	case *controlplane.Task_NodeLifecycleState:
		ne := th.createNodeEvent()
		nodeLifecycleState := getNodeLifecycleState(ne)
		lifecycleState := nodeLifecycleState.GetLifecycleState()

		state := x.NodeLifecycleState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.LifecycleType_CREATE_POD:
			setCreatePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_UPDATE_POD:
			setUpdatePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_DELETE_POD:
			setDeletePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_GET_POD:
			setGetPodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_GET_POD_STATUS:
			setGetPodStatusAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_GET_PODS:
			setGetPodsAction(nodeLifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_PING:
			setPingAction(nodeLifecycleState, state.Action, state.Percentage)
		}

	case *controlplane.Task_ResourcePressure:
		resourceState := th.createNodeEvent().GetNodeState().GetResourceState()
		resourceState.EnableResourceAlteration = true

		rp := x.ResourcePressure

		if rp.CpuUsage < 0 {
			return errors.New("CPU usage should be at least 0")
		}
		resourceState.CpuUsage = rp.CpuUsage

		if rp.MemoryUsage < 0 {
			return errors.New("memory usage should be at least 0")
		}
		resourceState.MemoryUsage = rp.MemoryUsage

		if rp.StorageUsage < 0 {
			return errors.New("storage usage should be at least 0")
		}
		resourceState.StorageUsage = rp.StorageUsage

		if rp.EphemeralStorageUsage < 0 {
			return errors.New("ephemeral storage usage should be at least 0")
		}
		resourceState.EphemeralStorageUsage = rp.EphemeralStorageUsage

	// Pod events
	case *controlplane.Task_PodLifecycleState:
		pe := th.createPodEvent()
		lifecycleState := getPodLifecycleState(pe).GetLifecycleState()

		state := x.PodLifecycleState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New("percentage should be between 0 and 100")
		}

		switch state.Type {
		case events.LifecycleType_CREATE_POD:
			setCreatePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_UPDATE_POD:
			setUpdatePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_DELETE_POD:
			setDeletePodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_GET_POD:
			setGetPodAction(lifecycleState, state.Action, state.Percentage)

		case events.LifecycleType_GET_POD_STATUS:
			setGetPodStatusAction(lifecycleState, state.Action, state.Percentage)

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
