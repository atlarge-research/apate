package events

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)

// LifecycleState methods
func setCreatePodAction(lifecycleState *apatelet.LifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.CreatePodAction = action
	lifecycleState.CreatePodActionPercentage = percentage
}

func setUpdatePodAction(lifecycleState *apatelet.LifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.UpdatePodAction = action
	lifecycleState.UpdatePodActionPercentage = percentage
}

func setDeletePodAction(lifecycleState *apatelet.LifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.DeletePodAction = action
	lifecycleState.DeletePodActionPercentage = percentage
}

func setGetPodAction(lifecycleState *apatelet.LifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.GetPodAction = action
	lifecycleState.GetPodActionPercentage = percentage
}

func setGetPodStatusAction(lifecycleState *apatelet.LifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.GetPodStatusAction = action
	lifecycleState.GetPodStatusActionPercentage = percentage
}

// NodeLifecycleState specific methods
func setGetPodsAction(lifecycleState *apatelet.NodeState_NodeLifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.GetPodsAction = action
	lifecycleState.GetPodsActionPercentage = percentage
}

func setPingAction(lifecycleState *apatelet.NodeState_NodeLifecycleState, action scenario.LifecycleAction, percentage int32) {
	lifecycleState.PingAction = action
	lifecycleState.PingActionPercentage = percentage
}

// getNodeLifecycleState returns the NodeLifecycleState of the nodeEvent
func getNodeLifecycleState(nodeEvent *apatelet.NodeEvent) *apatelet.NodeState_NodeLifecycleState {
	return nodeEvent.GetNodeState().GetNodeLifecycleState()
}

// getPodLifecycleState returns the PodLifecycleState of the podEvent
func getPodLifecycleState(podEvent *apatelet.PodEvent) *apatelet.PodState_PodLifecycleState {
	return podEvent.GetPodState().GetPodLifecycleState()
}

// applyNodeAction applies an action and percentage to all lifecycle states on the nodeEvent
func applyNodeAction(nodeEvent *apatelet.NodeEvent, action scenario.LifecycleAction, percentage int32) {
	nodeLifecycleState := nodeEvent.GetNodeState().GetNodeLifecycleState()
	setGetPodsAction(nodeLifecycleState, action, percentage)
	setPingAction(nodeLifecycleState, action, percentage)

	lifecycleState := nodeLifecycleState.GetLifecycleState()
	setCreatePodAction(lifecycleState, action, percentage)
	setUpdatePodAction(lifecycleState, action, percentage)
	setDeletePodAction(lifecycleState, action, percentage)
	setGetPodAction(lifecycleState, action, percentage)
	setGetPodStatusAction(lifecycleState, action, percentage)
}
