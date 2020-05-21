package translate

// TODO remove this when moving node to CRD

import (
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	"github.com/docker/go-units"
)

// DesugarTimestamp takes in a formatted duration string and returns the duration specified by this string in nanoseconds.
func DesugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Nanoseconds()), nil
}

// GetInBytes takes in a formatted size string and returns the size specified by this string in bytes.
// It errors if the size is below zero.
func GetInBytes(unit string, unitName string) (int64, error) {
	unitInt, err := units.RAMInBytes(unit)
	if err != nil {
		return 0, err
	} else if unitInt < 0 {
		return 0, errors.Errorf("%s usage should be at least 0", unitName)
	}
	return unitInt, nil
}

// EventFlags maps event flags to their any value
type EventFlags map[events.EventFlag]*anypb.Any

func newEventFlags() EventFlags {
	return make(map[events.EventFlag]*anypb.Any)
}

func (ef *EventFlags) flags(value interface{}, flags []events.EventFlag) {
	for _, flag := range flags {
		ef.flag(value, flag)
	}
}

func (ef EventFlags) flag(value interface{}, flag events.EventFlag) {
	mv := any.MarshalOrDie(value)
	ef[flag] = mv
}

var nodeEventFlags = []events.NodeEventFlag{
	events.NodeCreatePodResponse,
	events.NodeUpdatePodResponse,
	events.NodeDeletePodResponse,
	events.NodeGetPodResponse,
	events.NodeGetPodStatusResponse,
	events.NodeGetPodsResponse,
	events.NodePingResponse,
}
