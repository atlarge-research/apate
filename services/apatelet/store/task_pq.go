package store

import (
	"reflect"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
)

// Task is a task in the PQ
type Task struct {
	// The timestamp on which this task should be executed, relative to the start of the scenario
	RelativeTimestamp int64

	// Whether we are dealing with a pod (then PodTask should be non-nil) or a node (then NodeTask should be non-nil)
	IsPod    bool
	PodTask  *PodTask
	NodeTask NodeTask
}

// NodeTask is a task that should be executed on a node level
type NodeTask *apatelet.Task // TODO change when moving node to CRD

// PodTask is a task that should be executed on a pod elvel
type PodTask struct {
	// The label of the CRD, should be <namespace>/<name>
	Label string

	// The internal state to be set
	State *v1.EmulatedPodState
}

// taskQueue is a thread-safe priority queue based on a min-heap for tasks in the private scenario
type taskQueue struct {
	tasks []*Task
	lock  sync.RWMutex
}

func newTaskQueue() *taskQueue {
	return &taskQueue{
		tasks: make([]*Task, 0),
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

	return q.tasks[i].RelativeTimestamp < q.tasks[j].RelativeTimestamp
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

	if task, ok := x.(*Task); ok {
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
