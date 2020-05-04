// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: controlplane/scenario.proto

package controlplane

import (
	context "context"
	events "github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Scenario represents a full scenario as given in the scenario configuration file.
// It can be easily constructed from yaml or json (pkg/scenario/deserialize) and has to be
// converted to a private scenario (pkg/scenario/normalize).
type PublicScenario struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of node types.
	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	// List of node groups.
	NodeGroups []*NodeGroup `protobuf:"bytes,2,rep,name=node_groups,json=nodeGroups,proto3" json:"node_groups,omitempty"`
	// A scenario consists of an number of tasks.
	Tasks []*Task `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *PublicScenario) Reset() {
	*x = PublicScenario{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_scenario_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicScenario) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicScenario) ProtoMessage() {}

func (x *PublicScenario) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_scenario_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicScenario.ProtoReflect.Descriptor instead.
func (*PublicScenario) Descriptor() ([]byte, []int) {
	return file_controlplane_scenario_proto_rawDescGZIP(), []int{0}
}

func (x *PublicScenario) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *PublicScenario) GetNodeGroups() []*NodeGroup {
	if x != nil {
		return x.NodeGroups
	}
	return nil
}

func (x *PublicScenario) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

// A task is a part of a scenario describing what happens at a certain point in time.
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	/// The name of a task can be used to revert it.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Time offset after the start of the scenario
	// Specify as string with unit like 10s, 2h, 20ms.
	// If no unit is given, defaults to seconds.
	Time string `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// Revert the task
	Revert bool `protobuf:"varint,3,opt,name=revert,proto3" json:"revert,omitempty"`
	// Which nodes this task applies to
	// This field will be ignored when revert is true
	NodeGroups []string `protobuf:"bytes,4,rep,name=node_groups,json=nodeGroups,proto3" json:"node_groups,omitempty"`
	// The node event to be executed
	// This field will be ignored when revert is true
	//
	// Types that are assignable to NodeEvent:
	//	*Task_NodeFailure
	//	*Task_NetworkLatency
	//	*Task_TimeoutKeepHeartbeat
	//	*Task_NoTimeoutNoHeartbeat
	//	*Task_NodeResponseState
	//	*Task_ResourcePressure
	NodeEvent isTask_NodeEvent `protobuf_oneof:"node_event"`
	// The pod events to be executed
	// This field will be ignored when revert is true
	PodConfigs []*PodConfig `protobuf:"bytes,11,rep,name=pod_configs,json=podConfigs,proto3" json:"pod_configs,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_scenario_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_scenario_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_controlplane_scenario_proto_rawDescGZIP(), []int{1}
}

func (x *Task) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Task) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *Task) GetRevert() bool {
	if x != nil {
		return x.Revert
	}
	return false
}

func (x *Task) GetNodeGroups() []string {
	if x != nil {
		return x.NodeGroups
	}
	return nil
}

func (m *Task) GetNodeEvent() isTask_NodeEvent {
	if m != nil {
		return m.NodeEvent
	}
	return nil
}

func (x *Task) GetNodeFailure() *events.NodeFailure {
	if x, ok := x.GetNodeEvent().(*Task_NodeFailure); ok {
		return x.NodeFailure
	}
	return nil
}

func (x *Task) GetNetworkLatency() *events.NetworkLatency {
	if x, ok := x.GetNodeEvent().(*Task_NetworkLatency); ok {
		return x.NetworkLatency
	}
	return nil
}

func (x *Task) GetTimeoutKeepHeartbeat() *events.TimeoutKeepHeartbeat {
	if x, ok := x.GetNodeEvent().(*Task_TimeoutKeepHeartbeat); ok {
		return x.TimeoutKeepHeartbeat
	}
	return nil
}

func (x *Task) GetNoTimeoutNoHeartbeat() *events.NoTimeoutNoHeartbeat {
	if x, ok := x.GetNodeEvent().(*Task_NoTimeoutNoHeartbeat); ok {
		return x.NoTimeoutNoHeartbeat
	}
	return nil
}

func (x *Task) GetNodeResponseState() *events.ResponseState {
	if x, ok := x.GetNodeEvent().(*Task_NodeResponseState); ok {
		return x.NodeResponseState
	}
	return nil
}

func (x *Task) GetResourcePressure() *events.ResourcePressure {
	if x, ok := x.GetNodeEvent().(*Task_ResourcePressure); ok {
		return x.ResourcePressure
	}
	return nil
}

func (x *Task) GetPodConfigs() []*PodConfig {
	if x != nil {
		return x.PodConfigs
	}
	return nil
}

type isTask_NodeEvent interface {
	isTask_NodeEvent()
}

type Task_NodeFailure struct {
	NodeFailure *events.NodeFailure `protobuf:"bytes,5,opt,name=node_failure,json=nodeFailure,proto3,oneof"`
}

type Task_NetworkLatency struct {
	NetworkLatency *events.NetworkLatency `protobuf:"bytes,6,opt,name=network_latency,json=networkLatency,proto3,oneof"`
}

type Task_TimeoutKeepHeartbeat struct {
	TimeoutKeepHeartbeat *events.TimeoutKeepHeartbeat `protobuf:"bytes,7,opt,name=timeout_keep_heartbeat,json=timeoutKeepHeartbeat,proto3,oneof"`
}

type Task_NoTimeoutNoHeartbeat struct {
	NoTimeoutNoHeartbeat *events.NoTimeoutNoHeartbeat `protobuf:"bytes,8,opt,name=no_timeout_no_heartbeat,json=noTimeoutNoHeartbeat,proto3,oneof"`
}

type Task_NodeResponseState struct {
	NodeResponseState *events.ResponseState `protobuf:"bytes,9,opt,name=node_response_state,json=nodeResponseState,proto3,oneof"`
}

type Task_ResourcePressure struct {
	ResourcePressure *events.ResourcePressure `protobuf:"bytes,10,opt,name=resource_pressure,json=resourcePressure,proto3,oneof"`
}

func (*Task_NodeFailure) isTask_NodeEvent() {}

func (*Task_NetworkLatency) isTask_NodeEvent() {}

func (*Task_TimeoutKeepHeartbeat) isTask_NodeEvent() {}

func (*Task_NoTimeoutNoHeartbeat) isTask_NodeEvent() {}

func (*Task_NodeResponseState) isTask_NodeEvent() {}

func (*Task_ResourcePressure) isTask_NodeEvent() {}

type PodConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metadata.name from a configuration
	MetadataName string `protobuf:"bytes,1,opt,name=metadata_name,json=metadataName,proto3" json:"metadata_name,omitempty"`
	// Types that are assignable to PodEvent:
	//	*PodConfig_PodResponseState
	//	*PodConfig_PodStatusUpdate
	PodEvent isPodConfig_PodEvent `protobuf_oneof:"pod_event"`
}

func (x *PodConfig) Reset() {
	*x = PodConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_scenario_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodConfig) ProtoMessage() {}

func (x *PodConfig) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_scenario_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodConfig.ProtoReflect.Descriptor instead.
func (*PodConfig) Descriptor() ([]byte, []int) {
	return file_controlplane_scenario_proto_rawDescGZIP(), []int{2}
}

func (x *PodConfig) GetMetadataName() string {
	if x != nil {
		return x.MetadataName
	}
	return ""
}

func (m *PodConfig) GetPodEvent() isPodConfig_PodEvent {
	if m != nil {
		return m.PodEvent
	}
	return nil
}

func (x *PodConfig) GetPodResponseState() *events.ResponseState {
	if x, ok := x.GetPodEvent().(*PodConfig_PodResponseState); ok {
		return x.PodResponseState
	}
	return nil
}

func (x *PodConfig) GetPodStatusUpdate() *events.PodStatusUpdate {
	if x, ok := x.GetPodEvent().(*PodConfig_PodStatusUpdate); ok {
		return x.PodStatusUpdate
	}
	return nil
}

type isPodConfig_PodEvent interface {
	isPodConfig_PodEvent()
}

type PodConfig_PodResponseState struct {
	PodResponseState *events.ResponseState `protobuf:"bytes,2,opt,name=pod_response_state,json=podResponseState,proto3,oneof"`
}

type PodConfig_PodStatusUpdate struct {
	PodStatusUpdate *events.PodStatusUpdate `protobuf:"bytes,3,opt,name=pod_status_update,json=podStatusUpdate,proto3,oneof"`
}

func (*PodConfig_PodResponseState) isPodConfig_PodEvent() {}

func (*PodConfig_PodStatusUpdate) isPodConfig_PodEvent() {}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the node. This can be referred to in nodegroups.
	NodeType string `protobuf:"bytes,1,opt,name=node_type,json=nodeType,proto3" json:"node_type,omitempty"`
	// The amount of memory a node gets.
	// Specify as string with unit like 12G, 42M, 200K, or in bytes (without unit)
	Memory string `protobuf:"bytes,2,opt,name=memory,proto3" json:"memory,omitempty"`
	// The amount of milli CPUs allocated to the node
	Cpu int64 `protobuf:"varint,3,opt,name=cpu,proto3" json:"cpu,omitempty"`
	// The amount of storage a node gets.
	// Specify as string with unit like 12G, 42M, 200K, or in bytes (without unit)
	Storage string `protobuf:"bytes,4,opt,name=storage,proto3" json:"storage,omitempty"`
	// The amount of ephemeral storage a node gets.
	// Specify as string with unit like 12G, 42M, 200K, or in bytes (without unit)
	EphemeralStorage string `protobuf:"bytes,5,opt,name=ephemeral_storage,json=ephemeralStorage,proto3" json:"ephemeral_storage,omitempty"`
	// Maximum number of pods per node.
	MaxPods int64 `protobuf:"varint,6,opt,name=max_pods,json=maxPods,proto3" json:"max_pods,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_scenario_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_scenario_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_controlplane_scenario_proto_rawDescGZIP(), []int{3}
}

