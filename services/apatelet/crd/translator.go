// Package crd provides functions and types to control Apate through CRDs.
package crd

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// SetPodFlags sets all flags for a pod.
func SetPodFlags(st *store.Store, pt *store.PodTask) error {
	if !isResponseUnset(pt.State.CreatePodResponse) {
		(*st).SetPodFlag(pt.Label, events.PodCreatePodResponse, translateResponse(pt.State.CreatePodResponse))
	}

	if !isResponseUnset(pt.State.UpdatePodResponse) {
		(*st).SetPodFlag(pt.Label, events.PodUpdatePodResponse, translateResponse(pt.State.UpdatePodResponse))
	}

	if !isResponseUnset(pt.State.DeletePodResponse) {
		(*st).SetPodFlag(pt.Label, events.PodDeletePodResponse, translateResponse(pt.State.DeletePodResponse))
	}

	if !isResponseUnset(pt.State.GetPodResponse) {
		(*st).SetPodFlag(pt.Label, events.PodGetPodResponse, translateResponse(pt.State.GetPodResponse))
	}

	if !isResponseUnset(pt.State.GetPodStatusResponse) {
		(*st).SetPodFlag(pt.Label, events.PodGetPodStatusResponse, translateResponse(pt.State.GetPodStatusResponse))
	}

	if pt.State.PodResources != nil {
		resources, err := translatePodResources(pt.State.PodResources)
		if err != nil {
			return err
		}
		(*st).SetPodFlag(pt.Label, events.PodResources, resources)
	}

	if !isPodStatusUnset(pt.State.PodStatus) {
		(*st).SetPodFlag(pt.Label, events.PodStatus, translatePodStatus(pt.State.PodStatus))
	}

	return nil
}

func isResponseUnset(response v1.EmulatedPodResponse) bool {
	return response != v1.ResponseError && response != v1.ResponseNormal && response != v1.ResponseTimeout
}

func isPodStatusUnset(status v1.EmulatedPodStatus) bool {
	return status != v1.PodStatusFailed && status != v1.PodStatusPending && status != v1.PodStatusRunning && status != v1.PodStatusSucceeded && status != v1.PodStatusUnknown
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

func translatePodResources(input *v1.EmulatedPodResourceUsage) (*stats.PodStats, error) {
	memory, err := translate.GetInBytes(input.Memory, "memory")
	if err != nil {
		return nil, err
	}
	memoryUint := uint64(memory)

	storage, err := translate.GetInBytes(input.Storage, "storage")
	if err != nil {
		return nil, err
	}
	storageUint := uint64(storage)

	ephemeralStorage, err := translate.GetInBytes(input.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return nil, err
	}
	ephemeralStorageUint := uint64(ephemeralStorage)

	return &stats.PodStats{
		CPU: &stats.CPUStats{
			Time: metav1.Time{
				Time: time.Now(),
			},
			UsageNanoCores: &input.CPU,
		},
		Memory: &stats.MemoryStats{
			Time: metav1.Time{
				Time: time.Now(),
			},
			UsageBytes: &memoryUint,
		},
		VolumeStats: []stats.VolumeStats{
			{
				FsStats: stats.FsStats{
					Time: metav1.Time{
						Time: time.Now(),
					},
					UsedBytes: &storageUint,
				},
			},
		},
		EphemeralStorage: &stats.FsStats{
			Time: metav1.Time{
				Time: time.Now(),
			},
			UsedBytes: &ephemeralStorageUint,
		},
	}, nil
}
