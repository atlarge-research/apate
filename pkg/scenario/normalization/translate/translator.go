// Package translate contains utilities to translate events of tasks between public API and internal API formats
// TODO remove this when moving node to CRD
package translate

import (
	"errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	ef "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

const percentageRangeErrorMessage = "percentage should be between 0 and 100"

// EventTranslator is a utility to translate events between events sent through the public api to events understood by the Apatelets
type EventTranslator struct {
	originalTask *controlplane.Task
	newTask      *apatelet.Task
}

// NewEventTranslator constructs a new EventTranslator
func NewEventTranslator(originalTask *controlplane.Task, newTask *apatelet.Task) *EventTranslator {
	return &EventTranslator{
		originalTask: originalTask,
		newTask:      newTask,
	}
}

// TranslateEvent translates events sent through the public api to events understood by the Apatelets
func (et *EventTranslator) TranslateEvent() error {
	if et.originalTask.Event == nil {
		return nil
	}

	nef := newEventFlags()

	// et.originalTask.Event can be one of many types (see generated protobuf code)
	// ne will be the cast version of this event to the corresponding event, depending on the case
	switch ne := et.originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		nef.flags(scenario.Response_RESPONSE_TIMEOUT, nodeEventFlags)

	case *controlplane.Task_NetworkLatency:
		if ne.NetworkLatency.LatencyMsec < 0 {
			return errors.New("latency should be at least 0")
		}

		nef.flag(true, ef.NodeAddedLatencyEnabled)
		nef.flag(ne.NetworkLatency.LatencyMsec, ef.NodeAddedLatencyMsec)

	case *controlplane.Task_TimeoutKeepHeartbeat:
		nef.flags(scenario.Response_RESPONSE_TIMEOUT, nodeEventFlags)

		// Reset ping
		nef.flag(scenario.Response_RESPONSE_NORMAL, ef.NodePingResponse)

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		nef.flag(scenario.Response_RESPONSE_TIMEOUT, ef.NodePingResponse)

	case *controlplane.Task_NodeResponseState:
		state := ne.NodeResponseState

		if state.Percentage < 0 || state.Percentage > 100 {
			return errors.New(percentageRangeErrorMessage)
		}

		switch state.Type {
		case controlplane.RequestType_CREATE_POD:
			nef.flag(state.Response, ef.NodeCreatePodResponse)

		case controlplane.RequestType_UPDATE_POD:
			nef.flag(state.Response, ef.NodeUpdatePodResponse)

		case controlplane.RequestType_DELETE_POD:
			nef.flag(state.Response, ef.NodeDeletePodResponse)

		case controlplane.RequestType_GET_POD:
			nef.flag(state.Response, ef.NodeGetPodResponse)

		case controlplane.RequestType_GET_POD_STATUS:
			nef.flag(state.Response, ef.NodeGetPodStatusResponse)

		case controlplane.RequestType_GET_PODS:
			nef.flag(state.Response, ef.NodeGetPodsResponse)

		case controlplane.RequestType_PING:
			nef.flag(state.Response, ef.NodePingResponse)
		}

	case *controlplane.Task_CustomFlags:
		nef = ne.CustomFlags.CustomFlags
	}

	et.newTask.NodeEventFlags = nef

	return nil
}
