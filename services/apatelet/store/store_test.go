package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

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
	st.EnqueueTasks([]*Task{{
		RelativeTimestamp: 0,
		IsPod:             false,
		PodTask:           nil,
		NodeTask:          task,
	}})

	// Retrieve single task and verify it was the original one
	retrieved, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task, retrieved)

	// Also verify it was removed
	assert.Equal(t, 0, st.LenTasks())
}

// TestGetSingleTask ensures a polled task is not deleted
func TestPollSingleTask(t *testing.T) {
	task := &apatelet.Task{RelativeTimestamp: 424242}
	st := NewStore()

	// Enqueue single task
	st.EnqueueTasks([]*Task{{
		424242,
		false, nil, task,
	}})

	// Poll single task and verify the timestamp is correct
	retrieved, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task.RelativeTimestamp, retrieved)

	// Also verify it was not removed
	assert.Equal(t, 1, st.LenTasks())
}

// TestMultipleTasks ensures the priority queue actually sorts the tasks properly (earliest task first)
func TestMultipleTasks(t *testing.T) {
	task1 := &apatelet.Task{RelativeTimestamp: 213123}
	task2 := &apatelet.Task{RelativeTimestamp: 4242}
	task3 := &apatelet.Task{RelativeTimestamp: 83481234}
	st := NewStore()

	// Enqueue tasks
	st.EnqueueTasks([]*Task{
		{
			213123, false, nil, task1,
		},
		{
			4242, false, nil, task2,
		},
		{
			83481234, false, nil, task3,
		},
	})

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.RelativeTimestamp, firstTaskTime)

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
	assert.Equal(t, task3.RelativeTimestamp, lastTaskTime)
}

// TestArrayWithNil ensures an array containing nills will not destroy the pq
func TestArrayWithNil(t *testing.T) {
	task1 := &apatelet.Task{RelativeTimestamp: 213123}
	task2 := &apatelet.Task{RelativeTimestamp: 4242}
	st := NewStore()

	// Enqueue tasks
	st.EnqueueTasks([]*Task{nil, {
		213123, false, nil, task1,
	}, nil, {
		4242, false, nil, task2,
	}, nil, nil})

	// Ensure there are two tasks
	assert.Equal(t, 2, st.LenTasks())

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.RelativeTimestamp, firstTaskTime)

	// Retrieve first task, and confirm it's task 2
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task2, firstTask)

	// Ensure task 1 is still in the queue
	assert.Equal(t, 1, st.LenTasks())
}

// Nodes

// TestUnsetNodeFlag ensures the correct default value is returned for an unset flag (0), and an error is also returned
func TestUnsetNodeFlag(t *testing.T) {
	st := NewStore()

	// Retrieve unset flag and verify default value and err
	val, err := st.GetNodeFlag(42)
	assert.Nil(t, val)
	assert.Error(t, err)
}

// TestSetNodeFlag ensures the value for a flag is updated properly
func TestSetNodeFlag(t *testing.T) {
	st := NewStore()

	// Set flag
	st.SetNodeFlag(42, 15)
	st.SetNodeFlag(42, false)
	st.SetNodeFlag(42, "k8s")

	// Retrieve unset flag and verify default value and err
	val, err := st.GetNodeFlag(42)
	assert.Equal(t, "k8s", val)
	assert.NoError(t, err)

	_, err = st.GetNodeFlag(44)
	assert.Error(t, err, "flag not set")
}

// pods

// TestUnsetPodFlag ensures the correct default value is returned for an unset flag (0), and an error is also returned
func TestUnsetPodFlag(t *testing.T) {
	st := NewStore()

	// Retrieve unset flag and verify default value and err
	val, err := st.GetPodFlag("a", 42)
	assert.Nil(t, val)
	assert.Error(t, err)
}

// TestSetPodFlag ensures the value for a flag is updated properly
func TestSetPodFlag(t *testing.T) {
	st := NewStore()

	// Set flag
	st.SetPodFlag("a", 42, 15)
	st.SetPodFlag("a", 42, false)
	st.SetPodFlag("b", 42, "k8s")

	// Retrieve unset flag and verify default value and err
	val, err := st.GetPodFlag("a", 42)
	assert.Equal(t, false, val)
	assert.NoError(t, err)

	val, err = st.GetPodFlag("b", 42)
	assert.Equal(t, "k8s", val)
	assert.NoError(t, err)

	_, err = st.GetPodFlag("b", 44)
	assert.Error(t, err, "flag not set")
}
