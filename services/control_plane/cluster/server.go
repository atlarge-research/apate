package cluster

import (
	"sync"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

//TODO: Add more information for node

// Node represents a virtual kubelet in the Apate cluster
type Node struct {
	connectionInfo service.ConnectionInfo
	UUID           uuid.UUID
}

//TODO: Multi-master soon :tm:

// ApateCluster represents the entire apate cluster
type ApateCluster struct {
	nodes    map[uuid.UUID]*Node
	nodeLock sync.RWMutex
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo) *Node {
	return &Node{
		connectionInfo: info,
		UUID:           uuid.New(),
	}
}

// NewApateCluster creates a new empty cluster
func NewApateCluster() *ApateCluster {
	return &ApateCluster{
		nodes: make(map[uuid.UUID]*Node),
	}
}

// AddNode adds the given Node to the apate cluster
func (c *ApateCluster) AddNode(node *Node) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	c.nodes[node.UUID] = node
}

// RemoveNode removes the given Node from the apate cluster
func (c *ApateCluster) RemoveNode(node *Node) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	delete(c.nodes, node.UUID)
}

// GetNode returns the node with the given uuid
func (c *ApateCluster) GetNode(uuid uuid.UUID) *Node {
	c.nodeLock.RLock()
	defer c.nodeLock.RUnlock()

	return c.nodes[uuid]
}

// GetNodes returns an array containing all nodes in the apate cluster
func (c *ApateCluster) GetNodes() []*Node {
	c.nodeLock.RLock()
	defer c.nodeLock.RUnlock()

	nodes := make([]*Node, 0, len(c.nodes))

	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}
