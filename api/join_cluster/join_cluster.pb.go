// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: join_cluster/join_cluster.proto

package join_cluster

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
		mi := &file_join_cluster_join_cluster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinInformation) ProtoMessage() {}

func (x *JoinInformation) ProtoReflect() protoreflect.Message {
	mi := &file_join_cluster_join_cluster_proto_msgTypes[0]
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
	return file_join_cluster_join_cluster_proto_rawDescGZIP(), []int{0}
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

var File_join_cluster_join_cluster_proto protoreflect.FileDescriptor

var file_join_cluster_join_cluster_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x6a,
	0x6f, 0x69, 0x6e, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0f, 0x4a, 0x6f, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6b, 0x75, 0x62, 0x65,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55,
	0x55, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x55,
	0x55, 0x49, 0x44, 0x32, 0x5b, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x4c, 0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x23, 0x2e, 0x61, 0x70, 0x61, 0x74,
	0x65, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4a,
	0x6f, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00,
	0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f,
	0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x6b,
	0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f,
	0x69, 0x6e, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_join_cluster_join_cluster_proto_rawDescOnce sync.Once
	file_join_cluster_join_cluster_proto_rawDescData = file_join_cluster_join_cluster_proto_rawDesc
)

func file_join_cluster_join_cluster_proto_rawDescGZIP() []byte {
	file_join_cluster_join_cluster_proto_rawDescOnce.Do(func() {
		file_join_cluster_join_cluster_proto_rawDescData = protoimpl.X.CompressGZIP(file_join_cluster_join_cluster_proto_rawDescData)
	})
	return file_join_cluster_join_cluster_proto_rawDescData
}

var file_join_cluster_join_cluster_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_join_cluster_join_cluster_proto_goTypes = []interface{}{
	(*JoinInformation)(nil), // 0: apate.join_cluster.JoinInformation
	(*empty.Empty)(nil),     // 1: google.protobuf.Empty
}
var file_join_cluster_join_cluster_proto_depIdxs = []int32{
	1, // 0: apate.join_cluster.JoinCluster.joinCluster:input_type -> google.protobuf.Empty
	0, // 1: apate.join_cluster.JoinCluster.joinCluster:output_type -> apate.join_cluster.JoinInformation
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_join_cluster_join_cluster_proto_init() }
func file_join_cluster_join_cluster_proto_init() {
	if File_join_cluster_join_cluster_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_join_cluster_join_cluster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_join_cluster_join_cluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_join_cluster_join_cluster_proto_goTypes,
		DependencyIndexes: file_join_cluster_join_cluster_proto_depIdxs,
		MessageInfos:      file_join_cluster_join_cluster_proto_msgTypes,
	}.Build()
	File_join_cluster_join_cluster_proto = out.File
	file_join_cluster_join_cluster_proto_rawDesc = nil
	file_join_cluster_join_cluster_proto_goTypes = nil
	file_join_cluster_join_cluster_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// JoinClusterClient is the client API for JoinCluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JoinClusterClient interface {
	// Joins the node to the Apate cluster.
	// Will return information needed to identify yourself
	// And to join the kubernetes cluster.
	JoinCluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinInformation, error)
}

type joinClusterClient struct {
	cc grpc.ClientConnInterface
}

func NewJoinClusterClient(cc grpc.ClientConnInterface) JoinClusterClient {
	return &joinClusterClient{cc}
}

func (c *joinClusterClient) JoinCluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinInformation, error) {
	out := new(JoinInformation)
	err := c.cc.Invoke(ctx, "/apate.join_cluster.JoinCluster/joinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JoinClusterServer is the server API for JoinCluster service.
type JoinClusterServer interface {
	// Joins the node to the Apate cluster.
	// Will return information needed to identify yourself
	// And to join the kubernetes cluster.
	JoinCluster(context.Context, *empty.Empty) (*JoinInformation, error)
}

// UnimplementedJoinClusterServer can be embedded to have forward compatible implementations.
type UnimplementedJoinClusterServer struct {
}

func (*UnimplementedJoinClusterServer) JoinCluster(context.Context, *empty.Empty) (*JoinInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}

func RegisterJoinClusterServer(s *grpc.Server, srv JoinClusterServer) {
	s.RegisterService(&_JoinCluster_serviceDesc, srv)
}

func _JoinCluster_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JoinClusterServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.join_cluster.JoinCluster/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JoinClusterServer).JoinCluster(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _JoinCluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.join_cluster.JoinCluster",
	HandlerType: (*JoinClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "joinCluster",
			Handler:    _JoinCluster_JoinCluster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "join_cluster/join_cluster.proto",
}
