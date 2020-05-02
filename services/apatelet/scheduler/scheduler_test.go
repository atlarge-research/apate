package scheduler

//go:generate sh -c "cd ../../../ && make mockgen"

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"testing"
)

func TestTaskHandlerSimple(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	a, _ := any.Marshal(scenario.Response_ERROR)

	// Test task:
	task := apatelet.Task{
		EventFlags: map[int32]*anypb.Any{
			events.NodeCreatePodResponse: a,
		},
	}

	// Set up expectations
	ms.EXPECT().SetFlag(events.NodeCreatePodResponse, scenario.Response_ERROR)

	var s store.Store = ms
	sched := Scheduler{&s}

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, &task)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestTaskHandlerMultiple(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	m1, err := any.Marshal(scenario.Response_ERROR)
	assert.NoError(t, err)
	m2, err := any.Marshal(42)
	assert.NoError(t, err)

	// Test task:
	task := apatelet.Task{
		EventFlags: map[int32]*anypb.Any{
			events.NodeCreatePodResponse: m1,
			events.NodeAddedLatencyMsec:  m2,
		},
	}

	// Set up expectations
	ms.EXPECT().SetFlag(events.NodeCreatePodResponse, scenario.Response_ERROR)
	ms.EXPECT().SetFlag(events.NodeAddedLatencyMsec, int64(42))

	var s store.Store = ms
	sched := Scheduler{&s}

	// Run code under test
	ech := make(chan error)

	sched.taskHandler(ech, &task)

	select {
	case <-ech:
		t.Fail()
	default:
	}
}

func TestTaskHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	a := anypb.Any{
		TypeUrl: "invalid-any",
		Value:   nil,
	}

	// Test task:
	task := apatelet.Task{
		EventFlags: map[int32]*anypb.Any{
			events.NodeCreatePodResponse: &a,
		},
	}

	// Set up expectations

	var s store.Store = ms
	sched := Scheduler{&s}

	// Run code under test
	ech := make(chan error, 1)

	sched.taskHandler(ech, &task)

	select {
	case err := <-ech:
		assert.Error(t, err)
	default:
		t.Fail()
	}
}
