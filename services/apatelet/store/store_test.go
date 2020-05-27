package store

import (
	"testing"
	"time"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/stretchr/testify/assert"
)

// TestEmptyQueue ensures the store starts with an empty queue
func TestEmptyQueue(t *testing.T) {
	t.Parallel()

	st := NewStore()

	// Make sure the amount of tasks starts at zero
	assert.Equal(t, 0, st.LenTasks())

	// Make sure both poll and get return an error
	_, _, pollErr := st.PeekTask()
	_, getErr := st.PopTask()
	assert.NoError(t, pollErr)
	assert.Error(t, getErr)
}

// TestGetSingleTask ensures a retrieved task is also deleted
func TestGetSingleTask(t *testing.T) {
	t.Parallel()

	task := &nodeconfigv1.NodeConfigurationState{}
	st := NewStore()

	// Enqueue single task
	err := st.SetNodeTasks([]*Task{NewNodeTask(0, task)})
	assert.NoError(t, err)

	// Retrieve single task and verify it was the original one
	retrieved, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(0, task), retrieved)

	// Also verify it was removed
	assert.Equal(t, 0, st.LenTasks())
}

// TestGetSingleTask ensures a polled task is not deleted
func TestPollSingleTask(t *testing.T) {
	t.Parallel()

	timestamp := time.Duration(424242)
	task := &nodeconfigv1.NodeConfigurationState{}
	st := NewStore()

	// Enqueue single task
	err := st.SetNodeTasks([]*Task{NewNodeTask(timestamp, task)})
	assert.NoError(t, err)

	// Poll single task and verify the timestamp is correct
	retrieved, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, timestamp, retrieved)
	assert.True(t, found)

	// Also verify it was not removed
	assert.Equal(t, 1, st.LenTasks())
}

// TestMultipleTasks ensures the priority queue actually sorts the tasks properly (earliest task first)
func TestMultipleTasks(t *testing.T) {
	t.Parallel()

	task1Time := time.Duration(213123)
	task2Time := time.Duration(4242)
	task3Time := time.Duration(83481234)
	task1 := &nodeconfigv1.NodeConfigurationState{}
	task2 := &nodeconfigv1.NodeConfigurationState{}
	task3 := &nodeconfigv1.NodeConfigurationState{}

	st := NewStore()

	// Enqueue tasks
	err := st.SetNodeTasks([]*Task{
		NewNodeTask(task1Time, task1),
		NewNodeTask(task2Time, task2),
		NewNodeTask(task3Time, task3),
	})
	assert.NoError(t, err)

	// Poll first task, which should be task 2
	firstTaskTime, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2Time, firstTaskTime)
	assert.True(t, found)

	// Retrieve first two tasks
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(4242, task2), firstTask)

	secondTaskTime, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task1Time, secondTaskTime)
	assert.True(t, found)

	secondTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, NewNodeTask(213123, task1), secondTask)

	// Verify there is still one task left
	lastTaskTime, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, 1, st.LenTasks())
	assert.Equal(t, task3Time, lastTaskTime)
	assert.True(t, found)
}

// TestArrayWithNil ensures an array containing nills will not destroy the pq
func TestArrayWithNil(t *testing.T) {
	t.Parallel()

	task1 := NewNodeTask(213123, &nodeconfigv1.NodeConfigurationState{})
	task2 := NewNodeTask(4242, &nodeconfigv1.NodeConfigurationState{})
	st := NewStore()

	// Enqueue tasks
	err := st.SetNodeTasks([]*Task{nil, task1, task2, nil, nil})
	assert.NoError(t, err)

	// Ensure there are two tasks
	assert.Equal(t, 2, st.LenTasks())

	// Peek first task, which should be task 2
	firstTaskTime, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task2.RelativeTimestamp, firstTaskTime)
	assert.True(t, found)

	// Retrieve first task, and confirm it's task 2
	firstTask, err := st.PopTask()
	assert.NoError(t, err)
	assert.Equal(t, task2, firstTask)

	// Ensure task 1 is still in the queue
	assert.Equal(t, 1, st.LenTasks())
}

// TestReplaceNodeTask tests tasks are indeed replaced
func TestReplaceNodeTask(t *testing.T) {
	t.Parallel()

	task1 := NewNodeTask(213123, &nodeconfigv1.NodeConfigurationState{})
	task2 := NewNodeTask(4242, &nodeconfigv1.NodeConfigurationState{})
	task3 := NewNodeTask(53156, &nodeconfigv1.NodeConfigurationState{})
	st := NewStore()

	// Enqueue tasks
	err := st.SetNodeTasks([]*Task{nil, task1, task2, nil, nil})
	assert.NoError(t, err)

	// Ensure there are two tasks
	assert.Equal(t, 2, st.LenTasks())

	// Set new tasks
	err = st.SetNodeTasks([]*Task{task3})
	assert.NoError(t, err)

	// Ensure there is one task
	assert.Equal(t, 1, st.LenTasks())

	// Peek first task, which should be task 3
	firstTaskTime, found, err := st.PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, task3.RelativeTimestamp, firstTaskTime)
	assert.True(t, found)
}

