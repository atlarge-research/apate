// Package events lists all the events supported by the Apatelet
package events

// EventFlag is a flag to be used by the Apatelet
type EventFlag int32

const (
	// NodeCreatePodResponse determines how to respond to the CreatePod request
	// Can also be influenced on pod level
	NodeCreatePodResponse EventFlag = iota
	// NodeCreatePodResponsePercentage determines what percentage of CreatePod requests have the above mentioned response
	NodeCreatePodResponsePercentage

	// NodeUpdatePodResponse determines how to respond to the UpdatePod request
	// Can also be influenced on pod level
	NodeUpdatePodResponse
	// NodeUpdatePodResponsePercentage determines what percentage of UpdatePod requests have the above mentioned response
	NodeUpdatePodResponsePercentage

	// NodeDeletePodResponse determines how to respond to the DeletePod request
	// Can also be influenced on pod level
	NodeDeletePodResponse
	// NodeDeletePodResponsePercentage determines what percentage of DeletePod requests have the above mentioned response
	NodeDeletePodResponsePercentage

	// NodeGetPodResponse determines how to respond to the GetPod request
	// Can also be influenced on pod level
	NodeGetPodResponse
	// NodeGetPodResponsePercentage determines what percentage of GetPod requests have the above mentioned response
	NodeGetPodResponsePercentage

	// NodeGetPodStatusResponse determines how to respond to the GetPodStatus request
	// Can also be influenced on pod level
	NodeGetPodStatusResponse
	// NodeGetPodStatusResponsePercentage determines what percentage of GetPodStatus requests have the above mentioned response
	NodeGetPodStatusResponsePercentage

	// NodeGetPodsResponse determines how to respond to the GetPod request
	NodeGetPodsResponse
	// NodeGetPodsResponsePercentage determines what percentage of GetPod requests have the above mentioned response
	NodeGetPodsResponsePercentage

	// NodePingResponse determines how to respond to the Ping request
	NodePingResponse
	// NodePingResponsePercentage determines what percentage of Ping requests have the above mentioned response
	NodePingResponsePercentage

	// NodeEnableResourceAlteration defines the amount of resources the node has in use
	// Whether to enable the next fields
	// If not enabled, the next fields will be ignored
	NodeEnableResourceAlteration

	// NodeMemoryUsage is the amount of bytes of memory used
	// Will default to 0
	NodeMemoryUsage

	// NodeCPUUsage is the of milli CPUs used
	// Will default to 0
	NodeCPUUsage

	// NodeStorageUsage is the amount of bytes of storage used
	// Will default to 0
	NodeStorageUsage

	// NodeEphemeralStorageUsage is the amount of bytes of ephermal storage used
	// Will default to 0
	NodeEphemeralStorageUsage

	// NodeAddedLatencyEnabled defines the amount of added network latency the node experiences
	// Whether to enable added latency
	// If not enabled, the next field will be ignored
	NodeAddedLatencyEnabled

	// NodeAddedLatencyMsec is the amount of msec of latency
	// Will default to 0 msec
	NodeAddedLatencyMsec

	// PodCreatePodResponse determines how to respond to the CreatePod request
	// Can also be influenced on node level
	PodCreatePodResponse
	// PodCreatePodResponsePercentage determines what percentage of CreatePod requests have the above mentioned response
	PodCreatePodResponsePercentage

	// PodUpdatePodResponse determines how to respond to the UpdatePod request
	// Can also be influenced on node level
	PodUpdatePodResponse
	// PodUpdatePodResponsePercentage determines what percentage of UpdatePod requests have the above mentioned response
	PodUpdatePodResponsePercentage

	// PodDeletePodResponse determines how to respond to the DeletePod request
	// Can also be influenced on node level
	PodDeletePodResponse
	// PodDeletePodResponsePercentage determines what percentage of DeletePod requests have the above mentioned response
	PodDeletePodResponsePercentage

	// PodGetPodResponse determines how to respond to the GetPod request
	// Can also be influenced on node level
	PodGetPodResponse
	// PodGetPodResponsePercentage determines what percentage of GetPod requests have the above mentioned response
	PodGetPodResponsePercentage

	// PodGetPodStatusResponse determines how to respond to the GetPodStatus request
	// Can also be influenced on node level
	PodGetPodStatusResponse
	// PodGetPodStatusResponsePercentage determines what percentage of GetPodStatus requests have the above mentioned response
	PodGetPodStatusResponsePercentage

	// PodUpdatePodStatus updates the status of a certain percentage of pods in the current configuration
	// Can be left empty to keep the status unchanged
	// If left empty, the pod_status_percentage will be ignored
	PodUpdatePodStatus

	// PodUpdatePodStatusPercentage sets the percentage of pods in the current deployment getting the new status
	PodUpdatePodStatusPercentage

	// PodUpdateStartTime updates the starting time of this pod (see corev1.PodStatus.StartTime)
	// ISO8601 format
	// Can be left empty to keep the starting time unchanged
	PodUpdateStartTime
)
