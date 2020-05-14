// Package crd provides functions and types to control Apate through CRDs.
package crd

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// SetPodFlags sets all flags for a pod.
func SetPodFlags(st *store.Store, pt *store.PodTask) error {
	if pt.State.CreatePodResponse != v1.ResponseUnset {
		(*st).SetPodFlag(pt.Label, events.PodCreatePodResponse, translateResponse(pt.State.CreatePodResponse))
	}

	if pt.State.UpdatePodResponse != v1.ResponseUnset {
		(*st).SetPodFlag(pt.Label, events.PodUpdatePodResponse, translateResponse(pt.State.UpdatePodResponse))
	}

	if pt.State.DeletePodResponse != v1.ResponseUnset {
		(*st).SetPodFlag(pt.Label, events.PodDeletePodResponse, translateResponse(pt.State.DeletePodResponse))
	}

	if pt.State.GetPodResponse != v1.ResponseUnset {
		(*st).SetPodFlag(pt.Label, events.PodGetPodResponse, translateResponse(pt.State.GetPodResponse))
	}

	if pt.State.GetPodStatusResponse != v1.ResponseUnset {
		(*st).SetPodFlag(pt.Label, events.PodGetPodStatusResponse, translateResponse(pt.State.GetPodStatusResponse))
	}

	if pt.State.PodResources != nil {
		(*st).SetPodFlag(pt.Label, events.PodStatus, pt.State.PodResources)
	}

	if pt.State.PodStatus != v1.PodStatusUnset {
		(*st).SetPodFlag(pt.Label, events.PodStatus, translatePodStatus(pt.State.PodStatus))
	}

	return nil
}

func translateResponse(input v1.EmulatedPodResponse) scenario.Response {
	switch input {
	case v1.ResponseNormal:
		return scenario.Response_RESPONSE_NORMAL
	case v1.ResponseError:
		return scenario.Response_RESPONSE_ERROR
	case v1.ResponseTimeout:
		return scenario.Response_RESPONSE_TIMEOUT
	case v1.ResponseUnset:
		fallthrough
	default:
		return scenario.Response_RESPONSE_UNSET
	}
}

func translatePodStatus(input v1.EmulatedPodStatus) scenario.PodStatus {
	switch input {
	case v1.PodStatusPending:
		return scenario.PodStatus_POD_STATUS_PENDING
	case v1.PodStatusRunning:
		return scenario.PodStatus_POD_STATUS_RUNNING
	case v1.PodStatusSucceeded:
		return scenario.PodStatus_POD_STATUS_SUCCEEDED
	case v1.PodStatusFailed:
		return scenario.PodStatus_POD_STATUS_FAILED
	case v1.PodStatusUnknown:
		return scenario.PodStatus_POD_STATUS_UNKNOWN
	case v1.PodStatusUnset:
		fallthrough
	default:
		return scenario.PodStatus_POD_STATUS_UNSET
	}
}
