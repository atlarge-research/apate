// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: control_plane/cluster_operations.proto

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

type JoinInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The kube config which can be used to join the kubernetes cluster
	KubeConfig []byte `protobuf:"bytes,1,opt,name=kubeConfig,proto3" json:"kubeConfig,omitempty"`
	// The context used for joining the cluster
	KubeContext string `protobuf:"bytes,2,opt,name=kubeContext,proto3" json:"kubeContext,omitempty"`
	// The UUID that will be used from the control plane to identify this node
	NodeUUID string `protobuf:"bytes,3,opt,name=nodeUUID,proto3" json:"nodeUUID,omitempty"`
}

func (x *JoinInformation) Reset() {
	*x = JoinInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_plane_cluster_operations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinInformation) ProtoMessage() {}

func (x *JoinInformation) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_cluster_operations_proto_msgTypes[0]
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
	return file_control_plane_cluster_operations_proto_rawDescGZIP(), []int{0}
}

func (x *JoinInformation) GetKubeConfig() []byte {
	if x != nil {
		return x.KubeConfig
	}
	return nil
}

func (x *JoinInformation) GetKubeContext() string {
	if x != nil {
		return x.KubeContext
	}
	return ""
}

func (x *JoinInformation) GetNodeUUID() string {
	if x != nil {
		return x.NodeUUID
	}
	return ""
}

type LeaveInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The UUID that is be used from the control plane to identify this node
	NodeUUID string `protobuf:"bytes,1,opt,name=nodeUUID,proto3" json:"nodeUUID,omitempty"`
}

func (x *LeaveInformation) Reset() {
	*x = LeaveInformation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_plane_cluster_operations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveInformation) ProtoMessage() {}

func (x *LeaveInformation) ProtoReflect() protoreflect.Message {
	mi := &file_control_plane_cluster_operations_proto_msgTypes[1]
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
	return file_control_plane_cluster_operations_proto_rawDescGZIP(), []int{1}
}

func (x *LeaveInformation) GetNodeUUID() string {
	if x != nil {
		return x.NodeUUID
	}
	return ""
}

var File_control_plane_cluster_operations_proto protoreflect.FileDescriptor

var file_control_plane_cluster_operations_proto_rawDesc = []byte{
	0x0a, 0x26, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0f, 0x4a, 0x6f,
	0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a,
	0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x20, 0x0a,
	0x0b, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x55, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x55, 0x49, 0x44, 0x22, 0x2e, 0x0a, 0x10, 0x4c,
	0x65, 0x61, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55, 0x55, 0x49, 0x44, 0x32, 0xb3, 0x01, 0x0a, 0x11,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x4d, 0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x24, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4a,
	0x6f, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00,
	0x12, 0x4f, 0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x25, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x2d,
	0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_control_plane_cluster_operations_proto_rawDescOnce sync.Once
	file_control_plane_cluster_operations_proto_rawDescData = file_control_plane_cluster_operations_proto_rawDesc
)

func file_control_plane_cluster_operations_proto_rawDescGZIP() []byte {
	file_control_plane_cluster_operations_proto_rawDescOnce.Do(func() {
		file_control_plane_cluster_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_plane_cluster_operations_proto_rawDescData)
	})
	return file_control_plane_cluster_operations_proto_rawDescData
}

var file_control_plane_cluster_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_control_plane_cluster_operations_proto_goTypes = []interface{}{
	(*JoinInformation)(nil),  // 0: apate.control_plane.JoinInformation
	(*LeaveInformation)(nil), // 1: apate.control_plane.LeaveInformation
	(*empty.Empty)(nil),      // 2: google.protobuf.Empty
}
var file_control_plane_cluster_operations_proto_depIdxs = []int32{
	2, // 0: apate.control_plane.ClusterOperations.joinCluster:input_type -> google.protobuf.Empty
	1, // 1: apate.control_plane.ClusterOperations.leaveCluster:input_type -> apate.control_plane.LeaveInformation
	0, // 2: apate.control_plane.ClusterOperations.joinCluster:output_type -> apate.control_plane.JoinInformation
	2, // 3: apate.control_plane.ClusterOperations.leaveCluster:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_control_plane_cluster_operations_proto_init() }
func file_control_plane_cluster_operations_proto_init() {
	if File_control_plane_cluster_operations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_control_plane_cluster_operations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_control_plane_cluster_operations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_control_plane_cluster_operations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_plane_cluster_operations_proto_goTypes,
		DependencyIndexes: file_control_plane_cluster_operations_proto_depIdxs,
		MessageInfos:      file_control_plane_cluster_operations_proto_msgTypes,
	}.Build()
	File_control_plane_cluster_operations_proto = out.File
	file_control_plane_cluster_operations_proto_rawDesc = nil
	file_control_plane_cluster_operations_proto_goTypes = nil
	file_control_plane_cluster_operations_proto_depIdxs = nil
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
	JoinCluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinInformation, error)
	// Removes the node from the cluster
	LeaveCluster(ctx context.Context, in *LeaveInformation, opts ...grpc.CallOption) (*empty.Empty, error)
}

type clusterOperationsClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterOperationsClient(cc grpc.ClientConnInterface) ClusterOperationsClient {
	return &clusterOperationsClient{cc}
}

func (c *clusterOperationsClient) JoinCluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinInformation, error) {
	out := new(JoinInformation)
	err := c.cc.Invoke(ctx, "/apate.control_plane.ClusterOperations/joinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterOperationsClient) LeaveCluster(ctx context.Context, in *LeaveInformation, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.control_plane.ClusterOperations/leaveCluster", in, out, opts...)
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
	JoinCluster(context.Context, *empty.Empty) (*JoinInformation, error)
	// Removes the node from the cluster
	LeaveCluster(context.Context, *LeaveInformation) (*empty.Empty, error)
}

// UnimplementedClusterOperationsServer can be embedded to have forward compatible implementations.
type UnimplementedClusterOperationsServer struct {
}

func (*UnimplementedClusterOperationsServer) JoinCluster(context.Context, *empty.Empty) (*JoinInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (*UnimplementedClusterOperationsServer) LeaveCluster(context.Context, *LeaveInformation) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCluster not implemented")
}

func RegisterClusterOperationsServer(s *grpc.Server, srv ClusterOperationsServer) {
	s.RegisterService(&_ClusterOperations_serviceDesc, srv)
}

func _ClusterOperations_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterOperationsServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.control_plane.ClusterOperations/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterOperationsServer).JoinCluster(ctx, req.(*empty.Empty))
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
		FullMethod: "/apate.control_plane.ClusterOperations/LeaveCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterOperationsServer).LeaveCluster(ctx, req.(*LeaveInformation))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterOperations_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.control_plane.ClusterOperations",
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
	Metadata: "control_plane/cluster_operations.proto",
}
