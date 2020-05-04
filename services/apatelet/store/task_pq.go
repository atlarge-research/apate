package store

import (
	"reflect"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

// taskQueue is a thread-safe priority queue based on a min-heap for tasks in the private scenario
type taskQueue struct {
	tasks []*apatelet.Task
	lock  sync.RWMutex
}

func newTaskQueue() *taskQueue {
	return &taskQueue{
		tasks: make([]*apatelet.Task, 0),
	}
}

// Len is the number of elements in the queue
func (q *taskQueue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.tasks)
}

// Less reports whether the element with index i should sort before the element with index j
func (q *taskQueue) Less(i, j int) bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.tasks[i].AbsoluteTimestamp < q.tasks[j].AbsoluteTimestamp
}

// Swap swaps the elements with indexes i and j
func (q *taskQueue) Swap(i, j int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.tasks[i], q.tasks[j] = q.tasks[j], q.tasks[i]
}

// Push pushes a new task to the queue
func (q *taskQueue) Push(x interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// No-op if x is nil or not a task
	if x == nil || (reflect.ValueOf(x).Kind() == reflect.Ptr && reflect.ValueOf(x).IsNil()) {
		return
	}

	if task, ok := x.(*apatelet.Task); ok {
		q.tasks = append(q.tasks, task)
	}
}

// Pop returns the first task in the queue and removes it
func (q *taskQueue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	old := q.tasks
	n := len(old)
	task := old[n-1]
	old[n-1] = nil
	q.tasks = old[0 : n-1]
	return task
}

// First returns the first task in the queue without removing it
func (q *taskQueue) First() interface{} {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.tasks[0]
}
