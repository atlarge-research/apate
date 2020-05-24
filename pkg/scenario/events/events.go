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

	// NodeUpdatePodResponse determines how to respond to the UpdatePod request
	// Can also be influenced on pod level
	NodeUpdatePodResponse

	// NodeDeletePodResponse determines how to respond to the DeletePod request
	// Can also be influenced on pod level
	NodeDeletePodResponse

	// NodeGetPodResponse determines how to respond to the GetPod request
	// Can also be influenced on pod level
	NodeGetPodResponse

	// NodeGetPodStatusResponse determines how to respond to the GetPodStatus request
	// Can also be influenced on pod level
	NodeGetPodStatusResponse

	// NodeGetPodsResponse determines how to respond to the GetPod request
	NodeGetPodsResponse

	// NodePingResponse determines how to respond to the Ping request
	NodePingResponse

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

	// PodResources are resources for the pod. See stats.PodStats
	PodResources

	// PodStatus updates the status of a certain percentage of pods in the current configuration
	// Can be left empty to keep the status unchanged
	// If left empty, the pod_status_percentage will be ignored
	PodStatus
)
