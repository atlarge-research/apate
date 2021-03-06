package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// NodeConfigurationLabel defines the label which can used to find out to which node configuration crd a node belongs
	// by external tools
	NodeConfigurationLabel = "apate-name"
	// NodeConfigurationLabelNamespace defines the namespace of the above label
	NodeConfigurationLabelNamespace = "apate-namespace"

	// EmulatedLabel is used by taints and node labels
	EmulatedLabel = "emulated"
	// EmulatedLabelValue is the value associated with the EmulatedLabel
	EmulatedLabelValue = "yes"
	// NodeIDLabel defines the label which can be used to get the uuid of the node
	NodeIDLabel = "apate-uuid"
)

// NodeConfiguration is a definition of a NodeConfiguration resource
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=nodeconfigurations,shortName=nc,singular=nodeconfiguration
type NodeConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Spec NodeConfigurationSpec `json:"spec"`
}

// NodeConfigurationList is a list of NodeConfiguration
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodeConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NodeConfiguration `json:"items"`
}

// NodeConfigurationSpec is the spec which belongs to NodeConfiguration
type NodeConfigurationSpec struct {
	// A way to directly update the node state
	// +kubebuilder:validation:Optional
	NodeConfigurationState `json:",inline,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int64 `json:"replicas"`

	// +kubebuilder:validation:Required
	Resources NodeResources `json:"resources"`

	// The tasks to be executed on this node
	// +kubebuilder:validation:Optional
	Tasks []NodeConfigurationTask `json:"tasks,omitempty"`
}

// NodeResources specifies the resources the node has available
type NodeResources struct {
	// +kubebuilder:validation:Required
	Memory string `json:"memory,omitempty"`

	// +kubebuilder:validation:Required
	CPU int64 `json:"cpu,omitempty"`

	// +kubebuilder:validation:Required
	Storage string `json:"storage,omitempty"`

	// +kubebuilder:validation:Required
	EphemeralStorage string `json:"ephemeral_storage,omitempty"`

	// +kubebuilder:validation:Required
	MaxPods int64 `json:"max_pods,omitempty"`
}

// NodeConfigurationTask is a single task which modifies the node state on the given timestamp
type NodeConfigurationTask struct {
	// The timestamp at which the task is executed
	// Any time.ParseDuration format is accepted, such as "10ms" or "42s"
	// +kubebuilder:validation:Required
	Timestamp string `json:"timestamp"`

	// The desired state of the node after this task
	// +kubebuilder:validation:Required
	State NodeConfigurationState `json:"state"`
}

// NodeConfigurationState is the state of the node, used for determining how to respond to request from kubernetes.
// This state includes some built-in states, which Apate will translate to direct state for ease of use.
// Said built-in states take precedence over the custom state
type NodeConfigurationState struct {
	// If set, NodeFailed will result in timeouts for all requests by kubernetes
	// effectively taking down the node
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	NodeFailed bool `json:"node_failed,omitempty"`

	// NetworkLatency determines how much added latency will be introduced to requests by kubernetes.
	// Any time.ParseDuration format is accepted, such as "10ms" or "42s"
	// The default is unset. Any invalid or negative integer will also be interpreted as unset.
	// +kubebuilder:default=unset
	// +kubebuilder:validation:Optional
	NetworkLatency string `json:"network_latency,omitempty"`

	// If set, HeartbeatFailed will result in the node no longer responding to pings
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	HeartbeatFailed bool `json:"heartbeat_failed,omitempty"`

	// CustomState specifies a custom state
	// +kubebuilder:validation:Optional
	CustomState *NodeConfigurationCustomState `json:"custom_state,omitempty"`
}

// NodeConfigurationCustomState is the state of the node, used for determining how to respond to request from kubernetes.
// This state will not be translated or anything similar, as this is a direct mapping to the actual state of the apatelet
type NodeConfigurationCustomState struct {
	// CreatePodResponse determines how to respond to the CreatePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	CreatePodResponse NodeResponse `json:"create_pod_response,omitempty"`

	// UpdatePodResponse determines how to respond to the UpdatePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	UpdatePodResponse NodeResponse `json:"update_pod_response,omitempty"`

	// DeletePodResponse determines how to respond to the DeletePod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	DeletePodResponse NodeResponse `json:"delete_pod_response,omitempty"`

	// PodGetPodResponse determines how to respond to the GetPod request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	GetPodResponse NodeResponse `json:"get_pod_response,omitempty"`

	// GetPodsResponse determines how to respond to the GetPods request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	GetPodsResponse NodeResponse `json:"get_pods_response,omitempty"`

	// GetPodStatusResponse determines how to respond to the GetPodStatus request
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	GetPodStatusResponse NodeResponse `json:"get_pod_status_response,omitempty"`

	// NodePingResponse determines how to respond to a heartbeat ping
	// +kubebuilder:default=UNSET
	// +kubebuilder:validation:Optional
	NodePingResponse NodeResponse `json:"node_ping_response,omitempty"`
}

// NodeResponse can be NORMAL, TIMEOUT, ERROR or UNSET, and describes how a node should respond to a pod related request
// +kubebuilder:validation:Enum=NORMAL;TIMEOUT;ERROR;UNSET
type NodeResponse string

// Enum variants for PodResponse
const (
	ResponseNormal  NodeResponse = "NORMAL"
	ResponseTimeout NodeResponse = "TIMEOUT"
	ResponseError   NodeResponse = "ERROR"
	ResponseUnset   NodeResponse = "UNSET"
)
