package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

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