func (x *Node) GetNodeType() string {
	if x != nil {
		return x.NodeType
	}
	return ""
}

func (x *Node) GetMemory() string {
	if x != nil {
		return x.Memory
	}
	return ""
}

func (x *Node) GetCpu() int64 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *Node) GetStorage() string {
	if x != nil {
		return x.Storage
	}
	return ""
}

func (x *Node) GetEphemeralStorage() string {
	if x != nil {
		return x.EphemeralStorage
	}
	return ""
}

func (x *Node) GetMaxPods() int64 {
	if x != nil {
		return x.MaxPods
	}
	return 0
}

// A node group specifies a group of many nodes with the same properties.
// A nodegroup refers to a node and how many times that type of node is needed.
type NodeGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the group. This can be referred to in scenarios.
	GroupName string `protobuf:"bytes,1,opt,name=group_name,json=groupName,proto3" json:"group_name,omitempty"`
	// The type of node in this group.
	NodeType string `protobuf:"bytes,2,opt,name=node_type,json=nodeType,proto3" json:"node_type,omitempty"`
	// How many times you want this type of node.
	Amount int32 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *NodeGroup) Reset() {
	*x = NodeGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_scenario_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeGroup) ProtoMessage() {}

func (x *NodeGroup) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_scenario_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeGroup.ProtoReflect.Descriptor instead.
func (*NodeGroup) Descriptor() ([]byte, []int) {
	return file_controlplane_scenario_proto_rawDescGZIP(), []int{4}
}

