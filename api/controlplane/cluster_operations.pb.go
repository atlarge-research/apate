// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: controlplane/cluster_operations.proto

package controlplane

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

type ApateletInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The port of the apatelet api
	Port int32 `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *ApateletInformation) Reset() {
	*x = ApateletInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_cluster_operations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApateletInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApateletInformation) ProtoMessage() {}

func (x *ApateletInformation) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_cluster_operations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApateletInformation.ProtoReflect.Descriptor instead.
func (*ApateletInformation) Descriptor() ([]byte, []int) {
	return file_controlplane_cluster_operations_proto_rawDescGZIP(), []int{0}
}

func (x *ApateletInformation) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type NodeHardware struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The amount of bytes of memory
	Memory int64 `protobuf:"varint,1,opt,name=memory,proto3" json:"memory,omitempty"`
	// The amount of milli CPUs in Kubernetes
	Cpu int64 `protobuf:"varint,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
	// The amount of bytes of storage
	Storage int64 `protobuf:"varint,3,opt,name=storage,proto3" json:"storage,omitempty"`
	// The amount of bytes of ephemeral storage
	EphemeralStorage int64 `protobuf:"varint,4,opt,name=ephemeral_storage,json=ephemeralStorage,proto3" json:"ephemeral_storage,omitempty"`
	// The max amount of pods in Kubernetes
	MaxPods int64 `protobuf:"varint,5,opt,name=max_pods,json=maxPods,proto3" json:"max_pods,omitempty"`
}

func (x *NodeHardware) Reset() {
	*x = NodeHardware{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_cluster_operations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeHardware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeHardware) ProtoMessage() {}

func (x *NodeHardware) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_cluster_operations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeHardware.ProtoReflect.Descriptor instead.
func (*NodeHardware) Descriptor() ([]byte, []int) {
	return file_controlplane_cluster_operations_proto_rawDescGZIP(), []int{1}
}

func (x *NodeHardware) GetMemory() int64 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *NodeHardware) GetCpu() int64 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *NodeHardware) GetStorage() int64 {
	if x != nil {
		return x.Storage
	}
	return 0
}

func (x *NodeHardware) GetEphemeralStorage() int64 {
	if x != nil {
		return x.EphemeralStorage
	}
	return 0
}

func (x *NodeHardware) GetMaxPods() int64 {
	if x != nil {
		return x.MaxPods
	}
	return 0
}

type JoinInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The kube config which can be used to join the kubernetes cluster
	KubeConfig []byte `protobuf:"bytes,1,opt,name=kube_config,json=kubeConfig,proto3" json:"kube_config,omitempty"`
	// The UUID that will be used from the control plane to identify this node
	NodeUuid string `protobuf:"bytes,2,opt,name=node_uuid,json=nodeUuid,proto3" json:"node_uuid,omitempty"`
	// The hardware that was 'allocated' to the apatelet
	Hardware *NodeHardware `protobuf:"bytes,3,opt,name=hardware,proto3" json:"hardware,omitempty"`
}

func (x *JoinInformation) Reset() {
	*x = JoinInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_cluster_operations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinInformation) ProtoMessage() {}

func (x *JoinInformation) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_cluster_operations_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinInformation.ProtoReflect.Descriptor instead.
func (*JoinInformation) Descriptor() ([]byte, []int) {
	return file_controlplane_cluster_operations_proto_rawDescGZIP(), []int{2}
}

func (x *JoinInformation) GetKubeConfig() []byte {
	if x != nil {
		return x.KubeConfig
	}
	return nil
}

func (x *JoinInformation) GetNodeUuid() string {
	if x != nil {
		return x.NodeUuid
	}
	return ""
}

func (x *JoinInformation) GetHardware() *NodeHardware {
	if x != nil {
		return x.Hardware
	}
	return nil
}

type LeaveInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The UUID that is be used from the control plane to identify this node
	NodeUuid string `protobuf:"bytes,1,opt,name=node_uuid,json=nodeUuid,proto3" json:"node_uuid,omitempty"`
}

func (x *LeaveInformation) Reset() {
	*x = LeaveInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_cluster_operations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveInformation) ProtoMessage() {}

func (x *LeaveInformation) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_cluster_operations_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveInformation.ProtoReflect.Descriptor instead.
func (*LeaveInformation) Descriptor() ([]byte, []int) {
	return file_controlplane_cluster_operations_proto_rawDescGZIP(), []int{3}
}

func (x *LeaveInformation) GetNodeUuid() string {
	if x != nil {
		return x.NodeUuid
	}
	return ""
}

var File_controlplane_cluster_operations_proto protoreflect.FileDescriptor

var file_controlplane_cluster_operations_proto_rawDesc = []byte{
	0x0a, 0x25, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x29, 0x0a, 0x13, 0x41, 0x70, 0x61, 0x74,
	0x65, 0x6c, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x22, 0x9a, 0x01, 0x0a, 0x0c, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x61, 0x72, 0x64,
	0x77, 0x61, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x63, 0x70, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x63, 0x70, 0x75, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x70, 0x68, 0x65,
	0x6d, 0x65, 0x72, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x10, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x6f, 0x64,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d, 0x61, 0x78, 0x50, 0x6f, 0x64, 0x73,
	0x22, 0x8d, 0x01, 0x0a, 0x0f, 0x4a, 0x6f, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x75,
	0x69, 0x64, 0x12, 0x3c, 0x0a, 0x08, 0x68, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x61,
	0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x52, 0x08, 0x68, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65,
	0x22, 0x2f, 0x0a, 0x10, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x75, 0x69,
	0x64, 0x32, 0xc2, 0x01, 0x0a, 0x11, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x5d, 0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x41, 0x70, 0x61, 0x74,
	0x65, 0x6c, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x23, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x24, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4c, 0x65, 0x61, 0x76,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x65, 0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controlplane_cluster_operations_proto_rawDescOnce sync.Once
	file_controlplane_cluster_operations_proto_rawDescData = file_controlplane_cluster_operations_proto_rawDesc
)

