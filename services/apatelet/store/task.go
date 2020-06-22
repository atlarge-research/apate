package store

import (
	"time"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"

	"github.com/pkg/errors"
)

// TaskTypeCheck is function which is able to determine if the given task is of a certain type
type TaskTypeCheck func(*Task) (bool, error)

// Task is a task in the PQ
type Task struct {
	// The timestamp on which this task should be executed, relative to the start of the scenario
	RelativeTimestamp time.Duration

	PodTask  *PodTask
	NodeTask *NodeTask
}

// NodeTask is a task that should be executed on a node level
type NodeTask struct {
	State *nodeconfigv1.NodeConfigurationState
}

// IsNode returns true if the task is a node task
func (t *Task) IsNode() (bool, error) {
	isPod, err := t.IsPod()
	return !isPod, errors.Wrap(err, "isPod during IsNode failed")
}

// PodTask is a task that should be executed on a pod level
type PodTask struct {
	// The label of the CRD, should be <namespace>/<name>
	Label string
	State *podconfigv1.PodConfigurationState
}

// IsPod returns whether we are dealing with a pod (then PodTask should be non-nil) or a node (then NodeTask should be non-nil)
func (t *Task) IsPod() (bool, error) {
	if t.PodTask == nil && t.NodeTask == nil {
		return false, errors.New("pod task & node task are both nil")
	}
	if t.PodTask != nil && t.NodeTask != nil {
		return false, errors.New("pod task & node task are both non-nil")
	}
	return t.PodTask != nil, nil
}

// NewNodeTask creates a new task for a node event
func NewNodeTask(relativeTime time.Duration, state *nodeconfigv1.NodeConfigurationState) *Task {
	return &Task{
		RelativeTimestamp: relativeTime,
		PodTask:           nil,
		NodeTask: &NodeTask{
			State: state,
		},
	}
}

// NewPodTask creates a new task for a pod event
func NewPodTask(relativeTime time.Duration, label string, state *podconfigv1.PodConfigurationState) *Task {
	return &Task{
		RelativeTimestamp: relativeTime,
		PodTask: &PodTask{
			Label: label,
			State: state,
		},
	}
}
