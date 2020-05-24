// Package pod provides functions and types to deal with the PodConfiguration CRD
package pod

import (
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// SetPodFlags sets all flags for a pod.
func SetPodFlags(st *store.Store, label string, pt *v1.PodConfigurationState) error {
	if !isResponseUnset(pt.CreatePodResponse) {
		(*st).SetPodFlag(label, events.PodCreatePodResponse, translateResponse(pt.CreatePodResponse))
	}

	if !isResponseUnset(pt.UpdatePodResponse) {
		(*st).SetPodFlag(label, events.PodUpdatePodResponse, translateResponse(pt.UpdatePodResponse))
	}

	if !isResponseUnset(pt.DeletePodResponse) {
		(*st).SetPodFlag(label, events.PodDeletePodResponse, translateResponse(pt.DeletePodResponse))
	}

	if !isResponseUnset(pt.GetPodResponse) {
		(*st).SetPodFlag(label, events.PodGetPodResponse, translateResponse(pt.GetPodResponse))
	}

	if !isResponseUnset(pt.GetPodStatusResponse) {
		(*st).SetPodFlag(label, events.PodGetPodStatusResponse, translateResponse(pt.GetPodStatusResponse))
	}

	if pt.PodResources != nil {
		resources, err := translatePodResources(pt.PodResources)
		if err != nil {
			return errors.Wrap(err, "failed to translate pod resources")
		}
		(*st).SetPodFlag(label, events.PodResources, resources)
	}

	if !isPodStatusUnset(pt.PodStatus) {
		(*st).SetPodFlag(label, events.PodStatus, translatePodStatus(pt.PodStatus))
	}

	return nil
}

func isResponseUnset(response v1.PodResponse) bool {
	return response != v1.ResponseError && response != v1.ResponseNormal && response != v1.ResponseTimeout
}

func isPodStatusUnset(status v1.PodStatus) bool {
	return status != v1.PodStatusFailed && status != v1.PodStatusPending && status != v1.PodStatusRunning && status != v1.PodStatusSucceeded && status != v1.PodStatusUnknown
}

func translateResponse(input v1.PodResponse) scenario.Response {
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

func translatePodStatus(input v1.PodStatus) scenario.PodStatus {
	switch input {
	case v1.PodStatusPending:
		return scenario.PodStatusPending
	case v1.PodStatusRunning:
		return scenario.PodStatusRunning
	case v1.PodStatusSucceeded:
		return scenario.PodStatusSucceeded
	case v1.PodStatusFailed:
		return scenario.PodStatusFailed
	case v1.PodStatusUnknown:
		return scenario.PodStatusUnknown
	case v1.PodStatusUnset:
		fallthrough
	default:
		return scenario.PodStatusUnset
	}
}

func translatePodResources(input *v1.PodResources) (*stats.PodStats, error) {
	memory, err := scenario.GetInBytes(input.Memory, "memory")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate memory bytes (was %v)", input.EphemeralStorage)
	}
	memoryUint := uint64(memory)

	storage, err := scenario.GetInBytes(input.Storage, "storage")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate storage bytes (was %v)", input.EphemeralStorage)
	}
	storageUint := uint64(storage)

	ephemeralStorage, err := scenario.GetInBytes(input.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate ephemeral storage bytes (was %v)", input.EphemeralStorage)
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
