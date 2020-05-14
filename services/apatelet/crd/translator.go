package crd

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

func SetPodFlags(st *store.Store, pt *store.PodTask) error {
	if pt.State.CreatePodResponse != v1.RESPONSE_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodCreatePodResponse, translateResponse(pt.State.CreatePodResponse))
	}

	if pt.State.UpdatePodResponse != v1.RESPONSE_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodUpdatePodResponse, translateResponse(pt.State.UpdatePodResponse))
	}

	if pt.State.DeletePodResponse != v1.RESPONSE_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodDeletePodResponse, translateResponse(pt.State.DeletePodResponse))
	}

	if pt.State.GetPodResponse != v1.RESPONSE_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodGetPodResponse, translateResponse(pt.State.GetPodResponse))
	}

	if pt.State.GetPodStatusResponse != v1.RESPONSE_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodGetPodStatusResponse, translateResponse(pt.State.GetPodStatusResponse))
	}

	// Resource usage
	if err := setResourceBytes(st, pt.Label, pt.State.ResourceUsage.Memory, events.PodMemoryUsage); err != nil {
		return err
	}

	if pt.State.ResourceUsage.CPU != -1 {
		(*st).SetPodFlag(pt.Label, events.PodCPUUsage, pt.State.ResourceUsage.CPU)
	}

	if err := setResourceBytes(st, pt.Label, pt.State.ResourceUsage.Storage, events.PodStorageUsage); err != nil {
		return err
	}

	if err := setResourceBytes(st, pt.Label, pt.State.ResourceUsage.EphemeralStorage, events.PodEphemeralStorageUsage); err != nil {
		return err
	}

	if pt.State.PodStatus != v1.POD_STATUS_UNSET {
		(*st).SetPodFlag(pt.Label, events.PodStatus, translatePodStatus(pt.State.PodStatus))
	}

	return nil
}

func setResourceBytes(st *store.Store, label string, unit string, flag events.PodEventFlag) error {
	if unit != "-1" {
		bytes, err := translate.GetInBytes(unit, "storage")
		if err != nil {
			return err
		}
		(*st).SetPodFlag(label, flag, bytes)
	}

	return nil
}

func translateResponse(input v1.EmulatedPodResponse) scenario.Response {
	switch input {
	case v1.RESPONSE_NORMAL:
		return scenario.Response_RESPONSE_NORMAL
	case v1.RESPONSE_ERROR:
		return scenario.Response_RESPONSE_ERROR
	case v1.RESPONSE_TIMEOUT:
		return scenario.Response_RESPONSE_TIMEOUT
	case v1.RESPONSE_UNSET:
		fallthrough
	default:
		return scenario.Response_RESPONSE_UNSET
	}
}

func translatePodStatus(input v1.EmulatedPodStatus) scenario.PodStatus {
	switch input {
	case v1.POD_STATUS_PENDING:
		return scenario.PodStatus_POD_STATUS_PENDING
	case v1.POD_STATUS_RUNNING:
		return scenario.PodStatus_POD_STATUS_RUNNING
	case v1.POD_STATUS_SUCCEEDED:
		return scenario.PodStatus_POD_STATUS_SUCCEEDED
	case v1.POD_STATUS_FAILED:
		return scenario.PodStatus_POD_STATUS_FAILED
	case v1.POD_STATUS_UNKNOWN:
		return scenario.PodStatus_POD_STATUS_UNKNOWN
	case v1.POD_STATUS_UNSET:
		fallthrough
	default:
		return scenario.PodStatus_POD_STATUS_UNSET
	}
}
