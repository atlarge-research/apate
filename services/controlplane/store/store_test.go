package store

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
)

// TestEnqueue ensures adding a single node resources results in a single allocation
func TestEnqueue(t *testing.T) {
	store := NewStore()
	err := store.AddResourcesToQueue([]scenario.NodeResources{{}})
	assert.NoError(t, err)

	// Retrieve only resource
	first, err := store.GetResourceFromQueue()
	assert.NotNil(t, first)
	assert.NoError(t, err)

	// Attempt to get more, should fail
	second, err := store.GetResourceFromQueue()
	assert.Nil(t, second)
	assert.Error(t, err)
}

// TestEmptyGet ensures the store returns an error if no resources were enqueued
func TestEmptyGet(t *testing.T) {
	store := NewStore()

	// Attempt to get non-existing resource, should fail
	res, err := store.GetResourceFromQueue()
	assert.Nil(t, res)
	assert.Error(t, err)
}

// TestEmptyNodeMap ensures there are no nodes after start
func TestEmptyNodeMap(t *testing.T) {
	store := NewStore()

	// Verify there are no nodes by default
	nodes, err := store.GetNodes()
	assert.NoError(t, err)
	assert.Len(t, nodes, 0)
}

// TestAddNodeGet ensures a node can be retrieved by its uuid after it has been added
func TestAddNodeGet(t *testing.T) {
	store := NewStore()

	// Add created node
	id := uuid.New()
	expected := *NewNode(*service.NewConnectionInfo("yeet", 42, false),
		&scenario.NodeResources{UUID: id}, "6")
	err := store.AddNode(&expected)
	assert.NoError(t, err)

	// Verify node was added
	actual, err := store.GetNode(id)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

// TestAddNodeList ensures an added node appears in the list of nodes
func TestAddNodeList(t *testing.T) {
	store := NewStore()
	node := *NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{}, "5")

	err := store.AddNode(&node)
	assert.NoError(t, err)

	list, err := store.GetNodes()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, node, list[0])
}

// TestGetNodeWrongUuid ensures retrieving a node with an unknown uuid results in an error
func TestGetNodeWrongUuid(t *testing.T) {
	store := NewStore()

	// Verify error on node existing node
	node, err := store.GetNode(uuid.New())
	assert.Equal(t, Node{}, node)
	assert.Error(t, err)
}

// TestAddNodeDuplicateUuid ensures that a node with a duplicate uuid will not be aded
func TestAddNodeDuplicateUuid(t *testing.T) {
	store := NewStore()
	expected := NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{}, "4")

	// Add first time
	err := store.AddNode(expected)
	assert.NoError(t, err)

	// Then again with duplicate uuid, should fail
	err = store.AddNode(expected)
	assert.Error(t, err)
}

// TestRemoveNode ensures a removed node is no longer in the list and can no longer be retrieved
func TestRemoveNode(t *testing.T) {
	selector := "1231415"
	store := NewStore()
	node := NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{}, selector)

	// Add node
	err := store.AddNode(node)
	assert.NoError(t, err)

	// Remove node
	err = store.RemoveNode(node.UUID)
	assert.NoError(t, err)

	// Verify there are no nodes left
	list, err := store.GetNodes()
	assert.NoError(t, err)
	assert.Len(t, list, 0)

	// Verify it cannot be retrieved
	res, err := store.GetNode(node.UUID)
	assert.Equal(t, Node{}, res)
	assert.Error(t, err)

	// Verify it's no longer in the selector
	nodes, err := store.GetNodesBySelector(selector)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(nodes))
}

// TestDeleteNoNode ensures removing a node that does not exist keeps the store intact
func TestDeleteNoNode(t *testing.T) {
	store := NewStore()
	node := *NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{}, "3")

	err := store.AddNode(&node)
	assert.NoError(t, err)

	// Remove random node
	err = store.RemoveNode(NewNode(*service.NewConnectionInfo("yeet", 42, false),
		&scenario.NodeResources{UUID: uuid.New()}, "awdaw2").UUID)
	assert.NoError(t, err)

	// Check if the original node is still intact
	list, err := store.GetNodes()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, node, list[0])
}

// TestClearNodes ensures nodes are no longer in the list and can no longer be retrieved when the store is cleared
func TestClearNodes(t *testing.T) {
	selector := "1231145"
	store := NewStore()
	node := NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{}, selector)

	// Add node
	err := store.AddNode(node)
	assert.NoError(t, err)

	// Remove nodes
	err = store.ClearNodes()
	assert.NoError(t, err)

	// Verify there are no nodes left
	list, err := store.GetNodes()
	assert.NoError(t, err)
	assert.Len(t, list, 0)

	// Verify it cannot be retrieved
	res, err := store.GetNode(node.UUID)
	assert.Equal(t, Node{}, res)
	assert.Error(t, err)

	// Verify it's no longer in the selector
	nodes, err := store.GetNodesBySelector(selector)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(nodes))
}

// TestGetNodesBySelector ensures nodes can be retrieved based on their selector
func TestGetNodesBySelector(t *testing.T) {
	selector := "awidya8wdya9wd7iyh"
	store := NewStore()
	node := NewNode(*service.NewConnectionInfo("yeet", 42, false), &scenario.NodeResources{UUID: uuid.New()}, selector)
	node2 := NewNode(*service.NewConnectionInfo("awdadaw", 124, false), &scenario.NodeResources{UUID: uuid.New()}, "123")

	// Add nodes
	err := store.AddNode(node)
	assert.NoError(t, err)
	err = store.AddNode(node2)
	assert.NoError(t, err)

	// Verify it gets returned based on its selector
	nodes, err := store.GetNodesBySelector(selector)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, *node, nodes[0])
}
