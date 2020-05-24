// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.11.4
// source: apatelet/scenario.proto

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

// The top level object which defines how the different Apatelet will emulate certain deployments
type ApateletScenario struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The absolute timestamp at which the scenario will start
	StartTime int64 `protobuf:"varint,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
}

func (x *ApateletScenario) Reset() {
	*x = ApateletScenario{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apatelet_scenario_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApateletScenario) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApateletScenario) ProtoMessage() {}

func (x *ApateletScenario) ProtoReflect() protoreflect.Message {
	mi := &file_apatelet_scenario_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApateletScenario.ProtoReflect.Descriptor instead.
func (*ApateletScenario) Descriptor() ([]byte, []int) {
	return file_apatelet_scenario_proto_rawDescGZIP(), []int{0}
}

func (x *ApateletScenario) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

var File_apatelet_scenario_proto protoreflect.FileDescriptor

var file_apatelet_scenario_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x2f, 0x73, 0x63, 0x65, 0x6e, 0x61,
	0x72, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x70, 0x61, 0x74, 0x65,
	0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x31, 0x0a, 0x10, 0x41, 0x70, 0x61, 0x74, 0x65, 0x6c,
	0x65, 0x74, 0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x32, 0x57, 0x0a, 0x08, 0x53, 0x63, 0x65,
	0x6e, 0x61, 0x72, 0x69, 0x6f, 0x12, 0x4b, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63,
	0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x61,
	0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x2e, 0x41, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74,
	0x53, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65,
	0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x70, 0x61, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apatelet_scenario_proto_rawDescOnce sync.Once
	file_apatelet_scenario_proto_rawDescData = file_apatelet_scenario_proto_rawDesc
)

func file_apatelet_scenario_proto_rawDescGZIP() []byte {
	file_apatelet_scenario_proto_rawDescOnce.Do(func() {
		file_apatelet_scenario_proto_rawDescData = protoimpl.X.CompressGZIP(file_apatelet_scenario_proto_rawDescData)
	})
	return file_apatelet_scenario_proto_rawDescData
}

var file_apatelet_scenario_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_apatelet_scenario_proto_goTypes = []interface{}{
	(*ApateletScenario)(nil), // 0: apate.apatelet.ApateletScenario
	(*empty.Empty)(nil),      // 1: google.protobuf.Empty
}
var file_apatelet_scenario_proto_depIdxs = []int32{
	0, // 0: apate.apatelet.Scenario.startScenario:input_type -> apate.apatelet.ApateletScenario
	1, // 1: apate.apatelet.Scenario.startScenario:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apatelet_scenario_proto_init() }
func file_apatelet_scenario_proto_init() {
	if File_apatelet_scenario_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apatelet_scenario_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApateletScenario); i {
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
			RawDescriptor: file_apatelet_scenario_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apatelet_scenario_proto_goTypes,
		DependencyIndexes: file_apatelet_scenario_proto_depIdxs,
		MessageInfos:      file_apatelet_scenario_proto_msgTypes,
	}.Build()
	File_apatelet_scenario_proto = out.File
	file_apatelet_scenario_proto_rawDesc = nil
	file_apatelet_scenario_proto_goTypes = nil
	file_apatelet_scenario_proto_depIdxs = nil
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
	// Starts a scenario on the current Apatelet
	// This will be called on every Apatelet
	StartScenario(ctx context.Context, in *ApateletScenario, opts ...grpc.CallOption) (*empty.Empty, error)
}

type scenarioClient struct {
	cc grpc.ClientConnInterface
}

func NewScenarioClient(cc grpc.ClientConnInterface) ScenarioClient {
	return &scenarioClient{cc}
}

func (c *scenarioClient) StartScenario(ctx context.Context, in *ApateletScenario, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/apate.apatelet.Scenario/startScenario", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScenarioServer is the server API for Scenario service.
type ScenarioServer interface {
	// Starts a scenario on the current Apatelet
	// This will be called on every Apatelet
	StartScenario(context.Context, *ApateletScenario) (*empty.Empty, error)
}

// UnimplementedScenarioServer can be embedded to have forward compatible implementations.
type UnimplementedScenarioServer struct {
}

func (*UnimplementedScenarioServer) StartScenario(context.Context, *ApateletScenario) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartScenario not implemented")
}

func RegisterScenarioServer(s *grpc.Server, srv ScenarioServer) {
	s.RegisterService(&_Scenario_serviceDesc, srv)
}

func _Scenario_StartScenario_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApateletScenario)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScenarioServer).StartScenario(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apate.apatelet.Scenario/StartScenario",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScenarioServer).StartScenario(ctx, req.(*ApateletScenario))
	}
	return interceptor(ctx, in, info, handler)
}

var _Scenario_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apate.apatelet.Scenario",
	HandlerType: (*ScenarioServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "startScenario",
			Handler:    _Scenario_StartScenario_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apatelet/scenario.proto",
}
