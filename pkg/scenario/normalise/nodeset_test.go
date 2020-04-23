package normalise

import (
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"

	"testing"
)

var nodegroup1 = public.NodeGroup{
	Groupname: "test1",
	Nodetype:  "test",
	Amount:    10,
}

var nodegroup2 = public.NodeGroup{
	Groupname: "test2",
	Nodetype:  "test",
	Amount:    10,
}

var nodegroup3 = public.NodeGroup{
	Groupname: "test3",
	Nodetype:  "test",
	Amount:    10,
}

var nodegroups = []*public.NodeGroup{
	&nodegroup1,
	&nodegroup2,
	&nodegroup3,
}

func TestDesugarNodeSetAll(t *testing.T) {
	r, err := desugarNodeSet([]string{
		"test1",
		"test2",
	}, nodegroups)

	assert.NoError(t, err)

	assert.Equal(t, r, []string{
		"test1",
		"test2",
	})
}

func TestDesugarNode(t *testing.T) {
	r, err := desugarNodeSet([]string{
		"all",
	}, nodegroups)

	assert.NoError(t, err)

	assert.Equal(t, r, []string{
		"test1",
		"test2",
		"test3",
	})
}

func TestDesugarNodeSetDuplicate(t *testing.T) {
	_, err := desugarNodeSet([]string{
		"test1",
		"test2",
		"test3",
		"test4",
	}, nodegroups)

	assert.Error(t, err)
}

func TestDesugarNodeSetNotPresent(t *testing.T) {
	_, err := desugarNodeSet([]string{
		"test1",
		"test2",
		"test3",
		"test5",
	}, nodegroups)

	assert.Error(t, err)
}
