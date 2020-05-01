// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
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

	// GetFlag returns the value of the given flag
	GetFlag(int) (interface{}, error)

	// SetFlag sets the value of the given flag
	SetFlag(int, interface{})
}

type store struct {
	queue    *taskQueue
	flags    map[int]interface{}
	flagLock sync.RWMutex
}

// NewStore returns an empty store
func NewStore() Store {
	return &store{
		queue: newTaskQueue(),
		flags: make(map[int]interface{}),
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
		return int64(task.Timestamp), nil
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

func (s *store) GetFlag(id int) (interface{}, error) {
	s.flagLock.RLock()
	defer s.flagLock.RUnlock()

	if val, ok := s.flags[id]; ok {
		return val, nil
	}

	return nil, errors.New("flag not set")
}

func (s *store) SetFlag(id int, val interface{}) {
	s.flagLock.Lock()
	defer s.flagLock.Unlock()

	s.flags[id] = val
}
