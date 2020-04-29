package normalization

import (
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"testing"
)

var nodegroup1 = controlplane.NodeGroup{
	GroupName: "test1",
	NodeType:  "test",
	Amount:    10,
}

var nodegroup2 = controlplane.NodeGroup{
	GroupName: "test2",
	NodeType:  "test",
	Amount:    10,
}

var nodegroup3 = controlplane.NodeGroup{
	GroupName: "test3",
	NodeType:  "test",
	Amount:    10,
}

var nodegroups = []*controlplane.NodeGroup{
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

	assert.Equal(t, []string{
		"test1",
		"test2",
	}, r)
}

func TestDesugarNode(t *testing.T) {
	r, err := desugarNodeGroups([]string{
		"all",
	}, nodegroups)

	assert.NoError(t, err)

	assert.Equal(t, []string{
		"test1",
		"test2",
		"test3",
	}, r)
}

func TestDesugarNodeSetDuplicate(t *testing.T) {
	r, err := desugarNodeGroups([]string{
		"test1",
		"test2",
		"test3",
		"test3",
	}, nodegroups)

	assert.NoError(t, err)

	assert.Equal(t, []string{
		"test1",
		"test2",
		"test3",
	}, r)
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
