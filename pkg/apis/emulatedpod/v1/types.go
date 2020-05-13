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

	// +kubebuilder:default=NORMAL
	CreatePodResponse    EmulatedPodResponse `json:"create_pod_response,omitempty"`

	// +kubebuilder:default=NORMAL
	UpdatePodResponse    EmulatedPodResponse `json:"update_pod_response,omitempty"`

	// +kubebuilder:default=NORMAL
	DeletePodResponse    EmulatedPodResponse `json:"delete_pod_response,omitempty"`

	// +kubebuilder:default=NORMAL
	GetPodResponse       EmulatedPodResponse `json:"get_pod_response,omitempty"`

	// +kubebuilder:default=NORMAL
	GetPodStatusResponse EmulatedPodResponse `json:"get_pod_status_response,omitempty"`

	// +kubebuilder:validation:Optional
	ResourceUsage EmulatedPodResourceUsage `json:"resource_usage,omitempty"`

	// +kubebuilder:default=POD_RUNNING
	PodStatus EmulatedPodStatus `json:"pod_status,omitempty"`
}

type EmulatedPodResourceUsage struct {
	// +kubebuilder:default=0
	Memory           string `json:"memory,omitempty"`

	// +kubebuilder:default=0
	CPU              int64 `json:"cpu,omitempty"`

	// +kubebuilder:default=0
	Storage          string `json:"storage,omitempty"`

	// +kubebuilder:default=0
	EphemeralStorage string `json:"ephemeral_storage,omitempty"`
}

// +kubebuilder:validation:Enum=POD_PENDING;POD_RUNNING;POD_SUCCEEDED;POD_FAILED;POD_UNKNOWN
type EmulatedPodStatus string

// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR
type EmulatedPodResponse string
