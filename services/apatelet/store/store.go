// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
)

// Store represents the state of the apatelet
type Store interface {
	// SetStartTime sets the value of the start time in the store. All other timing is based in this
	SetStartTime(int64)

	// SetNodeTasks adds or updates node tasks
	// Existing node tasks will be removed if not in the list of tasks
	SetNodeTasks([]*Task) error

	// SetPodTasks adds or updates pod CRD tasks to the queue based on their label (<namespace>/<name>)
	// Existing pod tasks will be removed if not in the list of tasks
	SetPodTasks(string, []*Task) error

	// RemovePodTasks removes pod CRD tasks from the queue based on their label (<namespace>/<name>)
	RemovePodTasks(string) error

	// LenTasks returns the amount of tasks left to be picked up
	LenTasks() int

	// PeekTask returns the start time of the next task in the priority queue, without removing it from the queue
	PeekTask() (int64, error)

	// PopTask returns the first task to be executed and removes it from the queue
	PopTask() (*Task, error)

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
const FlagNotFoundError = throw.ConstException("flag not found")

type flags map[events.EventFlag]interface{}
type podFlags map[string]flags

type store struct {
	startTime int64
	queue     *taskQueue
	queueLock sync.RWMutex

	nodeFlags    flags
	nodeFlagLock sync.RWMutex

	podFlags    podFlags
	podFlagLock sync.RWMutex
}

// NewStore returns an empty store
func NewStore() Store {
	q := newTaskQueue()
	heap.Init(q)

	return &store{
		queue:     q,
		nodeFlags: make(flags),
		podFlags:  make(podFlags),
	}
}

func (s *store) SetStartTime(time int64) {
	s.startTime = time
}

func (s *store) setTasksOfType(tasks []*Task, check TaskTypeCheck) error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	for i, task := range s.queue.tasks {
		typeCheck, err := check(task)
		if err != nil {
			return err
		}

		if typeCheck {
			if len(tasks) == 0 {
				heap.Remove(s.queue, i)
			} else {
				if tasks[0] != nil {
					s.queue.tasks[i] = tasks[0]
					// Replacing and then fixing instead of deleting all and pushing because it's slightly faster, see comments on heap.Fix
					heap.Fix(s.queue, i)
				}
				tasks = tasks[1:]
			}
		}
	}

	for _, remainingTask := range tasks {
		if remainingTask != nil {
			heap.Push(s.queue, remainingTask)
		}
	}

	return nil
}

func (s *store) SetNodeTasks(tasks []*Task) error {
	return s.setTasksOfType(tasks, func(task *Task) (bool, error) {
		isNode, err := task.IsNode()
		if err != nil {
			return false, err
		}

		return isNode, nil
	})
}

func (s *store) SetPodTasks(label string, tasks []*Task) error {
	return s.setTasksOfType(tasks, func(task *Task) (bool, error) {
		isPod, err := task.IsPod()
		if err != nil {
			return false, err
		}

		return isPod && task.PodTask.Label == label, nil
	})
}

func (s *store) RemovePodTasks(label string) error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	for i := len(s.queue.tasks) - 1; i >= 0; i-- {
		task := s.queue.tasks[i]

		isPod, err := task.IsPod()
		if err != nil {
			return err
		}

		if isPod && task.PodTask.Label == label {
			heap.Remove(s.queue, i)
		}
	}

	return nil
}

func (s *store) LenTasks() int {
	s.queueLock.RLock()
	defer s.queueLock.RUnlock()

	return s.queue.Len()
}

func (s *store) PeekTask() (int64, error) {
	s.queueLock.RLock()
	defer s.queueLock.RUnlock()

	if s.queue.Len() == 0 {
		return -1, errors.New("no tasks left")
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := s.queue.First().(*Task); ok {
		return task.RelativeTimestamp + s.startTime, nil
	}

	return -1, errors.New("array in pq magically changed to a different type")
}

func (s *store) PopTask() (*Task, error) {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	if s.queue.Len() == 0 {
		return nil, errors.New("no tasks left")
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := heap.Pop(s.queue).(*Task); ok {
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

func (s *store) GetPodFlag(label string, flag events.PodEventFlag) (interface{}, error) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	if val, ok := s.podFlags[label][flag]; ok {
		return val, nil
	}

	if dv, ok := defaultPodValues[flag]; ok {
		return dv, nil
	}

	return nil, FlagNotFoundError
}

func (s *store) SetPodFlag(label string, flag events.PodEventFlag, val interface{}) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	if conf, ok := s.podFlags[label]; ok {
		conf[flag] = val
	} else {
		s.podFlags[label] = make(flags)
		s.podFlags[label][flag] = val
	}
}

var defaultNodeValues = map[events.EventFlag]interface{}{
	events.NodeCreatePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeUpdatePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeDeletePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodResponse:       scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodStatusResponse: scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodsResponse:      scenario.Response_RESPONSE_NORMAL,
	events.NodePingResponse:         scenario.Response_RESPONSE_NORMAL,

	events.NodeAddedLatencyEnabled: false,
	events.NodeAddedLatencyMsec:    uint64(0),
}

var defaultPodValues = map[events.PodEventFlag]interface{}{
	events.PodCreatePodResponse:    scenario.Response_RESPONSE_UNSET,
	events.PodUpdatePodResponse:    scenario.Response_RESPONSE_UNSET,
	events.PodDeletePodResponse:    scenario.Response_RESPONSE_UNSET,
	events.PodGetPodResponse:       scenario.Response_RESPONSE_UNSET,
	events.PodGetPodStatusResponse: scenario.Response_RESPONSE_UNSET,

	events.PodResources: nil,

	events.PodStatus: scenario.PodStatus_POD_STATUS_UNSET,
}
