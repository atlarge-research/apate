// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.11.4
// source: apatelet/apatelet.proto

package apatelet

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

var File_apatelet_apatelet_proto protoreflect.FileDescriptor

var file_apatelet_apatelet_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x61, 0x74, 0x65,
	0x6c, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x70, 0x61, 0x74, 0x65,
	0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x4c, 0x0a, 0x08, 0x41, 0x70, 0x61, 0x74, 0x65, 0x6c,
	0x65, 0x74, 0x12, 0x40, 0x0a, 0x0c, 0x73, 0x74, 0x6f, 0x70, 0x41, 0x70, 0x61, 0x74, 0x65, 0x6c,
	0x65, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x2f, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70,
	0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_apatelet_apatelet_proto_goTypes = []interface{}{
	(*empty.Empty)(nil), // 0: google.protobuf.Empty
}
var file_apatelet_apatelet_proto_depIdxs = []int32{
	0, // 0: apate.apatelet.Apatelet.stopApatelet:input_type -> google.protobuf.Empty
	0, // 1: apate.apatelet.Apatelet.stopApatelet:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apatelet_apatelet_proto_init() }
func file_apatelet_apatelet_proto_init() {
	if File_apatelet_apatelet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_apatelet_apatelet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apatelet_apatelet_proto_goTypes,
		DependencyIndexes: file_apatelet_apatelet_proto_depIdxs,
	}.Build()
	File_apatelet_apatelet_proto = out.File
	file_apatelet_apatelet_proto_rawDesc = nil
	file_apatelet_apatelet_proto_goTypes = nil
	file_apatelet_apatelet_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ApateletClient is the client API for Apatelet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApateletClient interface {
	// This will signal that the apatelet should leave the cluster and stop
	StopApatelet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
}

type apateletClient struct {
	cc grpc.ClientConnInterface
}

func NewApateletClient(cc grpc.ClientConnInterface) ApateletClient {
	return &apateletClient{cc}
}

func (c *apateletClient) StopApatelet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.apatelet.Apatelet/stopApatelet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApateletServer is the server API for Apatelet service.
type ApateletServer interface {
	// This will signal that the apatelet should leave the cluster and stop
	StopApatelet(context.Context, *empty.Empty) (*empty.Empty, error)
}

// UnimplementedApateletServer can be embedded to have forward compatible implementations.
type UnimplementedApateletServer struct {
}

func (*UnimplementedApateletServer) StopApatelet(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopApatelet not implemented")
}

func RegisterApateletServer(s *grpc.Server, srv ApateletServer) {
	s.RegisterService(&_Apatelet_serviceDesc, srv)
}

func _Apatelet_StopApatelet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApateletServer).StopApatelet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.apatelet.Apatelet/StopApatelet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApateletServer).StopApatelet(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Apatelet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.apatelet.Apatelet",
	HandlerType: (*ApateletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "stopApatelet",
			Handler:    _Apatelet_StopApatelet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apatelet/apatelet.proto",
}
