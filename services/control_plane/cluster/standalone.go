// Package cluster provides a way to maintain state for the apate cluster, and other various services
// related to the apate cluster
package cluster

import (
	"sync"

	"github.com/google/uuid"
)

type apateClusterStandalone struct {
	nodes    map[uuid.UUID]Node
	nodeLock sync.RWMutex
}

// NewApateCluster creates a new empty cluster
func NewApateCluster() ApateCluster {
	return &apateClusterStandalone{
		nodes: make(map[uuid.UUID]Node),
	}
}

// AddNode adds the given Node to the apate cluster
func (c *apateClusterStandalone) AddNode(node *Node) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

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

	return c.nodes[uuid], nil
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
