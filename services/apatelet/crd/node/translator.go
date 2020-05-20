// Package node provides functions and types to deal with the NodeConfiguration on the apatelet
package node

import (
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// SetNodeFlags sets the correct flags for the apatelet
func SetNodeFlags(st *store.Store, state *v1.NodeConfigurationState) {
	// Set custom flags
	setCustomFlags(st, state.CustomState)

	// Check if the node should no longer respond to heartbeat
	if state.HeartbeatFailed {
		(*st).SetNodeFlag(events.NodePingResponse, scenario.ResponseTimeout)
	}

	// Check if there should be extra latency
	if state.NetworkLatency > 0 {
		(*st).SetNodeFlag(events.NodeAddedLatencyEnabled, true)
		(*st).SetNodeFlag(events.NodeAddedLatencyMsec, state.NetworkLatency)
	} else {
		(*st).SetNodeFlag(events.NodeAddedLatencyEnabled, false)
	}

	// Check if the node should fail
	if state.NodeFailed {
		(*st).SetNodeFlag(events.NodeCreatePodResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeUpdatePodResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeDeletePodResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeGetPodResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeGetPodStatusResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeGetPodsResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodePingResponse, scenario.ResponseTimeout)
		(*st).SetNodeFlag(events.NodeAddedLatencyEnabled, false)
	}
}

func setCustomFlags(st *store.Store, state *v1.NodeConfigurationDirectState) {
	if !isResponseUnset(state.CreatePodResponse) {
		(*st).SetNodeFlag(events.NodeCreatePodResponse, translateResponse(state.CreatePodResponse))
	}

	if !isResponseUnset(state.UpdatePodResponse) {
		(*st).SetNodeFlag(events.NodeUpdatePodResponse, translateResponse(state.UpdatePodResponse))
	}

	if !isResponseUnset(state.DeletePodResponse) {
		(*st).SetNodeFlag(events.NodeDeletePodResponse, translateResponse(state.DeletePodResponse))
	}

	if !isResponseUnset(state.GetPodResponse) {
		(*st).SetNodeFlag(events.NodeGetPodResponse, translateResponse(state.GetPodResponse))
	}

	if !isResponseUnset(state.GetPodStatusResponse) {
		(*st).SetNodeFlag(events.NodeGetPodStatusResponse, translateResponse(state.GetPodStatusResponse))
	}

	if !isResponseUnset(state.GetPodsResponse) {
		(*st).SetNodeFlag(events.NodeGetPodsResponse, translateResponse(state.GetPodsResponse))
	}

	if !isResponseUnset(state.NodePingResponse) {
		(*st).SetNodeFlag(events.NodePingResponse, translateResponse(state.NodePingResponse))
	}
}

func isResponseUnset(response v1.NodeResponse) bool {
	return response != v1.ResponseError && response != v1.ResponseNormal && response != v1.ResponseTimeout
}

func translateResponse(input v1.NodeResponse) scenario.Response {
	switch input {
	case v1.ResponseNormal:
		return scenario.ResponseNormal
	case v1.ResponseError:
		return scenario.ResponseError
	case v1.ResponseTimeout:
		return scenario.ResponseTimeout
	case v1.ResponseUnset:
		fallthrough
	default:
		return scenario.ResponseUnset
	}
}
