// Package node provides functions and types to deal with the NodeConfiguration on the apatelet
package node

import (
	"time"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/scenario"
	"github.com/atlarge-research/apate/pkg/scenario/events"
	"github.com/atlarge-research/apate/services/apatelet/store"
)

// SetNodeFlags sets the correct flags for the apatelet
func SetNodeFlags(st *store.Store, state *nodeconfigv1.NodeConfigurationState) {
	flags := make(store.Flags)

	// Set custom flags
	setCustomFlags(flags, state.CustomState)

	// Check if the node should no longer respond to heartbeat
	if state.HeartbeatFailed {
		flags[events.NodePingResponse] = scenario.ResponseTimeout
	}

	// Set latency
	latency, err := time.ParseDuration(state.NetworkLatency)
	if err == nil && latency >= 0 {
		// Ignore errors explicitly, only valid ints are seen as updates
		flags[events.NodeAddedLatency] = latency
	}

	// Check if the node should fail
	if state.NodeFailed {
		flags[events.NodeCreatePodResponse] = scenario.ResponseTimeout
		flags[events.NodeUpdatePodResponse] = scenario.ResponseTimeout
		flags[events.NodeDeletePodResponse] = scenario.ResponseTimeout
		flags[events.NodeGetPodResponse] = scenario.ResponseTimeout
		flags[events.NodeGetPodStatusResponse] = scenario.ResponseTimeout
		flags[events.NodeGetPodsResponse] = scenario.ResponseTimeout
		flags[events.NodePingResponse] = scenario.ResponseTimeout
	}

	(*st).SetNodeFlags(flags)
}

func setCustomFlags(flags store.Flags, state *nodeconfigv1.NodeConfigurationCustomState) {
	// Check if there were no custom flags
	if state == nil {
		return
	}

	if !isResponseUnset(state.CreatePodResponse) {
		flags[events.NodeCreatePodResponse] = translateResponse(state.CreatePodResponse)
	}

	if !isResponseUnset(state.UpdatePodResponse) {
		flags[events.NodeUpdatePodResponse] = translateResponse(state.UpdatePodResponse)
	}

	if !isResponseUnset(state.DeletePodResponse) {
		flags[events.NodeDeletePodResponse] = translateResponse(state.DeletePodResponse)
	}

	if !isResponseUnset(state.GetPodResponse) {
		flags[events.NodeGetPodResponse] = translateResponse(state.GetPodResponse)
	}

	if !isResponseUnset(state.GetPodStatusResponse) {
		flags[events.NodeGetPodStatusResponse] = translateResponse(state.GetPodStatusResponse)
	}

	if !isResponseUnset(state.GetPodsResponse) {
		flags[events.NodeGetPodsResponse] = translateResponse(state.GetPodsResponse)
	}

	if !isResponseUnset(state.NodePingResponse) {
		flags[events.NodePingResponse] = translateResponse(state.NodePingResponse)
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
