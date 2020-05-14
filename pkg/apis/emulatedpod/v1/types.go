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

	// +kubebuilder:validation:Optional
	DirectState EmulatedPodState `json:"direct_task,omitempty"`

	// +kubebuilder:validation:Optional
	Tasks []EmulatedPodTask `json:"tasks,omitempty"`
}

type EmulatedPodTask struct {
	// +kubebuilder:validation:Minimum=0
	Timestamp int64 `json:"timestamp"`

	// +kubebuilder:validation:Required
	State EmulatedPodState `json:"state"`
}

type EmulatedPodState struct {
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
	// +kubebuilder:default="-1B"
	Memory string `json:"memory,omitempty"`

	// +kubebuilder:default="-1B"
	CPU int64 `json:"cpu,omitempty"`

	// +kubebuilder:default="-1B"
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:default="-1B"
	EphemeralStorage string `json:"ephemeral_storage,omitempty"`
}

// +kubebuilder:validation:Enum=PENDING;RUNNING;SUCCEEDED;FAILED;UNKNOWN;UNSET
type EmulatedPodStatus string

const (
	POD_STATUS_PENDING   EmulatedPodStatus = "PENDING"
	POD_STATUS_RUNNING   EmulatedPodStatus = "RUNNING"
	POD_STATUS_SUCCEEDED EmulatedPodStatus = "SUCCEEDED"
	POD_STATUS_FAILED    EmulatedPodStatus = "FAILED"
	POD_STATUS_UNKNOWN   EmulatedPodStatus = "UNKNOWN"
	POD_STATUS_UNSET     EmulatedPodStatus = "UNSET"
)

// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR;UNSET
type EmulatedPodResponse string

const (
	RESPONSE_NORMAL  EmulatedPodResponse = "NORMAL"
	RESPONSE_TIMEOUT EmulatedPodResponse = "TIMEOUT"
	RESPONSE_ERROR   EmulatedPodResponse = "ERROR"
	RESPONSE_UNSET   EmulatedPodResponse = "UNSET"
)
