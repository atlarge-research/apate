package scheduler

//go:generate sh -c "cd ../../../ && make mockgen"

import (
	"context"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	nodeV1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandlerSimpleNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NodeTask{
		State: &nodeV1.NodeConfigurationState{
			CustomState: &nodeV1.NodeConfigurationDirectState{
				CreatePodResponse: nodeV1.ResponseError,
			},
		},
	}

	// Set up expectations
	ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, scenario.ResponseError)
	ms.EXPECT().SetNodeFlag(events.NodeAddedLatencyMsec, int64(0))

	var s store.Store = ms
	sched := New(context.Background(), &s)

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, store.NewNodeTask(0, &task))

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestTaskHandlerSimplePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NewPodTask(42, "la/clappe", &v1.PodConfigurationState{
		PodStatus: v1.PodStatusFailed,
	})

	// Set up expectations
	ms.EXPECT().SetPodFlag("la/clappe", events.PodStatus, scenario.PodStatusFailed)

	var s store.Store = ms
	sched := New(context.Background(), &s)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NodeTask{
		State: &nodeV1.NodeConfigurationState{
			NetworkLatency: 42,
			CustomState: &nodeV1.NodeConfigurationDirectState{
				CreatePodResponse: nodeV1.ResponseError,
			},
		},
	}

	// Set up expectations
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(0, nil).AnyTimes()
	ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, scenario.ResponseError)
	ms.EXPECT().SetNodeFlag(events.NodeAddedLatencyMsec, int64(42))

	var s store.Store = ms
	sched := New(context.Background(), &s)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NodeTask{
		State: &nodeV1.NodeConfigurationState{
			CustomState: &nodeV1.NodeConfigurationDirectState{
				CreatePodResponse: nodeV1.ResponseError,
			},
		},
	}

	// Expectations
	ms.EXPECT().PeekTask().Return(int64(0), nil)
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlag(gomock.Any(), gomock.Any()).AnyTimes()

	var s store.Store = ms
	sched := New(context.Background(), &s)

	// Run code under test
	ech := make(chan error, 1)

	sched.runner(ech)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestRunnerDontHandleOldTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(int64(10), nil)
	ms.EXPECT().PopTask().Return(store.NewPodTask(10, "la/clappe", &v1.PodConfigurationState{
		PodStatus: v1.PodStatusUnknown,
	}), nil)

	var s store.Store = ms
	sched := Scheduler{
		store:     &s,
		ctx:       context.Background(),
		readyCh:   make(chan struct{}),
		prevT:     40,
		startTime: 0,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(int64(-1), errors.New("some error"))

	// Setup
	var s store.Store = ms
	sched := New(context.Background(), &s)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().PeekTask().Return(int64(0), nil)
	ms.EXPECT().PopTask().Return(nil, errors.New("new error"))

	// Setup
	var s store.Store = ms
	sched := New(context.Background(), &s)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Test task:
	task := store.NodeTask{
		State: &nodeV1.NodeConfigurationState{
			CustomState: &nodeV1.NodeConfigurationDirectState{
				CreatePodResponse: nodeV1.ResponseError,
			},
		},
	}

	// Expectations
	ms.EXPECT().PeekTask().Return(int64(0), nil)
	ms.EXPECT().PopTask().Return(store.NewNodeTask(0, &task), nil)
	ms.EXPECT().SetNodeFlag(gomock.Any(), gomock.Any()).AnyTimes()

	// any further peeks are well into the future
	ms.EXPECT().PeekTask().Return(time.Now().Add(time.Hour*12).UnixNano(), nil).AnyTimes()

	var s store.Store = ms
	sched := New(ctx, &s)

	// Run code under test
	ech := sched.EnableScheduler()
	sched.StartScheduler(0)

	time.Sleep(time.Second)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}