func (x *NodeGroup) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

func (x *NodeGroup) GetNodeType() string {
	if x != nil {
		return x.NodeType
	}
	return ""
}

func (x *NodeGroup) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

var File_controlplane_scenario_proto protoreflect.FileDescriptor

var file_controlplane_scenario_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x73,
	0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61,
	0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6f, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb0, 0x01, 0x0a, 0x0e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x53,
	0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x12, 0x2e, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x3e, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61,
	0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x0a, 0x6e, 0x6f, 0x64,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x2e, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0xe3, 0x05, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x76, 0x65,
	0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x73, 0x12, 0x4b, 0x0a, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x48,
	0x00, 0x52, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x54,
	0x0a, 0x0f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63,
	0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4c, 0x61, 0x74, 0x65, 0x6e,
	0x63, 0x79, 0x48, 0x00, 0x52, 0x0e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4c, 0x61, 0x74,
	0x65, 0x6e, 0x63, 0x79, 0x12, 0x67, 0x0a, 0x16, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f,
	0x6b, 0x65, 0x65, 0x70, 0x5f, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x4b, 0x65, 0x65, 0x70, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x48, 0x00, 0x52, 0x14, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x4b, 0x65, 0x65, 0x70, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x68, 0x0a,
	0x17, 0x6e, 0x6f, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f, 0x6e, 0x6f, 0x5f, 0x68,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f,
	0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x4e, 0x6f, 0x54, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x4e, 0x6f, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x48,
	0x00, 0x52, 0x14, 0x6e, 0x6f, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x4e, 0x6f, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x5a, 0x0a, 0x13, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x00,
	0x52, 0x11, 0x6e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x5a, 0x0a, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x70, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b,
	0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x50, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65, 0x48, 0x00, 0x52, 0x10, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65, 0x12,
	0x3e, 0x0a, 0x0b, 0x70, 0x6f, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x0b,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x50, 0x6f, 0x64, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x0a, 0x70, 0x6f, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x42,
	0x0c, 0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0xf1, 0x01,
	0x0a, 0x09, 0x50, 0x6f, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x23, 0x0a, 0x0d, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x58, 0x0a, 0x12, 0x70, 0x6f, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61,
	0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x10, 0x70, 0x6f, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x58, 0x0a, 0x11, 0x70, 0x6f,
	0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x50, 0x6f, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x6f, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x70, 0x6f, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x22, 0xaf, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f,
	0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e,
	0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x63, 0x70,
	0x75, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x65,
	0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61,
	0x6c, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x5f,
	0x70, 0x6f, 0x64, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d, 0x61, 0x78, 0x50,
	0x6f, 0x64, 0x73, 0x22, 0x5f, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x32, 0x9b, 0x01, 0x0a, 0x08, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69,
	0x6f, 0x12, 0x4c, 0x0a, 0x0c, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69,
	0x6f, 0x12, 0x22, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x53, 0x63, 0x65,
	0x6e, 0x61, 0x72, 0x69, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x41, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65,
	0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controlplane_scenario_proto_rawDescOnce sync.Once
	file_controlplane_scenario_proto_rawDescData = file_controlplane_scenario_proto_rawDesc
)

