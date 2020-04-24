// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: control_plane/scenario.proto

package control_plane

import (
	context "context"
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
		mi := &file_control_plane_scenario_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicScenario) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicScenario) ProtoMessage() {}

func (x *PublicScenario) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_scenario_proto_msgTypes[0]
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
	return file_control_plane_scenario_proto_rawDescGZIP(), []int{0}
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
	// This field should not be set when revert is true
	NodeGroups []string `protobuf:"bytes,4,rep,name=node_groups,json=nodeGroups,proto3" json:"node_groups,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_plane_scenario_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_scenario_proto_msgTypes[1]
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
	return file_control_plane_scenario_proto_rawDescGZIP(), []int{1}
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

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the node. This can be referred to in nodegroups.
	NodeType string `protobuf:"bytes,1,opt,name=node_type,json=nodeType,proto3" json:"node_type,omitempty"`
	// The amount of ram a node gets.
	// Specify as string with unit like 12G, 42M, 200K, or in bytes (without unit)
	Ram string `protobuf:"bytes,2,opt,name=ram,proto3" json:"ram,omitempty"`
	// Percentage of cpu allocated to a node.
	CpuPercent int32 `protobuf:"varint,3,opt,name=cpu_percent,json=cpuPercent,proto3" json:"cpu_percent,omitempty"`
	// Maximum number of pods per node.
	MaxPods int32 `protobuf:"varint,4,opt,name=max_pods,json=maxPods,proto3" json:"max_pods,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_plane_scenario_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_scenario_proto_msgTypes[2]
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
	return file_control_plane_scenario_proto_rawDescGZIP(), []int{2}
}

func (x *Node) GetNodeType() string {
	if x != nil {
		return x.NodeType
	}
	return ""
}

func (x *Node) GetRam() string {
	if x != nil {
		return x.Ram
	}
	return ""
}

func (x *Node) GetCpuPercent() int32 {
	if x != nil {
		return x.CpuPercent
	}
	return 0
}

func (x *Node) GetMaxPods() int32 {
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
		mi := &file_control_plane_scenario_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeGroup) ProtoMessage() {}

func (x *NodeGroup) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_scenario_proto_msgTypes[3]
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
	return file_control_plane_scenario_proto_rawDescGZIP(), []int{3}
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

var File_control_plane_scenario_proto protoreflect.FileDescriptor

var file_control_plane_scenario_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f,
	0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xb3, 0x01, 0x0a, 0x0e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x53, 0x63, 0x65, 0x6e, 0x61,
	0x72, 0x69, 0x6f, 0x12, 0x2f, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e,
	0x6f, 0x64, 0x65, 0x73, 0x12, 0x3f, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x61, 0x74,
	0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x2f, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x67, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x22,
	0x71, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x70, 0x75, 0x5f, 0x70, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x63, 0x70, 0x75,
	0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x5f, 0x70,
	0x6f, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x61, 0x78, 0x50, 0x6f,
	0x64, 0x73, 0x22, 0x5f, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x32, 0x9c, 0x01, 0x0a, 0x08, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x12, 0x4d, 0x0a, 0x0c, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x12, 0x23, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x53, 0x63, 0x65,
	0x6e, 0x61, 0x72, 0x69, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x41, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65,
	0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_control_plane_scenario_proto_rawDescOnce sync.Once
	file_control_plane_scenario_proto_rawDescData = file_control_plane_scenario_proto_rawDesc
)

func file_control_plane_scenario_proto_rawDescGZIP() []byte {
	file_control_plane_scenario_proto_rawDescOnce.Do(func() {
		file_control_plane_scenario_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_plane_scenario_proto_rawDescData)
	})
	return file_control_plane_scenario_proto_rawDescData
}

var file_control_plane_scenario_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_control_plane_scenario_proto_goTypes = []interface{}{
	(*PublicScenario)(nil), // 0: apate.control_plane.PublicScenario
	(*Task)(nil),           // 1: apate.control_plane.Task
	(*Node)(nil),           // 2: apate.control_plane.Node
	(*NodeGroup)(nil),      // 3: apate.control_plane.NodeGroup
	(*empty.Empty)(nil),    // 4: google.protobuf.Empty
}
var file_control_plane_scenario_proto_depIdxs = []int32{
	2, // 0: apate.control_plane.PublicScenario.nodes:type_name -> apate.control_plane.Node
	3, // 1: apate.control_plane.PublicScenario.node_groups:type_name -> apate.control_plane.NodeGroup
	1, // 2: apate.control_plane.PublicScenario.tasks:type_name -> apate.control_plane.Task
	0, // 3: apate.control_plane.Scenario.loadScenario:input_type -> apate.control_plane.PublicScenario
	4, // 4: apate.control_plane.Scenario.startScenario:input_type -> google.protobuf.Empty
	4, // 5: apate.control_plane.Scenario.loadScenario:output_type -> google.protobuf.Empty
	4, // 6: apate.control_plane.Scenario.startScenario:output_type -> google.protobuf.Empty
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_control_plane_scenario_proto_init() }
func file_control_plane_scenario_proto_init() {
	if File_control_plane_scenario_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_control_plane_scenario_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_control_plane_scenario_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_control_plane_scenario_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_control_plane_scenario_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_control_plane_scenario_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_plane_scenario_proto_goTypes,
		DependencyIndexes: file_control_plane_scenario_proto_depIdxs,
		MessageInfos:      file_control_plane_scenario_proto_msgTypes,
	}.Build()
	File_control_plane_scenario_proto = out.File
	file_control_plane_scenario_proto_rawDesc = nil
	file_control_plane_scenario_proto_goTypes = nil
	file_control_plane_scenario_proto_depIdxs = nil
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
	err := c.cc.Invoke(ctx, "/apate.control_plane.Scenario/loadScenario", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scenarioClient) StartScenario(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.control_plane.Scenario/startScenario", in, out, opts...)
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
		FullMethod: "/apate.control_plane.Scenario/LoadScenario",
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
		FullMethod: "/apate.control_plane.Scenario/StartScenario",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScenarioServer).StartScenario(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Scenario_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.control_plane.Scenario",
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
	Metadata: "control_plane/scenario.proto",
}
