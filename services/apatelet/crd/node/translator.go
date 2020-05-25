// Package node provides functions and types to deal with the NodeConfiguration on the apatelet
package node

import (
	"time"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// SetNodeFlags sets the correct flags for the apatelet
func SetNodeFlags(st *store.Store, state *nodeconfigv1.NodeConfigurationState) {
	// Set custom flags
	setCustomFlags(st, state.CustomState)

	// Check if the node should no longer respond to heartbeat
	if state.HeartbeatFailed {
		(*st).SetNodeFlag(events.NodePingResponse, scenario.ResponseTimeout)
	}

	// Set latency
	latency, err := time.ParseDuration(state.NetworkLatency)
	if err == nil && latency >= 0 {
		// Ignore errors explicitly, only valid ints are seen as updates
		(*st).SetNodeFlag(events.NodeAddedLatency, latency)
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
	}
}

func setCustomFlags(st *store.Store, state *nodeconfigv1.NodeConfigurationCustomState) {
	// Check if there were no custom flags
	if state == nil {
		return
	}

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

func isResponseUnset(response nodeconfigv1.NodeResponse) bool {
	return response != nodeconfigv1.ResponseError && response != nodeconfigv1.ResponseNormal && response != nodeconfigv1.ResponseTimeout
}

func translateResponse(input nodeconfigv1.NodeResponse) scenario.Response {
	switch input {
	case nodeconfigv1.ResponseNormal:
		return scenario.ResponseNormal
	case nodeconfigv1.ResponseError:
		return scenario.ResponseError
	case nodeconfigv1.ResponseTimeout:
		return scenario.ResponseTimeout
	case nodeconfigv1.ResponseUnset:
		fallthrough
	default:
		return scenario.ResponseUnset
	}
}
