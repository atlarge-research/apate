package normalise

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/stretchr/testify/assert"

	"testing"
)

var nodegroup1 = control_plane.NodeGroup{
	GroupName: "test1",
	NodeType:  "test",
	Amount:    10,
}

var nodegroup2 = control_plane.NodeGroup{
	GroupName: "test2",
	NodeType:  "test",
	Amount:    10,
}

var nodegroup3 = control_plane.NodeGroup{
	GroupName: "test3",
	NodeType:  "test",
	Amount:    10,
}

var nodegroups = []*control_plane.NodeGroup{
	&nodegroup1,
	&nodegroup2,
	&nodegroup3,
}

func TestDesugarNodeSetAll(t *testing.T) {
	r, err := desugarNodeGroups([]string{
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
	r, err := desugarNodeGroups([]string{
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
	_, err := desugarNodeGroups([]string{
		"test1",
		"test2",
		"test3",
		"test3",
	}, nodegroups)

	assert.Error(t, err)
}

func TestDesugarNodeSetNotPresent(t *testing.T) {
	_, err := desugarNodeGroups([]string{
		"test1",
		"test2",
		"test3",
		"test5",
	}, nodegroups)

	assert.Error(t, err)
}

func TestDesugarNodeSetMultipleAll(t *testing.T) {
	_, err := desugarNodeGroups([]string{
		"all",
		"test1",
	}, nodegroups)

	assert.Error(t, err)
}
