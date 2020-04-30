package events

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)

// ResponseState methods
func setCreatePodResponse(ResponseState *apatelet.ResponseState, response scenario.Response, percentage int32) {
	ResponseState.CreatePodResponse = response
	ResponseState.CreatePodResponsePercentage = percentage
}

func setUpdatePodResponse(ResponseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	ResponseState.UpdatePodResponse = action
	ResponseState.UpdatePodResponsePercentage = percentage
}

func setDeletePodResponse(ResponseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	ResponseState.DeletePodResponse = action
	ResponseState.DeletePodResponsePercentage = percentage
}

func setGetPodResponse(ResponseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	ResponseState.GetPodResponse = action
	ResponseState.GetPodResponsePercentage = percentage
}

func setGetPodStatusResponse(ResponseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	ResponseState.GetPodStatusResponse = action
	ResponseState.GetPodStatusResponsePercentage = percentage
}

// NodeResponseState specific methods
func setGetPodsResponse(ResponseState *apatelet.NodeState_NodeResponseState, action scenario.Response, percentage int32) {
	ResponseState.GetPodsResponse = action
	ResponseState.GetPodsResponsePercentage = percentage
}

func setPingResponse(ResponseState *apatelet.NodeState_NodeResponseState, action scenario.Response, percentage int32) {
	ResponseState.PingResponse = action
	ResponseState.PingResponsePercentage = percentage
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

	ResponseState := nodeResponseState.GetResponseState()
	setCreatePodResponse(ResponseState, action, percentage)
	setUpdatePodResponse(ResponseState, action, percentage)
	setDeletePodResponse(ResponseState, action, percentage)
	setGetPodResponse(ResponseState, action, percentage)
	setGetPodStatusResponse(ResponseState, action, percentage)
}
