package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EmulatedPod is a definition of EmulatedPod resource.
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=emulatedpods,shortName=ep,singular=emulatedpod
type EmulatedPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Spec EmulatedPodSpec `json:"spec"`
}

// EmulatedPodList is a list of EmulatedPods.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EmulatedPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EmulatedPod `json:"items"`
}

// EmulatedPodSpec is the spec which belongs to the EmulatedPod CRD
type EmulatedPodSpec struct {
	// A direct way to update state, this will circumvent the timestamps / scenario
	// +kubebuilder:validation:Optional
	DirectState EmulatedPodState `json:"direct_task,omitempty"`

	// The tasks to be executed
	// +kubebuilder:validation:Optional
	Tasks []EmulatedPodTask `json:"tasks,omitempty"`
}

// EmulatedPodTask is a single task which updates a pod state and is executed at a timestamp
type EmulatedPodTask struct {
	// The timestamp at which the task is executed
	// +kubebuilder:validation:Minimum=0
	Timestamp int64 `json:"timestamp"`

	// The state to be set
	// +kubebuilder:validation:Required
	State EmulatedPodState `json:"state"`
}

// EmulatedPodState is the state to be set for the related pods
type EmulatedPodState struct {
	// CreatePodResponse determines how to respond to the CreatePod request
	// +kubebuilder:default=UNSET
	CreatePodResponse EmulatedPodResponse `json:"create_pod_response,omitempty"`

	// UpdatePodResponse determines how to respond to the UpdatePod request
	// +kubebuilder:default=UNSET
	UpdatePodResponse EmulatedPodResponse `json:"update_pod_response,omitempty"`

	// DeletePodResponse determines how to respond to the DeletePod request
	// +kubebuilder:default=UNSET
	DeletePodResponse EmulatedPodResponse `json:"delete_pod_response,omitempty"`

	// PodGetPodResponse determines how to respond to the GetPod request
	// +kubebuilder:default=UNSET
	GetPodResponse EmulatedPodResponse `json:"get_pod_response,omitempty"`

	// GetPodStatusResponse determines how to respond to the GetPodStatus request
	// +kubebuilder:default=UNSET
	GetPodStatusResponse EmulatedPodResponse `json:"get_pod_status_response,omitempty"`

	// PodResources sets the amount of resources the related pods are using
	// +kubebuilder:validation:Optional
	PodResources *PodStats `json:"pod_resources,omitempty"`

	// PodStatus updates the current pod status
	// +kubebuilder:default=UNSET
	PodStatus EmulatedPodStatus `json:"pod_status,omitempty"`
}

// EmulatedPodStatus can be PENDING, RUNNING, SUCCEEDED, FAILED, UNKNOWN or UNSET, and describes the state of a pod.
// +kubebuilder:validation:Enum=PENDING;RUNNING;SUCCEEDED;FAILED;UNKNOWN;UNSET
type EmulatedPodStatus string

// Enum variants for EmulatedPodStatus
const (
	PodStatusPending   EmulatedPodStatus = "PENDING"
	PodStatusRunning   EmulatedPodStatus = "RUNNING"
	PodStatusSucceeded EmulatedPodStatus = "SUCCEEDED"
	PodStatusFailed    EmulatedPodStatus = "FAILED"
	PodStatusUnknown   EmulatedPodStatus = "UNKNOWN"
	PodStatusUnset     EmulatedPodStatus = "UNSET"
)

// EmulatedPodResponse can be NORMAL, TIMEOUT, ERROR or UNSET, and describes how a pod should respond
// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR;UNSET
type EmulatedPodResponse string

// Enum variants for EmulatedPodResponse
const (
	ResponseNormal  EmulatedPodResponse = "NORMAL"
	ResponseTimeout EmulatedPodResponse = "TIMEOUT"
	ResponseError   EmulatedPodResponse = "ERROR"
	ResponseUnset   EmulatedPodResponse = "UNSET"
)
