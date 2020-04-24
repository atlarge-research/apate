// Package cluster provides state to the apate cluster
package cluster

import (
	"fmt"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalise"
	"sync"

	"github.com/google/uuid"
)

//TODO: Multi-master soon :tm:

// Store represents the entire apate cluster
type Store interface {
	// AddNode adds the given Node to the apate cluster
	AddNode(*Node) error

	// RemoveNode removes the given Node from the apate cluster
	RemoveNode(*Node) error

	// GetNode returns the node with the given uuid
	GetNode(uuid.UUID) (Node, error)

	// GetNodes returns an array containing all nodes in the apate cluster
	GetNodes() ([]Node, error)

	// ClearNodes removes all nodes from the apate cluster
	ClearNodes() error

	// AddResourceQueue adds a node resource to the queue
	AddResourceQueue(normalise.NodeResources) error
}

type apateClusterStandalone struct {
	nodes    map[uuid.UUID]Node
	nodeLock sync.RWMutex
}

// NewApateCluster creates a new empty cluster
func NewApateCluster() Store {
	return &apateClusterStandalone{
		nodes: make(map[uuid.UUID]Node),
	}
}

// AddNode adds the given Node to the apate cluster
func (c *apateClusterStandalone) AddNode(node *Node) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	// Check if node already exists (uuid collision)
	if _, ok := c.nodes[node.UUID]; ok {
		return fmt.Errorf("node with uuid '%s' already exists", node.UUID.String())
	}

	c.nodes[node.UUID] = *node

	return nil
}

// RemoveNode removes the given Node from the apate cluster
func (c *apateClusterStandalone) RemoveNode(node *Node) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	delete(c.nodes, node.UUID)
	return nil
}

// GetNode returns the node with the given uuid
func (c *apateClusterStandalone) GetNode(uuid uuid.UUID) (Node, error) {
	c.nodeLock.RLock()
	defer c.nodeLock.RUnlock()

	if node, ok := c.nodes[uuid]; ok {
		return node, nil
	}

	return Node{}, fmt.Errorf("node with uuid '%s' not found", uuid.String())
}

// GetNodes returns an array containing all nodes in the apate cluster
func (c *apateClusterStandalone) GetNodes() ([]Node, error) {
	c.nodeLock.RLock()
	defer c.nodeLock.RUnlock()

	nodes := make([]Node, 0, len(c.nodes))

	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (c *apateClusterStandalone) ClearNodes() error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	c.nodes = make(map[uuid.UUID]Node)
	return nil
}

func (c *apateClusterStandalone) AddResourceQueue(resources normalise.NodeResources) error {
	panic("awdawdadawda")
}
