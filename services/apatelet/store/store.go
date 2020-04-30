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

	// PollTask returns the start time of the next task in the priority queue, without removing it from the queue
	PollTask() (int64, error)

	// GetTask returns the first task to be executed and removes it from the queue
	GetTask() (*apatelet.Task, error)

	// GetFlag returns the value of the given flag, default value is zero for non-existing flags
	GetFlag(string) (int, error)

	// UpdateFlag updates the value of the given flag with the given value (the value is added to the current value),
	// or creates it with the given value if it didn't exist before
	UpdateFlag(string, int)
}

type store struct {
	queue    *taskQueue
	flags    map[string]int
	flagLock sync.RWMutex
}

// NewStore returns an empty store
func NewStore() Store {
	return &store{
		queue: newTaskQueue(),
		flags: make(map[string]int),
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

func (s *store) PollTask() (int64, error) {
	if s.queue.Len() == 0 {
		return -1, errors.New("no tasks left")
	}

	return int64(s.queue.Poll().(*apatelet.Task).Timestamp), nil
}

func (s *store) GetTask() (*apatelet.Task, error) {
	if s.queue.Len() == 0 {
		return nil, errors.New("no tasks left")
	}

	return heap.Pop(s.queue).(*apatelet.Task), nil
}

func (s *store) GetFlag(flag string) (int, error) {
	s.flagLock.RLock()
	defer s.flagLock.RUnlock()

	if val, ok := s.flags[flag]; ok {
		return val, nil
	}

	return 0, errors.New("flag defaulted to zero")
}

func (s *store) UpdateFlag(flag string, val int) {
	s.flagLock.Lock()
	defer s.flagLock.Unlock()

	s.flags[flag] += val
}
