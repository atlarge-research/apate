package store

import (
	"errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
)

// Task is a task in the PQ
type Task struct {
	// The timestamp on which this task should be executed, relative to the start of the scenario
	RelativeTimestamp int64

	PodTask  *PodTask
	NodeTask NodeTask
}

// NodeTask is a task that should be executed on a node level
type NodeTask *apatelet.Task // TODO change when moving node to CRD

// PodTask is a task that should be executed on a pod level
type PodTask struct {
	// The label of the CRD, should be <namespace>/<name>
	Label string
	State *v1.EmulatedPodState
}

// IsPod returns whether we are dealing with a pod (then PodTask should be non-nil) or a node (then NodeTask should be non-nil)
func (t *Task) IsPod() (bool, error) {
	if t.PodTask == nil && t.NodeTask == nil {
		return false, errors.New("pod task & node task are nil")
	}
	return t.PodTask != nil, nil
}

// NewNodeTask creates a new task for a node event
func NewNodeTask(relativeTime int64, task NodeTask) *Task {
	return &Task{
		RelativeTimestamp: relativeTime,
		PodTask:           nil,
		NodeTask:          task,
	}
}

// NewPodTask creates a new task for a pod event
func NewPodTask(relativeTime int64, label string, state *v1.EmulatedPodState) *Task {
	return &Task{
		RelativeTimestamp: relativeTime,
		PodTask: &PodTask{
			Label: label,
			State: state,
		},
	}
}
