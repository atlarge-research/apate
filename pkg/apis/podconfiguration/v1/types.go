package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PodConfigurationLabel defines the label which is used to define which pod configuration resource belongs to the pod
	PodConfigurationLabel = "apate"
)

// PodConfiguration is a definition of PodConfiguration resource.
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=podconfigurations,shortName=pc,singular=podconfiguration
type PodConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Spec PodConfigurationSpec `json:"spec"`
}

// PodConfigurationList is a list of PodConfigurations.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PodConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PodConfiguration `json:"items"`
}

// PodConfigurationSpec is the spec which belongs to the PodConfiguration CRD
type PodConfigurationSpec struct {
	// A direct way to update state, this will circumvent the timestamps / scenario
	// +kubebuilder:validation:Optional
	PodConfigurationState `json:",inline,omitempty"`

	// The tasks to be executed
	// +kubebuilder:validation:Optional
	Tasks []PodConfigurationTask `json:"tasks,omitempty"`
}

// PodConfigurationTask is a single task which updates a pod state and is executed at a timestamp
type PodConfigurationTask struct {
	// The timestamp at which the task is executed
	// Any time.ParseDuration format is accepted, such as "10ms" or "42s"
	// +kubebuilder:validation:Required
	Timestamp string `json:"timestamp"`

	// Indicates whether the timestamp is relative to the start of the pod or not.
	// If set to true, a timestamp of 10s means this task will be executed 10 seconds after the pod started.
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	RelativeToPod bool `json:"relative_to_pod"`

	// The state to be set
	// +kubebuilder:validation:Required
	State PodConfigurationState `json:"state"`
}

// PodConfigurationState is the state to be set for the related pods
type PodConfigurationState struct {
	// CreatePodResponse determines how to respond to the CreatePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	CreatePodResponse PodResponse `json:"create_pod_response,omitempty"`

	// UpdatePodResponse determines how to respond to the UpdatePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	UpdatePodResponse PodResponse `json:"update_pod_response,omitempty"`

	// DeletePodResponse determines how to respond to the DeletePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	DeletePodResponse PodResponse `json:"delete_pod_response,omitempty"`

	// PodGetPodResponse determines how to respond to the GetPod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	GetPodResponse PodResponse `json:"get_pod_response,omitempty"`

	// GetPodStatusResponse determines how to respond to the GetPodStatus request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	GetPodStatusResponse PodResponse `json:"get_pod_status_response,omitempty"`

	// PodResources sets the amount of resources the related pods are using
	// +kubebuilder:validation:Optional
	PodResources *PodResources `json:"pod_resources,omitempty"`

	// PodStatus updates the current pod status
	// +kubebuilder:default=UNSET
	PodStatus PodStatus `json:"pod_status,omitempty"`
}

// PodResources defines the current resource usage of the pod
type PodResources struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0B"
	Memory string `json:"memory,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=0
	CPU uint64 `json:"cpu,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0B"
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0B"
	EphemeralStorage string `json:"ephemeral_storage,omitempty"`
}

// PodStatus can be PENDING, RUNNING, SUCCEEDED, FAILED, UNKNOWN or UNSET, and describes the state of a pod.
// +kubebuilder:validation:Enum=PENDING;RUNNING;SUCCEEDED;FAILED;UNKNOWN;UNSET
type PodStatus string

// Enum variants for PodStatus
const (
	PodStatusPending   PodStatus = "PENDING"
	PodStatusRunning   PodStatus = "RUNNING"
	PodStatusSucceeded PodStatus = "SUCCEEDED"
	PodStatusFailed    PodStatus = "FAILED"
	PodStatusUnknown   PodStatus = "UNKNOWN"
	PodStatusUnset     PodStatus = "UNSET"
)

// PodResponse can be NORMAL, TIMEOUT, ERROR or UNSET, and describes how a pod should respond
// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR;UNSET
type PodResponse string

// Enum variants for PodResponse
const (
	ResponseNormal  PodResponse = "NORMAL"
	ResponseTimeout PodResponse = "TIMEOUT"
	ResponseError   PodResponse = "ERROR"
	ResponseUnset   PodResponse = "UNSET"
)
