// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"sync"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// Store represents the state of the apatelet
type Store interface {
	// SetStartTime sets the value of the starttime in the store. All other tiing is based in this
	SetStartTime(int64)

	// TODO remove this when moving node to CRD
	// EnqueueTasks creates a priority queue based on these tasks
	EnqueueTasks([]*Task)

	// TODO add node equivalent when node to CRD
	// EnqueuePodTasks adds or updates pod CRD tasks to the queue based on their label (<namespace>/<name>)
	EnqueuePodTasks(string, []*Task) error

	// TODO add node equivalent when node to CRD
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
	startTime    int64
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

func (s *store) SetStartTime(time int64) {
	s.startTime = time
}

// TODO remove this when moving node to CRD
func (s *store) EnqueueTasks(tasks []*Task) {
	for _, task := range tasks {
		s.queue.Push(task)
	}

	heap.Init(s.queue)
}

func (s *store) EnqueuePodTasks(label string, newTasks []*Task) error {
	for i, task := range s.queue.tasks {
		isPod, err := task.IsPod()
		if err != nil {
			return errors.Wrap(err, "failed to read pod type")
		}

		if isPod && task.PodTask.Label == label {
			if len(newTasks) == 0 {
				heap.Remove(s.queue, i)
			} else {
				s.queue.tasks[i] = newTasks[0]
				// Replacing and then fixing instead of deleting all and pushing because it's slightly faster, see comments on heap.Fix
				heap.Fix(s.queue, i)
				newTasks = newTasks[1:]
			}
		}
	}

	for _, remainingTask := range newTasks {
		heap.Push(s.queue, remainingTask)
	}

	return nil
}

func (s *store) RemovePodTasks(label string) error {
	for i := len(s.queue.tasks) - 1; i >= 0; i-- {
		task := s.queue.tasks[i]

		isPod, err := task.IsPod()
		if err != nil {
			return errors.Wrap(err, "failed to read pod type")
		}

		if isPod && task.PodTask.Label == label {
			heap.Remove(s.queue, i)
		}
	}

	return nil
}

func (s *store) LenTasks() int {
	return s.queue.Len()
}

func (s *store) PeekTask() (int64, error) {
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
	events.NodeCreatePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeUpdatePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeDeletePodResponse:    scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodResponse:       scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodStatusResponse: scenario.Response_RESPONSE_NORMAL,
	events.NodeGetPodsResponse:      scenario.Response_RESPONSE_NORMAL,
	events.NodePingResponse:         scenario.Response_RESPONSE_NORMAL,

	events.NodeAddedLatencyEnabled: false,
	events.NodeAddedLatencyMsec:    int64(0),
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
