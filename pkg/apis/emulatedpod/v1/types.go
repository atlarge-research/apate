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
	UsePodStartTime bool `json:"use_pod_start_time"`

	Tasks []EmulatedPodTask `json:"tasks,omitempty"`
}

type EmulatedPodTask struct {
	// +kubebuilder:validation:Minimum=0
	Timestamp int64 `json:"timestamp"`

	// +kubebuilder:default=UNSET
	CreatePodResponse EmulatedPodResponse `json:"create_pod_response,omitempty"`

	// +kubebuilder:default=UNSET
	UpdatePodResponse EmulatedPodResponse `json:"update_pod_response,omitempty"`

	// +kubebuilder:default=UNSET
	DeletePodResponse EmulatedPodResponse `json:"delete_pod_response,omitempty"`

	// +kubebuilder:default=UNSET
	GetPodResponse EmulatedPodResponse `json:"get_pod_response,omitempty"`

	// +kubebuilder:default=UNSET
	GetPodStatusResponse EmulatedPodResponse `json:"get_pod_status_response,omitempty"`

	// +kubebuilder:validation:Optional
	ResourceUsage *EmulatedPodResourceUsage `json:"resource_usage,omitempty"`

	// +kubebuilder:default=UNSET
	PodStatus EmulatedPodStatus `json:"pod_status,omitempty"`
}

type EmulatedPodResourceUsage struct {
	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:default=-1
	Memory string `json:"memory,omitempty"`

	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:default=-1
	CPU int64 `json:"cpu,omitempty"`

	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:default=-1
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:default=-1
	EphemeralStorage string `json:"ephemeral_storage,omitempty"`
}

// +kubebuilder:validation:Enum=POD_PENDING;POD_RUNNING;POD_SUCCEEDED;POD_FAILED;POD_UNKNOWN;UNSET
type EmulatedPodStatus string
const (
	POD_STATUS_PENDING EmulatedPodResponse = "POD_PENDING"
	POD_STATUS_RUNNING EmulatedPodResponse = "POD_RUNNING"
	POD_STATUS_SUCCEEDED EmulatedPodResponse = "POD_SUCCEEDED"
	POD_STATUS_FAILED EmulatedPodResponse = "POD_FAILED"
	POD_STATUS_UNKNOWN EmulatedPodStatus = "POD_UNKNOWN"
	POD_STATUS_UNSET EmulatedPodStatus = "UNSET"
)

// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR;UNSET
type EmulatedPodResponse string
const (
	RESPONSE_NORMAL EmulatedPodResponse = "NORMAL"
	RESPONSE_TIMEOUT EmulatedPodResponse = "TIMEOUT"
	RESPONSE_ERROR EmulatedPodResponse = "ERROR"
	RESPONSE_UNSET EmulatedPodResponse = "UNSET"
)
