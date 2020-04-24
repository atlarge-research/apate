package normalise

import "github.com/google/uuid"

type NodeResources struct {
	uuid       uuid.UUID
	Ram        int64
	CpuPercent int
	MaxPods    int
}
