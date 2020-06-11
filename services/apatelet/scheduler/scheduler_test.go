package scheduler

//go:generate sh -c "cd ../../../ && make mockgen"

import (
	"context"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandlerSimpleNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := &nodeconfigv1.NodeConfigurationState{

		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	// Set up expectations
	ms.EXPECT().SetNodeFlags(store.Flags{events.NodeCreatePodResponse: scenario.ResponseError})

	var s store.Store = ms
	sched := New(&s)

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, store.NewNodeTask(0, task))

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestTaskHandlerSimplePod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NewPodTask(42, "la/clappe", &podconfigv1.PodConfigurationState{
		PodStatus: podconfigv1.PodStatusFailed,
	})

	// Set up expectations
	ms.EXPECT().SetPodFlags("la/clappe", store.Flags{events.PodStatus: scenario.PodStatusFailed})

	var s store.Store = ms
	sched := New(&s)

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, task)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestTaskHandlerMultiple(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := nodeconfigv1.NodeConfigurationState{
		NetworkLatency: "42s",
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	// Set up expectations
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil).AnyTimes()
	ms.EXPECT().SetNodeFlags(store.Flags{
		events.NodeCreatePodResponse: scenario.ResponseError,
		events.NodeAddedLatency:      42 * time.Second,
	})

	var s store.Store = ms
	sched := New(&s)

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, store.NewNodeTask(0, &task))

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunner(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := nodeconfigv1.NodeConfigurationState{
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(0), true, nil).AnyTimes()
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlags(gomock.Any()).AnyTimes()

	var s store.Store = ms
	sched := New(&s)
	sched.WakeScheduler()

	// Run code under test
	ech := make(chan error, 1)

	sched.runner(ech)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunnerDelay(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := nodeconfigv1.NodeConfigurationState{
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	secondDelay := time.Second * 5

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(0), true, nil)
	ms.EXPECT().PeekTask().Return(secondDelay, true, nil)
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlags(gomock.Any()).AnyTimes()

	var s store.Store = ms
	sched := New(&s)
	sched.startTime = time.Now()

	// Run code under test
	ech := make(chan error, 1)

	done, delay := sched.runner(ech)
	assert.Equal(t, false, done)

	maxDelay := sched.startTime.Add(secondDelay).Sub(sched.startTime)
	assert.Less(t, delay.Nanoseconds(), maxDelay.Nanoseconds())

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunnerSleep(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := nodeconfigv1.NodeConfigurationState{
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(0), true, nil)
	ms.EXPECT().PeekTask().Return(time.Duration(0), false, nil).Times(3)
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlags(gomock.Any()).AnyTimes()

	var s store.Store = ms
	sched := New(&s)
	ech := sched.EnableScheduler(context.Background())

	sched.StartScheduler(0)
	time.Sleep(time.Millisecond * 500)
	sched.WakeScheduler()

	time.Sleep(time.Second)
	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunnerDontHandleOldTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(10), true, nil).AnyTimes()
	ms.EXPECT().PopTask().Return(store.NewPodTask(10, "la/clappe", &podconfigv1.PodConfigurationState{
		PodStatus: podconfigv1.PodStatusUnknown,
	}), nil)

	var s store.Store = ms
	sched := Scheduler{
		store:     &s,
		readyCh:   make(chan struct{}),
		prevT:     time.Unix(0, 40),
		startTime: time.Unix(0, 0),
	}

	// Run code under test
	ech := make(chan error, 1)

	sched.runner(ech)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunnerEarlyFail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(-1), true, errors.New("some error"))

	// Setup
	var s store.Store = ms
	sched := New(&s)

	// Test
	ech := make(chan error, 1)
	sched.runner(ech)

	// Verify
	select {
	case err := <-ech:
		assert.Error(t, err)
	default:
		t.Fail()
	}
}

func TestRunnerPopFail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(0), true, nil)
	ms.EXPECT().PopTask().Return(nil, errors.New("new error"))

	// Setup
	var s store.Store = ms
	sched := New(&s)

	// Test
	ech := make(chan error, 1)
	sched.runner(ech)

	// Verify
	select {
	case err := <-ech:
		assert.Error(t, err)
	default:
		t.Fail()
	}
}

func TestStartScheduler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := nodeconfigv1.NodeConfigurationState{
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse: nodeconfigv1.ResponseError,
		},
	}

	// Expectations
	ms.EXPECT().PeekTask().Return(time.Duration(0), true, nil)
	ms.EXPECT().PeekTask().Return(time.Duration(0), false, nil).AnyTimes()
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlags(gomock.Any()).AnyTimes()

	var s store.Store = ms
	sched := New(&s)

	// Run code under test
	ech := sched.EnableScheduler(ctx)
	sched.StartScheduler(0)

	time.Sleep(time.Second)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}
