package events

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)

// Note that these methods should preferably be implemented on their respective structs, but since these structs reside in /api and we don't want code there, we couldn't.
// An alternative would be wrapping the structs, but this would yield bad code.

// setCreatePodResponse sets the response to the CreatePod request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setCreatePodResponse(responseState *apatelet.ResponseState, response scenario.Response, percentage int32) {
	responseState.CreatePodResponse = response
	responseState.CreatePodResponsePercentage = percentage
}

// setUpdatePodResponse sets the response to the UpdatePod request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setUpdatePodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.UpdatePodResponse = action
	responseState.UpdatePodResponsePercentage = percentage
}

// setDeletePodResponse sets the response to the DeletePod request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setDeletePodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.DeletePodResponse = action
	responseState.DeletePodResponsePercentage = percentage
}

// setGetPodResponse sets the response to the GetPod request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setGetPodResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodResponse = action
	responseState.GetPodResponsePercentage = percentage
}

// setGetPodStatusResponse sets the response to the GetPodStatus request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setGetPodStatusResponse(responseState *apatelet.ResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodStatusResponse = action
	responseState.GetPodStatusResponsePercentage = percentage
}

// setGetPodsResponse sets the response to the GetPods request
// see https://godoc.org/github.com/virtual-kubelet/virtual-kubelet/node#PodLifecycleHandler
func setGetPodsResponse(responseState *apatelet.NodeState_NodeResponseState, action scenario.Response, percentage int32) {
	responseState.GetPodsResponse = action
	responseState.GetPodsResponsePercentage = percentage
}

// setPingResponse sets the response to the ping request
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
