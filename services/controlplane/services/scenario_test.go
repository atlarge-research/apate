package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

func TestConvertToAbsolute(t *testing.T) {
	fiveS := time.Second * time.Duration(5)
	minusFiveM := time.Minute * time.Duration(-5)
	fifteenS := time.Second * time.Duration(15)

	input := &apatelet.ApateletScenario{Task: []*apatelet.Task{
		{
			RelativeTimestamp: fiveS.Nanoseconds(),
			AbsoluteTimestamp: 0,
		},
		{
			RelativeTimestamp: minusFiveM.Nanoseconds(),
			AbsoluteTimestamp: 0,
		},
		{
			RelativeTimestamp: fifteenS.Nanoseconds(),
			AbsoluteTimestamp: 0,
		},
	}}

	startTime := time.Now().UnixNano()
	convertToAbsoluteTimestamp(input, startTime)

	assert.Equal(t, startTime+fiveS.Nanoseconds(), input.Task[0].AbsoluteTimestamp)
	assert.Equal(t, int64(0), input.Task[0].RelativeTimestamp)

	assert.Equal(t, startTime+minusFiveM.Nanoseconds(), input.Task[1].AbsoluteTimestamp)
	assert.Equal(t, int64(0), input.Task[1].RelativeTimestamp)

	assert.Equal(t, startTime+fifteenS.Nanoseconds(), input.Task[2].AbsoluteTimestamp)
	assert.Equal(t, int64(0), input.Task[2].RelativeTimestamp)
}

func TestFilter(t *testing.T) {
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	uuid3 := uuid.New().String()

	input := []*apatelet.Task{
		{
			NodeSet: map[string]bool{
				uuid1: true,
				uuid2: true,
				uuid3: true,
			},
		},
		{
			NodeSet: map[string]bool{
				uuid2: true,
				uuid3: true,
			},
		},
		{
			NodeSet: map[string]bool{
				uuid1: true,
				uuid3: true,
			},
		},
	}

	uuid1Tasks := filterTasksForNode(input, uuid1)
	assert.EqualValues(t, []*apatelet.Task{
		{
			NodeSet: map[string]bool{
				uuid1: true,
				uuid2: true,
				uuid3: true,
			},
		},
		{
			NodeSet: map[string]bool{
				uuid1: true,
				uuid3: true,
			},
		},
	}, uuid1Tasks)

	uuid2Tasks := filterTasksForNode(input, uuid2)
	assert.EqualValues(t, []*apatelet.Task{
		{
			NodeSet: map[string]bool{
				uuid1: true,
				uuid2: true,
				uuid3: true,
			},
		},
		{
			NodeSet: map[string]bool{
				uuid2: true,
				uuid3: true,
			},
		},
	}, uuid2Tasks)

	uuid3Tasks := filterTasksForNode(input, uuid3)
	assert.EqualValues(t, input, uuid3Tasks)
}
