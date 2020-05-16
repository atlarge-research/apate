package store

import (
	"testing"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

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
	st.EnqueueTasks([]*Task{NewNodeTask(0, task)})

	// Retrieve single task and verify it was the original one
	retrieved, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(0, task), retrieved)

	// Also verify it was removed
	assert.Equal(t, 0, st.LenTasks())
}

// TestGetSingleTask ensures a polled task is not deleted
func TestPollSingleTask(t *testing.T) {
	task := &apatelet.Task{RelativeTimestamp: 424242}
	st := NewStore()

	// Enqueue single task
	st.EnqueueTasks([]*Task{NewNodeTask(424242, task)})

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
	st.SetStartTime(100)

	// Enqueue tasks
	st.EnqueueTasks([]*Task{
		NewNodeTask(213123, task1),
		NewNodeTask(4242, task2),
		NewNodeTask(83481234, task3),
	})

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.RelativeTimestamp+100, firstTaskTime)

	// Retrieve first two tasks
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(4242, task2), firstTask)

	secondTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task1.RelativeTimestamp+100, secondTaskTime)

	secondTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(213123, task1), secondTask)

	// Verify there is still one task left
	lastTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, 1, st.LenTasks())
	assert.Equal(t, task3.RelativeTimestamp+100, lastTaskTime)
}

// TestArrayWithNil ensures an array containing nills will not destroy the pq
func TestArrayWithNil(t *testing.T) {
	task1 := &apatelet.Task{RelativeTimestamp: 213123}
	task2 := &apatelet.Task{RelativeTimestamp: 4242}
	st := NewStore()

	// Enqueue tasks
	st.EnqueueTasks([]*Task{nil, NewNodeTask(213123, task1), NewNodeTask(4242, task2), nil, nil})

	// Ensure there are two tasks
	assert.Equal(t, 2, st.LenTasks())

	// Poll first task, which should be task 2
	firstTaskTime, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.RelativeTimestamp, firstTaskTime)

	// Retrieve first task, and confirm it's task 2
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(4242, task2), firstTask)

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

func TestEnqueueCRDBasic(t *testing.T) {
	st := insertBaselineCRD(t)

	testPeekAndPop(t, &st, 42, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 100, true)
	testPeekAndPop(t, &st, 140, true)
	testPeekAndPop(t, &st, 200, false)
}

func TestEnqueueCRDUpdate(t *testing.T) {
	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding less means old CRDs are removed
	err := st.EnqueuePodTasks("la/clappe", []*Task{
		NewPodTask(10, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(20, "la/clappe", &v1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	testPeekAndPop(t, &st, 10, true)
	testPeekAndPop(t, &st, 20, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 200, false)
}

func TestEnqueueCRDUpdateMore(t *testing.T) {
	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding more means new CRDs are added
	err := st.EnqueuePodTasks("la/clappe", []*Task{
		NewPodTask(10, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(20, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(220, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(120, "la/clappe", &v1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	testPeekAndPop(t, &st, 10, true)
	testPeekAndPop(t, &st, 20, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 120, true)
	testPeekAndPop(t, &st, 200, false)
	testPeekAndPop(t, &st, 220, true)
}

func TestEnqueueCRDNewLabel(t *testing.T) {
	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding more means new CRDs are added
	err := st.EnqueuePodTasks("high/tech", []*Task{
		NewPodTask(44, "high/tech", &v1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	testPeekAndPop(t, &st, 42, true)
	testPeekAndPop(t, &st, 44, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 100, true)
	testPeekAndPop(t, &st, 140, true)
	testPeekAndPop(t, &st, 200, false)
}

func TestRemoveCRD(t *testing.T) {
	st := insertBaselineCRD(t)

	// Testing whether removig CRDs works, even when there are multiple
	err := st.EnqueuePodTasks("high/tech", []*Task{
		NewPodTask(44, "high/tech", &v1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	err = st.RemovePodTasks("la/clappe")
	assert.NoError(t, err)

	testPeekAndPop(t, &st, 44, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 200, false)
}

func insertBaselineCRD(t *testing.T) Store {
	st := NewStore()

	st.EnqueueTasks([]*Task{
		NewNodeTask(80, &apatelet.Task{}),
		NewNodeTask(200, &apatelet.Task{}),
	})

	// Testing whether adding new CRDs works
	err := st.EnqueuePodTasks("la/clappe", []*Task{
		NewPodTask(100, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(42, "la/clappe", &v1.PodConfigurationState{}),
		NewPodTask(140, "la/clappe", &v1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	return st
}

func testPeekAndPop(t *testing.T, st *Store, relTime int64, shouldBePod bool) {
	taskTime, err := (*st).PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, relTime, taskTime)

	task, err := (*st).PopTask()
	assert.NoError(t, err)
	assert.Equal(t, relTime, task.RelativeTimestamp)

	isPod, err := task.IsPod()
	assert.NoError(t, err)
	assert.Equal(t, shouldBePod, isPod)
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
