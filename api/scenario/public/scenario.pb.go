// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: scenario/public/scenario.proto

package public

import (
	context "context"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type SendScenarioResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendScenarioResponse) Reset() {
	*x = SendScenarioResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scenario_public_scenario_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendScenarioResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendScenarioResponse) ProtoMessage() {}

func (x *SendScenarioResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scenario_public_scenario_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendScenarioResponse.ProtoReflect.Descriptor instead.
func (*SendScenarioResponse) Descriptor() ([]byte, []int) {
	return file_scenario_public_scenario_proto_rawDescGZIP(), []int{0}
}

// Scenario represents a full scenario as given in the scenario configuration file.
// It can be easily constructed from yaml or json (pkg/scenario/deserialize) and has to be
// converted to a private scenario (pkg/scenario/normalize).
type Scenario struct {
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

func (x *Scenario) Reset() {
	*x = Scenario{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scenario_public_scenario_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Scenario) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Scenario) ProtoMessage() {}

func (x *Scenario) ProtoReflect() protoreflect.Message {
	mi := &file_scenario_public_scenario_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Scenario.ProtoReflect.Descriptor instead.
func (*Scenario) Descriptor() ([]byte, []int) {
	return file_scenario_public_scenario_proto_rawDescGZIP(), []int{1}
}

func (x *Scenario) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Scenario) GetNodeGroups() []*NodeGroup {
	if x != nil {
		return x.NodeGroups
	}
	return nil
}

func (x *Scenario) GetTasks() []*Task {
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
		mi := &file_scenario_public_scenario_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_scenario_public_scenario_proto_msgTypes[2]
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
	return file_scenario_public_scenario_proto_rawDescGZIP(), []int{2}
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
		mi := &file_scenario_public_scenario_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_scenario_public_scenario_proto_msgTypes[3]
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
	return file_scenario_public_scenario_proto_rawDescGZIP(), []int{3}
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
	NodeType string `protobuf:"bytes,3,opt,name=node_type,json=nodeType,proto3" json:"node_type,omitempty"`
	// How many times you want this type of node.
	Amount int32 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *NodeGroup) Reset() {
	*x = NodeGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scenario_public_scenario_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeGroup) ProtoMessage() {}

func (x *NodeGroup) ProtoReflect() protoreflect.Message {
	mi := &file_scenario_public_scenario_proto_msgTypes[4]
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
	return file_scenario_public_scenario_proto_rawDescGZIP(), []int{4}
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

var File_scenario_public_scenario_proto protoreflect.FileDescriptor

var file_scenario_public_scenario_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x2f, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x16, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x71, 0x0a, 0x08, 0x53, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x12, 0x1b, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x2b, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x1b,
	0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x67, 0x0a, 0x04, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x76, 0x65, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x76,
	0x65, 0x72, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x22, 0x71, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x61, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x70, 0x75, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x63, 0x70, 0x75, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x6d, 0x61, 0x78, 0x5f, 0x70, 0x6f, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x6d, 0x61, 0x78, 0x50, 0x6f, 0x64, 0x73, 0x22, 0x5f, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x44, 0x0a, 0x0e, 0x53, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x0c, 0x73, 0x65,
	0x6e, 0x64, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x12, 0x09, 0x2e, 0x53, 0x63, 0x65,
	0x6e, 0x61, 0x72, 0x69, 0x6f, 0x1a, 0x15, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4b,
	0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x74, 0x6c,
	0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x6f, 0x70,
	0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x6b, 0x75, 0x62,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_scenario_public_scenario_proto_rawDescOnce sync.Once
	file_scenario_public_scenario_proto_rawDescData = file_scenario_public_scenario_proto_rawDesc
)

func file_scenario_public_scenario_proto_rawDescGZIP() []byte {
	file_scenario_public_scenario_proto_rawDescOnce.Do(func() {
		file_scenario_public_scenario_proto_rawDescData = protoimpl.X.CompressGZIP(file_scenario_public_scenario_proto_rawDescData)
	})
	return file_scenario_public_scenario_proto_rawDescData
}

var file_scenario_public_scenario_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_scenario_public_scenario_proto_goTypes = []interface{}{
	(*SendScenarioResponse)(nil), // 0: SendScenarioResponse
	(*Scenario)(nil),             // 1: Scenario
	(*Task)(nil),                 // 2: Task
	(*Node)(nil),                 // 3: Node
	(*NodeGroup)(nil),            // 4: NodeGroup
}
var file_scenario_public_scenario_proto_depIdxs = []int32{
	3, // 0: Scenario.nodes:type_name -> Node
	4, // 1: Scenario.node_groups:type_name -> NodeGroup
	2, // 2: Scenario.tasks:type_name -> Task
	1, // 3: ScenarioSender.sendScenario:input_type -> Scenario
	0, // 4: ScenarioSender.sendScenario:output_type -> SendScenarioResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_scenario_public_scenario_proto_init() }
func file_scenario_public_scenario_proto_init() {
	if File_scenario_public_scenario_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_scenario_public_scenario_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendScenarioResponse); i {
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
		file_scenario_public_scenario_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Scenario); i {
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
		file_scenario_public_scenario_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_scenario_public_scenario_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_scenario_public_scenario_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_scenario_public_scenario_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_scenario_public_scenario_proto_goTypes,
		DependencyIndexes: file_scenario_public_scenario_proto_depIdxs,
		MessageInfos:      file_scenario_public_scenario_proto_msgTypes,
	}.Build()
	File_scenario_public_scenario_proto = out.File
	file_scenario_public_scenario_proto_rawDesc = nil
	file_scenario_public_scenario_proto_goTypes = nil
	file_scenario_public_scenario_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ScenarioSenderClient is the client API for ScenarioSender service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ScenarioSenderClient interface {
	SendScenario(ctx context.Context, in *Scenario, opts ...grpc.CallOption) (*SendScenarioResponse, error)
}

type scenarioSenderClient struct {
	cc grpc.ClientConnInterface
}

func NewScenarioSenderClient(cc grpc.ClientConnInterface) ScenarioSenderClient {
	return &scenarioSenderClient{cc}
}

func (c *scenarioSenderClient) SendScenario(ctx context.Context, in *Scenario, opts ...grpc.CallOption) (*SendScenarioResponse, error) {
	out := new(SendScenarioResponse)
	err := c.cc.Invoke(ctx, "/ScenarioSender/sendScenario", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScenarioSenderServer is the server API for ScenarioSender service.
type ScenarioSenderServer interface {
	SendScenario(context.Context, *Scenario) (*SendScenarioResponse, error)
}

// UnimplementedScenarioSenderServer can be embedded to have forward compatible implementations.
type UnimplementedScenarioSenderServer struct {
}

func (*UnimplementedScenarioSenderServer) SendScenario(context.Context, *Scenario) (*SendScenarioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendScenario not implemented")
}

func RegisterScenarioSenderServer(s *grpc.Server, srv ScenarioSenderServer) {
	s.RegisterService(&_ScenarioSender_serviceDesc, srv)
}

func _ScenarioSender_SendScenario_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Scenario)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScenarioSenderServer).SendScenario(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ScenarioSender/SendScenario",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScenarioSenderServer).SendScenario(ctx, req.(*Scenario))
	}
	return interceptor(ctx, in, info, handler)
}

var _ScenarioSender_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ScenarioSender",
	HandlerType: (*ScenarioSenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendScenario",
			Handler:    _ScenarioSender_SendScenario_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scenario/public/scenario.proto",
}
