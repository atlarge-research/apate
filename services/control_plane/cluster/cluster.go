package cluster

import (
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

//TODO: Add more information for node

// Node represents a virtual kubelet in the Apate cluster
type Node struct {
	connectionInfo service.ConnectionInfo
	UUID           uuid.UUID
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo) *Node {
	return &Node{
		connectionInfo: info,
		UUID:           uuid.New(),
	}
}

//TODO: Multi-master soon :tm:

// ApateCluster represents the entire apate cluster
type ApateCluster interface {
	// AddNode adds the given Node to the apate cluster
	AddNode(*Node) error

	// RemoveNode removes the given Node from the apate cluster
	RemoveNode(*Node) error

	// GetNode returns the node with the given uuid
	GetNode(uuid.UUID) (Node, error)

	// GetNodes returns an array containing all nodes in the apate cluster
	GetNodes() ([]Node, error)
}