func file_controlplane_scenario_proto_rawDescGZIP() []byte {
	file_controlplane_scenario_proto_rawDescOnce.Do(func() {
		file_controlplane_scenario_proto_rawDescData = protoimpl.X.CompressGZIP(file_controlplane_scenario_proto_rawDescData)
	})
	return file_controlplane_scenario_proto_rawDescData
}

var file_controlplane_scenario_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_controlplane_scenario_proto_goTypes = []interface{}{
	(*PublicScenario)(nil),              // 0: apate.controlplane.PublicScenario
	(*Task)(nil),                        // 1: apate.controlplane.Task
	(*PodConfig)(nil),                   // 2: apate.controlplane.PodConfig
	(*Node)(nil),                        // 3: apate.controlplane.Node
	(*NodeGroup)(nil),                   // 4: apate.controlplane.NodeGroup
	(*events.NodeFailure)(nil),          // 5: apate.controlplane.events.NodeFailure
	(*events.NetworkLatency)(nil),       // 6: apate.controlplane.events.NetworkLatency
	(*events.TimeoutKeepHeartbeat)(nil), // 7: apate.controlplane.events.TimeoutKeepHeartbeat
	(*events.NoTimeoutNoHeartbeat)(nil), // 8: apate.controlplane.events.NoTimeoutNoHeartbeat
	(*events.ResponseState)(nil),        // 9: apate.controlplane.events.ResponseState
	(*events.ResourcePressure)(nil),     // 10: apate.controlplane.events.ResourcePressure
	(*events.PodStatusUpdate)(nil),      // 11: apate.controlplane.events.PodStatusUpdate
	(*empty.Empty)(nil),                 // 12: google.protobuf.Empty
}
var file_controlplane_scenario_proto_depIdxs = []int32{
	3,  // 0: apate.controlplane.PublicScenario.nodes:type_name -> apate.controlplane.Node
	4,  // 1: apate.controlplane.PublicScenario.node_groups:type_name -> apate.controlplane.NodeGroup
	1,  // 2: apate.controlplane.PublicScenario.tasks:type_name -> apate.controlplane.Task
	5,  // 3: apate.controlplane.Task.node_failure:type_name -> apate.controlplane.events.NodeFailure
	6,  // 4: apate.controlplane.Task.network_latency:type_name -> apate.controlplane.events.NetworkLatency
	7,  // 5: apate.controlplane.Task.timeout_keep_heartbeat:type_name -> apate.controlplane.events.TimeoutKeepHeartbeat
	8,  // 6: apate.controlplane.Task.no_timeout_no_heartbeat:type_name -> apate.controlplane.events.NoTimeoutNoHeartbeat
	9,  // 7: apate.controlplane.Task.node_response_state:type_name -> apate.controlplane.events.ResponseState
	10, // 8: apate.controlplane.Task.resource_pressure:type_name -> apate.controlplane.events.ResourcePressure
	2,  // 9: apate.controlplane.Task.pod_configs:type_name -> apate.controlplane.PodConfig
	9,  // 10: apate.controlplane.PodConfig.pod_response_state:type_name -> apate.controlplane.events.ResponseState
	11, // 11: apate.controlplane.PodConfig.pod_status_update:type_name -> apate.controlplane.events.PodStatusUpdate
	0,  // 12: apate.controlplane.Scenario.loadScenario:input_type -> apate.controlplane.PublicScenario
	12, // 13: apate.controlplane.Scenario.startScenario:input_type -> google.protobuf.Empty
	12, // 14: apate.controlplane.Scenario.loadScenario:output_type -> google.protobuf.Empty
	12, // 15: apate.controlplane.Scenario.startScenario:output_type -> google.protobuf.Empty
	14, // [14:16] is the sub-list for method output_type
	12, // [12:14] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_controlplane_scenario_proto_init() }
func file_controlplane_scenario_proto_init() {
	if File_controlplane_scenario_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controlplane_scenario_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicScenario); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_controlplane_scenario_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_controlplane_scenario_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_controlplane_scenario_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_controlplane_scenario_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeGroup); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_controlplane_scenario_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Task_NodeFailure)(nil),
		(*Task_NetworkLatency)(nil),
		(*Task_TimeoutKeepHeartbeat)(nil),
		(*Task_NoTimeoutNoHeartbeat)(nil),
		(*Task_NodeResponseState)(nil),
		(*Task_ResourcePressure)(nil),
	}
	file_controlplane_scenario_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*PodConfig_PodResponseState)(nil),
		(*PodConfig_PodStatusUpdate)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_controlplane_scenario_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_controlplane_scenario_proto_goTypes,
		DependencyIndexes: file_controlplane_scenario_proto_depIdxs,
		MessageInfos:      file_controlplane_scenario_proto_msgTypes,
	}.Build()
	File_controlplane_scenario_proto = out.File
	file_controlplane_scenario_proto_rawDesc = nil
	file_controlplane_scenario_proto_goTypes = nil
	file_controlplane_scenario_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ScenarioClient is the client API for Scenario service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ScenarioClient interface {
	LoadScenario(ctx context.Context, in *PublicScenario, opts ...grpc.CallOption) (*empty.Empty, error)
	StartScenario(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
}

type scenarioClient struct {
	cc grpc.ClientConnInterface
}

func NewScenarioClient(cc grpc.ClientConnInterface) ScenarioClient {
	return &scenarioClient{cc}
}

func (c *scenarioClient) LoadScenario(ctx context.Context, in *PublicScenario, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.controlplane.Scenario/loadScenario", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scenarioClient) StartScenario(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.controlplane.Scenario/startScenario", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScenarioServer is the server API for Scenario service.
type ScenarioServer interface {
	LoadScenario(context.Context, *PublicScenario) (*empty.Empty, error)
	StartScenario(context.Context, *empty.Empty) (*empty.Empty, error)
}

// UnimplementedScenarioServer can be embedded to have forward compatible implementations.
type UnimplementedScenarioServer struct {
}

func (*UnimplementedScenarioServer) LoadScenario(context.Context, *PublicScenario) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadScenario not implemented")
}
func (*UnimplementedScenarioServer) StartScenario(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartScenario not implemented")
}

func RegisterScenarioServer(s *grpc.Server, srv ScenarioServer) {
	s.RegisterService(&_Scenario_serviceDesc, srv)
}

func _Scenario_LoadScenario_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicScenario)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScenarioServer).LoadScenario(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.controlplane.Scenario/LoadScenario",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScenarioServer).LoadScenario(ctx, req.(*PublicScenario))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scenario_StartScenario_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScenarioServer).StartScenario(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.controlplane.Scenario/StartScenario",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScenarioServer).StartScenario(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Scenario_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.controlplane.Scenario",
	HandlerType: (*ScenarioServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "loadScenario",
			Handler:    _Scenario_LoadScenario_Handler,
		},
		{
			MethodName: "startScenario",
			Handler:    _Scenario_StartScenario_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controlplane/scenario.proto",
}
