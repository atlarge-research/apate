package store

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/google/uuid"
)

// Node represents a Apatelet in the Apate cluster
type Node struct {
	ConnectionInfo service.ConnectionInfo
	UUID           uuid.UUID
	Status         health.Status
	Label          string
}

// NewNode creates a new Node based on the given connection information
func NewNode(info service.ConnectionInfo, resources *scenario.NodeResources, label string) *Node {
	return &Node{
		ConnectionInfo: info,
		UUID:           resources.UUID,
		Status:         health.Status_UNKNOWN,
		Label:          label,
	}
}
