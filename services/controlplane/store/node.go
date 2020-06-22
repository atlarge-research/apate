package store

import (
	"github.com/atlarge-research/apate/api/health"
	"github.com/atlarge-research/apate/internal/service"
	"github.com/atlarge-research/apate/pkg/scenario"

	"github.com/google/uuid"
)

// Node represents a Apatelet in the Apate cluster
type Node struct {
	ConnectionInfo service.ConnectionInfo
	UUID           uuid.UUID
	Status         health.Status
	Label          string
	Resources      *scenario.NodeResources
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo, resources *scenario.NodeResources, label string) *Node {
	return &Node{
		ConnectionInfo: info,
		UUID:           resources.UUID,
		Status:         health.Status_UNKNOWN,
		Label:          label,
		Resources:      resources,
	}
}
