// Package events lists all the events supported by the Apatelet
package events

// EventFlag is a flag to be used by the Apatelet
type EventFlag = int32

// NodeEventFlag is a node specific flag to be used by the Apatelet
type NodeEventFlag = EventFlag

const (
	// NodeCreatePodResponse determines how to respond to the CreatePod request
	// Can also be influenced on pod level
	NodeCreatePodResponse NodeEventFlag = iota
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

	// NodeAddedLatencyEnabled defines the amount of added network latency the node experiences
	// Whether to enable added latency
	// If not enabled, the next field will be ignored
	NodeAddedLatencyEnabled

	// NodeAddedLatencyMsec is the amount of msec of latency
	// Will default to 0 msec
	NodeAddedLatencyMsec
)

// PodEventFlag is a pod specific flag to be used by the Apatelet
type PodEventFlag = EventFlag

const (
	// PodCreatePodResponse determines how to respond to the CreatePod request
	// Can also be influenced on node level
	PodCreatePodResponse PodEventFlag = iota

	// PodUpdatePodResponse determines how to respond to the UpdatePod request
	// Can also be influenced on node level
	PodUpdatePodResponse

	// PodDeletePodResponse determines how to respond to the DeletePod request
	// Can also be influenced on node level
	PodDeletePodResponse

	// PodGetPodResponse determines how to respond to the GetPod request
	// Can also be influenced on node level
	PodGetPodResponse

	// PodGetPodStatusResponse determines how to respond to the GetPodStatus request
	// Can also be influenced on node level
	PodGetPodStatusResponse

	// PodMemoryUsage is the amount of bytes of memory used
	// Will default to 0
	PodMemoryUsage

	// PodCPUUsage is the of milli CPUs used
	// Will default to 0
	PodCPUUsage

	// PodStorageUsage is the amount of bytes of storage used
	// Will default to 0
	PodStorageUsage

	// PodEphemeralStorageUsage is the amount of bytes of ephermal storage used
	// Will default to 0
	PodEphemeralStorageUsage

	// PodStatus updates the status of a certain percentage of pods in the current configuration
	// Can be left empty to keep the status unchanged
	// If left empty, the pod_status_percentage will be ignored
	PodStatus
)