// TestInvalidTaskSet
func TestInvalidTaskSet(t *testing.T) {
	t.Parallel()

	task1 := &Task{
		RelativeTimestamp: 1,
		NodeTask: &NodeTask{
			State: nil,
		},
	}
	st := NewStore()

	// Enqueue something
	err := st.SetNodeTasks([]*Task{{}})
	assert.NoError(t, err)

	// Enqueue node tasks
	err = st.SetNodeTasks([]*Task{task1})
	assert.Error(t, err)

	// Enqueue pod tasks
	err = st.SetPodTasks("", []*Task{task1})
	assert.Error(t, err)
}

// Task

// TestInvalidTask ensures the task functionality works
func TestInvalidTask(t *testing.T) {
	t.Parallel()

	task := Task{}

	_, err := task.IsPod()
	assert.Error(t, err)

	task = Task{
		RelativeTimestamp: 1,
		PodTask: &PodTask{
			Label: "awawd",
		},
		NodeTask: &NodeTask{
			State: nil,
		},
	}

	_, err = task.IsNode()
	assert.Error(t, err)
}

// Nodes

// TestDefaultNodeFlag ensure the default node flags are working properly
func TestDefaultNodeFlag(t *testing.T) {
	t.Parallel()

	st := NewStore()

	// Retrieve default value
	val, err := st.GetNodeFlag(events.NodeCreatePodResponse)
	assert.Equal(t, scenario.ResponseNormal, val)
	assert.NoError(t, err)
}

// TestUnsetNodeFlag ensures the correct default value is returned for an unset flag (0), and an error is also returned
func TestUnsetNodeFlag(t *testing.T) {
	t.Parallel()

	st := NewStore()

	// Retrieve unset flag and verify default value and err
	val, err := st.GetNodeFlag(42)
	assert.Nil(t, val)
	assert.Error(t, err)
}

// TestSetNodeFlag ensures the value for a flag is updated properly
func TestSetNodeFlag(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	st := insertBaselineCRD(t)

	testPeekAndPop(t, &st, 42, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 100, true)
	testPeekAndPop(t, &st, 140, true)
	testPeekAndPop(t, &st, 200, false)
}

func TestEnqueueCRDUpdate(t *testing.T) {
	t.Parallel()

	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding less means old CRDs are removed
	err := st.SetPodTasks("la/clappe", []*Task{
		NewPodTask(10, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(20, "la/clappe", &podconfigv1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	testPeekAndPop(t, &st, 10, true)
	testPeekAndPop(t, &st, 20, true)
	testPeekAndPop(t, &st, 80, false)
	testPeekAndPop(t, &st, 200, false)
}

func TestEnqueueCRDUpdateMore(t *testing.T) {
	t.Parallel()

	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding more means new CRDs are added
	err := st.SetPodTasks("la/clappe", []*Task{
		NewPodTask(10, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(20, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(220, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(120, "la/clappe", &podconfigv1.PodConfigurationState{}),
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
	t.Parallel()

	st := insertBaselineCRD(t)

	// Testing whether updating CRDs works
	// And if adding more means new CRDs are added
	err := st.SetPodTasks("high/tech", []*Task{
		NewPodTask(44, "high/tech", &podconfigv1.PodConfigurationState{}),
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
	t.Parallel()

	st := insertBaselineCRD(t)

	// Testing whether removig CRDs works, even when there are multiple
	err := st.SetPodTasks("high/tech", []*Task{
		NewPodTask(44, "high/tech", &podconfigv1.PodConfigurationState{}),
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

	err := st.SetNodeTasks([]*Task{
		NewNodeTask(80, &nodeconfigv1.NodeConfigurationState{}),
		NewNodeTask(200, &nodeconfigv1.NodeConfigurationState{}),
	})
	assert.NoError(t, err)

	// Testing whether adding new CRDs works
	err = st.SetPodTasks("la/clappe", []*Task{
		NewPodTask(100, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(42, "la/clappe", &podconfigv1.PodConfigurationState{}),
		NewPodTask(140, "la/clappe", &podconfigv1.PodConfigurationState{}),
	})
	assert.NoError(t, err)

	return st
}

func testPeekAndPop(t *testing.T, st *Store, relTime time.Duration, shouldBePod bool) {
	taskTime, found, err := (*st).PeekTask()
	assert.NoError(t, err)
	assert.Equal(t, relTime, taskTime)
	assert.True(t, found)

	task, err := (*st).PopTask()
	assert.NoError(t, err)
	assert.Equal(t, relTime, task.RelativeTimestamp)

	isPod, err := task.IsPod()
	assert.NoError(t, err)
	assert.Equal(t, shouldBePod, isPod)
}

// pods

// TestDefaultPodFlag ensure the default pod flags are working properly
func TestDefaultPodFlag(t *testing.T) {
	t.Parallel()

	st := NewStore()

	// Retrieve default value
	val, err := st.GetPodFlag("", events.PodCreatePodResponse)
	assert.Equal(t, scenario.ResponseUnset, val)
	assert.NoError(t, err)
}

// TestUnsetPodFlag ensures the correct default value is returned for an unset flag (0), and an error is also returned
func TestUnsetPodFlag(t *testing.T) {
	t.Parallel()

	st := NewStore()

	// Retrieve unset flag and verify default value and err
	val, err := st.GetPodFlag("a", 42)
	assert.Nil(t, val)
	assert.Error(t, err)
}

// TestSetPodFlag ensures the value for a flag is updated properly
func TestSetPodFlag(t *testing.T) {
	t.Parallel()

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
