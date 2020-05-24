package scenario

// PodStatus specifies the pod status
type PodStatus int

const (
	// PodStatusUnset means there is no status set, used as a default
	PodStatusUnset PodStatus = iota
	// PodStatusPending means the pod is in the pending state
	PodStatusPending
	// PodStatusRunning means the pod is in the running state
	PodStatusRunning
	// PodStatusSucceeded means the pod is in the succeeded state
	PodStatusSucceeded
	// PodStatusFailed means the pod is in the failed state
	PodStatusFailed
	// PodStatusUnknown means the pod is in an unknown state
	PodStatusUnknown
)
