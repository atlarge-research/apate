package events

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)

// ResponseState methods
func setCreatePodResponse(responseState *apatelet.ResponseState, response scenario.Response, percentage int32) {
	responseState.CreatePodResponse = response
	responseState.CreatePodResponsePercentage = percentage
}

func setUpdatePodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.UpdatePodResponse = action
	responseState.UpdatePodResponsePercentage = percentage
}

func setDeletePodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.DeletePodResponse = action
	responseState.DeletePodResponsePercentage = percentage
}

func setGetPodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodResponse = action
	responseState.GetPodResponsePercentage = percentage
}

func setGetPodStatusResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodStatusResponse = action
	responseState.GetPodStatusResponsePercentage = percentage
}

// Node specific specific methods
func setGetPodsResponse(responseState *apatelet.NodeState_NodeResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodsResponse = action
	responseState.GetPodsResponsePercentage = percentage
}

func setPingResponse(responseState *apatelet.NodeState_NodeResponseState, action scenario.Response, percentage int32) {
	responseState.PingResponse = action
	responseState.PingResponsePercentage = percentage
}

// getNodeResponseState returns the NodeResponseState of the nodeEvent
func getNodeResponseState(nodeEvent *apatelet.NodeEvent) *apatelet.NodeState_NodeResponseState {
	return nodeEvent.GetNodeState().GetNodeResponseState()
}

// getPodResponseState returns the PodResponseState of the podEvent
func getPodResponseState(podEvent *apatelet.PodEvent) *apatelet.PodState_PodResponseState {
	return podEvent.GetPodState().GetPodResponseState()
}

// applyNodeResponse applies an action and percentage to all lifecycle states on the nodeEvent
func applyNodeResponse(nodeEvent *apatelet.NodeEvent, action scenario.Response, percentage int32) {
	nodeResponseState := nodeEvent.GetNodeState().GetNodeResponseState()
	setGetPodsResponse(nodeResponseState, action, percentage)
	setPingResponse(nodeResponseState, action, percentage)

	responseState := nodeResponseState.GetResponseState()
	setCreatePodResponse(responseState, action, percentage)
	setUpdatePodResponse(responseState, action, percentage)
	setDeletePodResponse(responseState, action, percentage)
	setGetPodResponse(responseState, action, percentage)
	setGetPodStatusResponse(responseState, action, percentage)
}
