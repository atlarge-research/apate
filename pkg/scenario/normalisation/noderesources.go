package normalisation

import "github.com/google/uuid"

// NodeResources describe the resources of a single node, including the UUID of that node
type NodeResources struct {
	UUID       uuid.UUID
	RAM        int64
	CPUPercent int
	MaxPods    int
}
