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

func (th *EventTranslator) createNodeEvent() *NodeEventWrapper {
	nodeEvent := &apatelet.NodeEvent{
		NodeState: &apatelet.NodeState{
			NodeLifecycleState: &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{},
			},
			ResourceState:     &apatelet.NodeState_ResourceState{},
			AddedLatencyState: &apatelet.NodeState_AddedLatencyState{},
		},
	}
	th.newTask.Event = &apatelet.Task_NodeEvent{NodeEvent: nodeEvent}
	n := &NodeEventWrapper{nodeEvent: nodeEvent}
	n.BaseEventWrapper = BaseEventWrapper{lifecycleGetter: n}
	return n
}

func (th *EventTranslator) createPodEvent() *PodEventWrapper {
	podEvent := &apatelet.PodEvent{
		PodState: &apatelet.PodState{
			PodLifecycleState: &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{},
			},
		},
	}
	th.newTask.Event = &apatelet.Task_PodEvent{PodEvent: podEvent}
	p := &PodEventWrapper{podEvent: podEvent}
	p.BaseEventWrapper = BaseEventWrapper{lifecycleGetter: p}
	return p
}

func (th *EventTranslator) TranslateEvent() error {
	if th.originalTask.Event == nil {
		return errors.New("you must pass an event to be executed")
	}

	switch x := th.originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		ne := th.createNodeEvent()
		ne.applyAction(scenario.LifecycleAction_TIMEOUT, 100)

	case *controlplane.Task_NetworkLatency:
		latencyState := th.createNodeEvent().nodeEvent.GetNodeState().GetAddedLatencyState()
		latencyState.AddedLatencyEnabled = true
		latencyState.AddedLatencyMsec = x.NetworkLatency.GetLatencyMsec()

	case *controlplane.Task_TimeoutKeepHeartbeat:
		ne := th.createNodeEvent()
		ne.applyAction(scenario.LifecycleAction_TIMEOUT, 100)
		ne.setPingAction(scenario.LifecycleAction_NORMAL, 0)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		ne := th.createNodeEvent()
		ne.setPingAction(scenario.LifecycleAction_TIMEOUT, 100)

	case *controlplane.Task_NodeLifecycleState:
		ne := th.createNodeEvent()
		state := x.NodeLifecycleState

		switch state.Type {
		case events.LifecycleType_CREATE_POD:
			ne.setCreatePodAction(state.Action, state.Percentage)

		case events.LifecycleType_UPDATE_POD:
			ne.setUpdatePodAction(state.Action, state.Percentage)

		case events.LifecycleType_DELETE_POD:
			ne.setDeletePodAction(state.Action, state.Percentage)

		case events.LifecycleType_GET_POD:
			ne.setGetPodAction(state.Action, state.Percentage)

		case events.LifecycleType_GET_POD_STATUS:
			ne.setGetPodStatusAction(state.Action, state.Percentage)

		case events.LifecycleType_GET_PODS:
			ne.setGetPodsAction(state.Action, state.Percentage)

		case events.LifecycleType_PING:
			ne.setPingAction(state.Action, state.Percentage)
		}

	case *controlplane.Task_ResourcePressure:
		resourceState := th.createNodeEvent().nodeEvent.GetNodeState().GetResourceState()
		resourceState.EnableResourceAlteration = true
		resourceState.CpuUsage = x.ResourcePressure.GetCpuUsage()
		resourceState.MemoryUsage = x.ResourcePressure.GetMemoryUsage()
		resourceState.StorageUsage = x.ResourcePressure.GetStorageUsage()
		resourceState.EphemeralStorageUsage = x.ResourcePressure.GetEphemeralStorageUsage()

		// Pod events
	case *controlplane.Task_PodLifecycleState:
		pe := th.createPodEvent()
		state := x.PodLifecycleState

		switch state.Type {
		case events.LifecycleType_CREATE_POD:
			pe.setCreatePodAction(state.Action, state.Percentage)

		case events.LifecycleType_UPDATE_POD:
			pe.setUpdatePodAction(state.Action, state.Percentage)

		case events.LifecycleType_DELETE_POD:
			pe.setDeletePodAction(state.Action, state.Percentage)

		case events.LifecycleType_GET_POD:
			pe.setGetPodAction(state.Action, state.Percentage)

		case events.LifecycleType_GET_POD_STATUS:
			pe.setGetPodStatusAction(state.Action, state.Percentage)

		default:
			return errors.New("can't alter the GetPods / Ping response on pod level")
		}

	case *controlplane.Task_PodStatusUpdate:
		th.createPodEvent().podEvent.GetPodState().PodStatus = x.PodStatusUpdate.NewStatus

	case *controlplane.Task_PodStartTimeUpdate:
		th.createPodEvent().podEvent.GetPodState().StartTime = x.PodStartTimeUpdate.NewStartTime
	}

	return nil
}
