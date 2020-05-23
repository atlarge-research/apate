// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"sync"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// Store represents the state of the apatelet
type Store interface {
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

type flags map[events.EventFlag]interface{}
type podFlags map[string]flags

type store struct {
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

func (s *store) setTasksOfType(newTasks []*Task, check TaskTypeCheck) error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	for i, task := range s.queue.tasks {
		typeCheck, err := check(task)
		if err != nil {
			return errors.Wrap(err, "failed to determine task type")
		}

		if typeCheck {
			if len(newTasks) == 0 {
				heap.Remove(s.queue, i)
			} else {
				if newTasks[0] != nil {
					s.queue.tasks[i] = newTasks[0]
					// Replacing and then fixing instead of deleting all and pushing because it's slightly faster, see comments on heap.Fix
					heap.Fix(s.queue, i)
				}
				newTasks = newTasks[1:]
			}
		}
	}

	for _, remainingTask := range newTasks {
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
			return errors.Wrap(err, "failed to determine task type")
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
		return -1, nil
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := s.queue.First().(*Task); ok {
		return task.RelativeTimestamp, nil
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

	return nil, errors.New("flag not found in get node flag")
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

	return nil, errors.New("flag not found in get pod flag")
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
	events.NodeCreatePodResponse:    scenario.ResponseNormal,
	events.NodeUpdatePodResponse:    scenario.ResponseNormal,
	events.NodeDeletePodResponse:    scenario.ResponseNormal,
	events.NodeGetPodResponse:       scenario.ResponseNormal,
	events.NodeGetPodStatusResponse: scenario.ResponseNormal,
	events.NodeGetPodsResponse:      scenario.ResponseNormal,
	events.NodePingResponse:         scenario.ResponseNormal,

	events.NodeAddedLatencyMsec:    uint64(0),
}

var defaultPodValues = map[events.PodEventFlag]interface{}{
	// Default is unset because then if no option is set we fall back to a node wide response

	events.PodCreatePodResponse:    scenario.ResponseUnset,
	events.PodUpdatePodResponse:    scenario.ResponseUnset,
	events.PodDeletePodResponse:    scenario.ResponseUnset,
	events.PodGetPodResponse:       scenario.ResponseUnset,
	events.PodGetPodStatusResponse: scenario.ResponseUnset,

	events.PodResources: nil,

	events.PodStatus: scenario.ResponseUnset,
}
