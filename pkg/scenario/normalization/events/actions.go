package events

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)


type LifecycleGetter interface {
	getLifecycleState() *apatelet.LifecycleState
}

type BaseEventWrapper struct {
	lifecycleGetter LifecycleGetter
}

func (n *BaseEventWrapper) setCreatePodAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := n.lifecycleGetter.getLifecycleState()
	lifecycleState.CreatePodAction = action
	lifecycleState.CreatePodActionPercentage = percentage
}

func (n *BaseEventWrapper) setUpdatePodAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := n.lifecycleGetter.getLifecycleState()
	lifecycleState.UpdatePodAction = action
	lifecycleState.UpdatePodActionPercentage = percentage
}

func (n *BaseEventWrapper) setDeletePodAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := n.lifecycleGetter.getLifecycleState()
	lifecycleState.DeletePodAction = action
	lifecycleState.DeletePodActionPercentage = percentage
}

func (n *BaseEventWrapper) setGetPodAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := n.lifecycleGetter.getLifecycleState()
	lifecycleState.GetPodAction = action
	lifecycleState.GetPodActionPercentage = percentage
}

func (n *BaseEventWrapper) setGetPodStatusAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := n.lifecycleGetter.getLifecycleState()
	lifecycleState.GetPodStatusAction = action
	lifecycleState.GetPodStatusActionPercentage = percentage
}

// Node
type NodeEventWrapper struct {
	BaseEventWrapper
	nodeEvent *apatelet.NodeEvent
}

func (ne *NodeEventWrapper) getLifecycleState() *apatelet.LifecycleState {
	return ne.nodeEvent.GetNodeState().GetNodeLifecycleState().GetLifecycleState()
}

func (ne *NodeEventWrapper) getAddedLatencyState() *apatelet.LifecycleState {
	return ne.nodeEvent.GetNodeState().GetNodeLifecycleState().GetLifecycleState()
}

func (ne *NodeEventWrapper) applyAction(action scenario.LifecycleAction, percentage int32) {
	ne.setCreatePodAction(action, percentage)
	ne.setUpdatePodAction(action, percentage)
	ne.setDeletePodAction(action, percentage)
	ne.setGetPodAction(action, percentage)
	ne.setGetPodStatusAction(action, percentage)
	ne.setGetPodsAction(action, percentage)
	ne.setPingAction(action, percentage)
}

func (ne *NodeEventWrapper) setGetPodsAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := ne.nodeEvent.GetNodeState().GetNodeLifecycleState()
	lifecycleState.GetPodsAction = action
	lifecycleState.GetPodsActionPercentage = percentage
}

func (ne *NodeEventWrapper) setPingAction(action scenario.LifecycleAction, percentage int32) {
	lifecycleState := ne.nodeEvent.GetNodeState().GetNodeLifecycleState()
	lifecycleState.PingAction = action
	lifecycleState.PingActionPercentage = percentage
}

// Pod
type PodEventWrapper struct {
	BaseEventWrapper
	podEvent *apatelet.PodEvent
}

func (n *PodEventWrapper) getLifecycleState() *apatelet.LifecycleState {
	return n.podEvent.GetPodState().GetPodLifecycleState().GetLifecycleState()
}

