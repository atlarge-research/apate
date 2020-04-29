package normalization

import "github.com/google/uuid"

// NodeResources describe the resources of a single node, including the UUID of that node
type NodeResources struct {
	// The UUID of the node
	UUID uuid.UUID

	// The amount of bytes of memory
	Memory int64

	// The amount of milli CPUs in Kubernetes
	CPU int64

	// The max amount of pods in Kubernetes
	MaxPods int64
}
