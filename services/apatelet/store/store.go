// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
)

// Store represents the state of the apatelet
type Store interface {
	// EnqueueTasks creates a priority queue based on these tasks
	EnqueueTasks([]*apatelet.Task)

	// LenTasks returns the amount of tasks left to be picked up
	LenTasks() int

	// PeekTask returns the start time of the next task in the priority queue, without removing it from the queue
	PeekTask() (int64, error)

	// PopTask returns the first task to be executed and removes it from the queue
	PopTask() (*apatelet.Task, error)

	// GetNodeFlag returns the value of the given node flag
	GetNodeFlag(events.NodeEventFlag) (interface{}, error)

	// SetNodeFlag sets the value of the given node flag
	SetNodeFlag(events.NodeEventFlag, interface{})
}

// FlagNotFoundError is raised whenever a flag is not set
const FlagNotFoundError = throw.ConstException("flag not found")

type flags map[events.EventFlag]interface{}

type store struct {
	queue        *taskQueue
	nodeFlags    flags
	nodeFlagLock sync.RWMutex
}

// NewStore returns an empty store
func NewStore() Store {
	return &store{
		queue:     newTaskQueue(),
		nodeFlags: make(flags),
	}
}

func (s *store) EnqueueTasks(tasks []*apatelet.Task) {
	for _, task := range tasks {
		s.queue.Push(task)
	}

	heap.Init(s.queue)
}

func (s *store) LenTasks() int {
	return s.queue.Len()
}

func (s *store) PeekTask() (int64, error) {
	if s.queue.Len() == 0 {
		return -1, errors.New("no tasks left")
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := s.queue.First().(*apatelet.Task); ok {
		return task.AbsoluteTimestamp, nil
	}

	return -1, errors.New("array in pq magically changed to a different type")
}

func (s *store) PopTask() (*apatelet.Task, error) {
	if s.queue.Len() == 0 {
		return nil, errors.New("no tasks left")
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := heap.Pop(s.queue).(*apatelet.Task); ok {
		return task, nil
	}

	return nil, errors.New("array in pq magically changed to a different type")
}

func (s *store) GetNodeFlag(id events.NodeEventFlag) (interface{}, error) {
	s.nodeFlagLock.RLock()
	defer s.nodeFlagLock.RUnlock()

	if val, ok := s.nodeFlags[id]; ok {
		return val, nil
	}

	if dv, ok := defaultNodeValues[id]; ok {
		return dv, nil
	}

	return nil, FlagNotFoundError
}

func (s *store) SetNodeFlag(id events.NodeEventFlag, val interface{}) {
	s.nodeFlagLock.Lock()
	defer s.nodeFlagLock.Unlock()

	s.nodeFlags[id] = val
}

var defaultNodeValues = map[events.EventFlag]interface{}{
	events.NodeCreatePodResponse:           scenario.Response_NORMAL,
	events.NodeCreatePodResponsePercentage: int32(0),

	events.NodeUpdatePodResponse:           scenario.Response_NORMAL,
	events.NodeUpdatePodResponsePercentage: int32(0),

	events.NodeDeletePodResponse:           scenario.Response_NORMAL,
	events.NodeDeletePodResponsePercentage: int32(0),

	events.NodeGetPodResponse:           scenario.Response_NORMAL,
	events.NodeGetPodResponsePercentage: int32(0),

	events.NodeGetPodStatusResponse:           scenario.Response_NORMAL,
	events.NodeGetPodStatusResponsePercentage: int32(0),

	events.NodeGetPodsResponse:           scenario.Response_NORMAL,
	events.NodeGetPodsResponsePercentage: int32(0),

	events.NodePingResponse:           scenario.Response_NORMAL,
	events.NodePingResponsePercentage: int32(0),

	events.NodeEnableResourceAlteration: false,
	events.NodeMemoryUsage:              int64(0),
	events.NodeCPUUsage:                 int64(0),
	events.NodeStorageUsage:             int64(0),
	events.NodeEphemeralStorageUsage:    int64(0),

	events.NodeAddedLatencyEnabled: false,
	events.NodeAddedLatencyMsec:    int64(0),
}
