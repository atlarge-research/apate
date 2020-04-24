package store

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"

	"github.com/google/uuid"
)

//TODO: Add more information for node

// Node represents a virtual kubelet in the Apate cluster
type Node struct {
	ConnectionInfo service.ConnectionInfo
	UUID           uuid.UUID
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo) *Node {
	return &Node{
		ConnectionInfo: info,
		UUID:           uuid.New(),
	}
}
