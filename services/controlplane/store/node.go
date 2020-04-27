package store

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"

	"github.com/google/uuid"
)

//TODO: Add more information for node

// Node represents a Apatelet in the Apate cluster
type Node struct {
	ConnectionInfo service.ConnectionInfo
	UUID           uuid.UUID
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo, resources *normalization.NodeResources) *Node {
	return &Node{
		ConnectionInfo: info,
		UUID:           resources.UUID,
	}
}
