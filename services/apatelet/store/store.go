// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"errors"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"sync"

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

	// GetPodFlag returns the value of the given pod flag for a configuration
	GetPodFlag(string, events.PodEventFlag) (interface{}, error)

	// SetNodeFlag sets the value of the given pod flag for a configuration
	SetPodFlag(string, events.PodEventFlag, interface{})
}

// FlagNotFoundError is raised whenever a flag is not set
const FlagNotFoundError = throw.Exception("flag not found")

type flags map[events.EventFlag]interface{}
type podFlags map[string]flags

type store struct {
	queue        *taskQueue
	nodeFlags    flags
	nodeFlagLock sync.RWMutex
	podFlags     podFlags
	podFlagLock  sync.RWMutex
}

// NewStore returns an empty store
func NewStore() Store {
	return &store{
		queue:     newTaskQueue(),
		nodeFlags: make(flags),
		podFlags:  make(podFlags),
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

func (s *store) GetPodFlag(configuration string, flag events.PodEventFlag) (interface{}, error) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	if val, ok := s.podFlags[configuration][flag]; ok {
		return val, nil
	}

	if dv, ok := defaultPodValues[flag]; ok {
		return dv, nil
	}

	return nil, FlagNotFoundError
}

func (s *store) SetPodFlag(configuration string, flag events.PodEventFlag, val interface{}) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	if conf, ok := s.podFlags[configuration]; ok {
		conf[flag] = val
	} else {
		s.podFlags[configuration] = make(flags)
		s.podFlags[configuration][flag] = val
	}
}

var defaultNodeValues = map[events.EventFlag]interface{} {
	events.NodeCreatePodResponse: scenario.Response_NORMAL,
	events.NodeCreatePodResponsePercentage: 0,

	events.NodeUpdatePodResponse: scenario.Response_NORMAL,
	events.NodeUpdatePodResponsePercentage: 0,

	events.NodeDeletePodResponse: scenario.Response_NORMAL,
	events.NodeDeletePodResponsePercentage: 0,

	events.NodeGetPodResponse: scenario.Response_NORMAL,
	events.NodeGetPodResponsePercentage: 0,

	events.NodeGetPodStatusResponse: scenario.Response_NORMAL,
	events.NodeGetPodStatusResponsePercentage: 0,

	events.NodeGetPodsResponse: scenario.Response_NORMAL,
	events.NodeGetPodsResponsePercentage: 0,

	events.NodePingResponse: scenario.Response_NORMAL,
	events.NodePingResponsePercentage: 0,

	events.NodeEnableResourceAlteration: false,
	events.NodeMemoryUsage: 0,
	events.NodeCPUUsage: 0,
	events.NodeStorageUsage: 0,
	events.NodeEphemeralStorageUsage: 0,

	events.NodeAddedLatencyEnabled: false,
	events.NodeAddedLatencyMsec: 0,
}

var defaultPodValues = map[events.PodEventFlag]interface{} {
	events.PodCreatePodResponse: scenario.Response_NORMAL,
	events.PodCreatePodResponsePercentage: 0,

	events.PodUpdatePodResponse: scenario.Response_NORMAL,
	events.PodUpdatePodResponsePercentage: 0,

	events.PodDeletePodResponse: scenario.Response_NORMAL,
	events.PodDeletePodResponsePercentage: 0,

	events.PodGetPodResponse: scenario.Response_NORMAL,
	events.PodGetPodResponsePercentage: 0,

	events.PodGetPodStatusResponse: scenario.Response_NORMAL,
	events.PodGetPodStatusResponsePercentage: 0,

	events.PodUpdatePodStatus: scenario.PodStatus_POD_RUNNING,
	events.PodUpdatePodStatusPercentage: 0,
}
