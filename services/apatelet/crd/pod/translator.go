// Package pod provides functions and types to deal with the PodConfiguration CRD
package pod

import (
	"github.com/finitum/node-cli/stats"

	"github.com/atlarge-research/apate/pkg/scenario"

	"github.com/pkg/errors"

	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/scenario/events"
	"github.com/atlarge-research/apate/services/apatelet/store"
)

// SetPodFlags sets all flags for a pod.
func SetPodFlags(st *store.Store, label string, pt *podconfigv1.PodConfigurationState) error {
	flags, err := TranslatePodFlags(pt)
	if err != nil {
		return errors.Wrap(err, "failed to translate pod state into flags")
	}

	(*st).SetPodFlags(label, flags)

	return nil
}

// TranslatePodFlags translates a pod configuration into a map of flags.
func TranslatePodFlags(pt *podconfigv1.PodConfigurationState) (store.Flags, error) {
	flags := make(store.Flags)
	if !isResponseUnset(pt.CreatePodResponse) {
		flags[events.PodCreatePodResponse] = translateResponse(pt.CreatePodResponse)
	}

	if !isResponseUnset(pt.UpdatePodResponse) {
		flags[events.PodUpdatePodResponse] = translateResponse(pt.UpdatePodResponse)
	}

	if !isResponseUnset(pt.DeletePodResponse) {
		flags[events.PodDeletePodResponse] = translateResponse(pt.DeletePodResponse)
	}

	if !isResponseUnset(pt.GetPodResponse) {
		flags[events.PodGetPodResponse] = translateResponse(pt.GetPodResponse)
	}

	if !isResponseUnset(pt.GetPodStatusResponse) {
		flags[events.PodGetPodStatusResponse] = translateResponse(pt.GetPodStatusResponse)
	}

	if pt.PodResources != nil {
		resources, err := translatePodResources(pt.PodResources)
		if err != nil {
			return nil, errors.Wrap(err, "failed to translate pod resources")
		}
		flags[events.PodResources] = resources
	}

	if !isPodStatusUnset(pt.PodStatus) {
		flags[events.PodStatus] = translatePodStatus(pt.PodStatus)
	}

	return flags, nil
}

func isResponseUnset(response podconfigv1.PodResponse) bool {
	return response != podconfigv1.ResponseError && response != podconfigv1.ResponseNormal && response != podconfigv1.ResponseTimeout
}

func isPodStatusUnset(status podconfigv1.PodStatus) bool {
	return status != podconfigv1.PodStatusFailed && status != podconfigv1.PodStatusPending && status != podconfigv1.PodStatusRunning && status != podconfigv1.PodStatusSucceeded && status != podconfigv1.PodStatusUnknown
}

func translateResponse(input podconfigv1.PodResponse) scenario.Response {
	switch input {
	case podconfigv1.ResponseNormal:
		return scenario.ResponseNormal
	case podconfigv1.ResponseError:
		return scenario.ResponseError
	case podconfigv1.ResponseTimeout:
		return scenario.ResponseTimeout
	case podconfigv1.ResponseUnset:
		fallthrough
	default:
		return scenario.ResponseUnset
	}
}

func translatePodStatus(input podconfigv1.PodStatus) scenario.PodStatus {
	switch input {
	case podconfigv1.PodStatusPending:
		return scenario.PodStatusPending
	case podconfigv1.PodStatusRunning:
		return scenario.PodStatusRunning
	case podconfigv1.PodStatusSucceeded:
		return scenario.PodStatusSucceeded
	case podconfigv1.PodStatusFailed:
		return scenario.PodStatusFailed
	case podconfigv1.PodStatusUnknown:
		return scenario.PodStatusUnknown
	case podconfigv1.PodStatusUnset:
		fallthrough
	default:
		return scenario.PodStatusUnset
	}
}

func translatePodResources(input *podconfigv1.PodResources) (*stats.PodStats, error) {
	memory, err := scenario.GetInBytes(input.Memory, "memory")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate memory bytes (was %v)", input.EphemeralStorage)
	}

	storage, err := scenario.GetInBytes(input.Storage, "storage")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate storage bytes (was %v)", input.EphemeralStorage)
	}

	ephemeralStorage, err := scenario.GetInBytes(input.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to translate ephemeral storage bytes (was %v)", input.EphemeralStorage)
	}

	return &stats.PodStats{
		PodRef:             stats.PodReference{},
		UsageNanoCores:     input.CPU,
		UsageBytesMemory:   uint64(memory),
		UsedBytesEphemeral: uint64(ephemeralStorage),
		UsedBytesStorage:   uint64(storage),
	}, nil
}