func file_controlplane_cluster_operations_proto_rawDescGZIP() []byte {
	file_controlplane_cluster_operations_proto_rawDescOnce.Do(func() {
		file_controlplane_cluster_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_controlplane_cluster_operations_proto_rawDescData)
	})
	return file_controlplane_cluster_operations_proto_rawDescData
}

var file_controlplane_cluster_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_controlplane_cluster_operations_proto_goTypes = []interface{}{
	(*ApateletInformation)(nil), // 0: apate.controlplane.ApateletInformation
	(*NodeHardware)(nil),        // 1: apate.controlplane.NodeHardware
	(*JoinInformation)(nil),     // 2: apate.controlplane.JoinInformation
	(*LeaveInformation)(nil),    // 3: apate.controlplane.LeaveInformation
	(*empty.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_controlplane_cluster_operations_proto_depIdxs = []int32{
	1, // 0: apate.controlplane.JoinInformation.hardware:type_name -> apate.controlplane.NodeHardware
	0, // 1: apate.controlplane.ClusterOperations.joinCluster:input_type -> apate.controlplane.ApateletInformation
	3, // 2: apate.controlplane.ClusterOperations.leaveCluster:input_type -> apate.controlplane.LeaveInformation
	2, // 3: apate.controlplane.ClusterOperations.joinCluster:output_type -> apate.controlplane.JoinInformation
	4, // 4: apate.controlplane.ClusterOperations.leaveCluster:output_type -> google.protobuf.Empty
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_controlplane_cluster_operations_proto_init() }
func file_controlplane_cluster_operations_proto_init() {
	if File_controlplane_cluster_operations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controlplane_cluster_operations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApateletInformation); i {
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
		file_controlplane_cluster_operations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeHardware); i {
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
		file_controlplane_cluster_operations_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinInformation); i {
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
		file_controlplane_cluster_operations_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveInformation); i {
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
			RawDescriptor: file_controlplane_cluster_operations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_controlplane_cluster_operations_proto_goTypes,
		DependencyIndexes: file_controlplane_cluster_operations_proto_depIdxs,
		MessageInfos:      file_controlplane_cluster_operations_proto_msgTypes,
	}.Build()
	File_controlplane_cluster_operations_proto = out.File
	file_controlplane_cluster_operations_proto_rawDesc = nil
	file_controlplane_cluster_operations_proto_goTypes = nil
	file_controlplane_cluster_operations_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ClusterOperationsClient is the client API for ClusterOperations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterOperationsClient interface {
	// Joins the node to the Apate cluster.
	// Will return information needed to identify yourself
	// And to join the kubernetes cluster.
	JoinCluster(ctx context.Context, in *ApateletInformation, opts ...grpc.CallOption) (*JoinInformation, error)
	// Removes the node from the cluster
	LeaveCluster(ctx context.Context, in *LeaveInformation, opts ...grpc.CallOption) (*empty.Empty, error)
}

type clusterOperationsClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterOperationsClient(cc grpc.ClientConnInterface) ClusterOperationsClient {
	return &clusterOperationsClient{cc}
}

func (c *clusterOperationsClient) JoinCluster(ctx context.Context, in *ApateletInformation, opts ...grpc.CallOption) (*JoinInformation, error) {
	out := new(JoinInformation)
	err := c.cc.Invoke(ctx, "/apate.controlplane.ClusterOperations/joinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterOperationsClient) LeaveCluster(ctx context.Context, in *LeaveInformation, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.controlplane.ClusterOperations/leaveCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterOperationsServer is the server API for ClusterOperations service.
type ClusterOperationsServer interface {
	// Joins the node to the Apate cluster.
	// Will return information needed to identify yourself
	// And to join the kubernetes cluster.
	JoinCluster(context.Context, *ApateletInformation) (*JoinInformation, error)
	// Removes the node from the cluster
	LeaveCluster(context.Context, *LeaveInformation) (*empty.Empty, error)
}

// UnimplementedClusterOperationsServer can be embedded to have forward compatible implementations.
type UnimplementedClusterOperationsServer struct {
}

func (*UnimplementedClusterOperationsServer) JoinCluster(context.Context, *ApateletInformation) (*JoinInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (*UnimplementedClusterOperationsServer) LeaveCluster(context.Context, *LeaveInformation) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCluster not implemented")
}

func RegisterClusterOperationsServer(s *grpc.Server, srv ClusterOperationsServer) {
	s.RegisterService(&_ClusterOperations_serviceDesc, srv)
}

func _ClusterOperations_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApateletInformation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterOperationsServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.controlplane.ClusterOperations/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterOperationsServer).JoinCluster(ctx, req.(*ApateletInformation))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterOperations_LeaveCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveInformation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterOperationsServer).LeaveCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.controlplane.ClusterOperations/LeaveCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterOperationsServer).LeaveCluster(ctx, req.(*LeaveInformation))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterOperations_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.controlplane.ClusterOperations",
	HandlerType: (*ClusterOperationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "joinCluster",
			Handler:    _ClusterOperations_JoinCluster_Handler,
		},
		{
			MethodName: "leaveCluster",
			Handler:    _ClusterOperations_LeaveCluster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controlplane/cluster_operations.proto",
}
