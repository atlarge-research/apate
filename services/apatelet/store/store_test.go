package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

// TODO: Change tests to use absolute timestamp instead of relative one
// TODO: Use int64 from absolute timestamp

// TestEmptyQueue ensures the store starts with an empty queue
func TestEmptyQueue(t *testing.T) {
	st := NewStore()

	// Make sure the amount of tasks starts at zero
	assert.Equal(t, 0, st.LenTasks())

	// Make sure both poll and get return an error
	_, pollErr := st.PeekTask()
	_, getErr := st.PopTask()
	assert.Error(t, pollErr)
	assert.Error(t, getErr)
}

// TestGetSingleTask ensures a retrieved task is also deleted
func TestGetSingleTask(t *testing.T) {
	task := &apatelet.Task{}
	st := NewStore()

	// Enqueue single task
	st.EnqueueTasks([]*apatelet.Task{task})

	// Retrieve single task and verify it was the original one
	retrieved, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task, retrieved)

	// Also verify it was removed
	assert.Equal(t, 0, st.LenTasks())
}

// TestGetSingleTask ensures a polled task is not deleted
func TestPollSingleTask(t *testing.T) {
	task := &apatelet.Task{AbsoluteTimestamp: 424242}
	st := NewStore()

	// Enqueue single task
	st.EnqueueTasks([]*apatelet.Task{task})

	// Poll single task and verify the timestamp is correct
	retrieved, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task.AbsoluteTimestamp, retrieved)

	// Also verify it was not removed
	assert.Equal(t, 1, st.LenTasks())
}

// TestMultipleTasks ensures the priority queue actually sorts the tasks properly (earliest task first)
func TestMultipleTasks(t *testing.T) {
	task1 := &apatelet.Task{AbsoluteTimestamp: 213123}
	task2 := &apatelet.Task{AbsoluteTimestamp: 4242}
	task3 := &apatelet.Task{AbsoluteTimestamp: 83481234}
	st := NewStore()

	// Enqueue tasks
	st.EnqueueTasks([]*apatelet.Task{task1, task2, task3})

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.AbsoluteTimestamp, firstTaskTime)

	// Retrieve first two tasks
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task2, firstTask)

	secondTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task1, secondTask)

	// Verify there is still one task left
	lastTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, 1, st.LenTasks())
	assert.Equal(t, task3.AbsoluteTimestamp, lastTaskTime)
}

// TestUnsetFlag ensures the correct default value is returned for an unset flag (0), and an error is also returned
func TestUnsetFlag(t *testing.T) {
	st := NewStore()

	// Retrieve unset flag and verify default value and err
	val, err := st.GetFlag(42)
	assert.Nil(t, val)
	assert.Error(t, err)
}

// TestSetFlag ensures the value for a flag is updated properly
func TestSetFlag(t *testing.T) {
	st := NewStore()

	// Set flag
	st.SetFlag(42, 15)
	st.SetFlag(42, false)
	st.SetFlag(42, "k8s")

	// Retrieve unset flag and verify default value and err
	val, err := st.GetFlag(42)
	assert.Equal(t, "k8s", val)
	assert.NoError(t, err)
}

// TestArrayWithNil ensures an array containing nills will not destroy the pq
func TestArrayWithNil(t *testing.T) {
	task1 := &apatelet.Task{AbsoluteTimestamp: 213123}
	task2 := &apatelet.Task{AbsoluteTimestamp: 4242}
	st := NewStore()

	// Enqueue tasks
	st.EnqueueTasks([]*apatelet.Task{nil, task1, nil, task2, nil, nil})

	// Ensure there are two tasks
	assert.Equal(t, 2, st.LenTasks())

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.AbsoluteTimestamp, firstTaskTime)

	// Retrieve first task, and confirm it's task 2
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task2, firstTask)

	// Ensure task 1 is still in the queue
	assert.Equal(t, 1, st.LenTasks())
}
