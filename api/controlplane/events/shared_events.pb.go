// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: controlplane/events/shared_events.proto

package events

import (
	scenario "github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	proto "github.com/golang/protobuf/proto"
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

// See https://github.com/virtual-kubelet/virtual-kubelet/#podlifecylcehandler
// Ping stands for the heartbeat requests sent by Kubernetes
type RequestType int32

const (
	RequestType_CREATE_POD     RequestType = 0
	RequestType_UPDATE_POD     RequestType = 1
	RequestType_DELETE_POD     RequestType = 2
	RequestType_GET_POD        RequestType = 3
	RequestType_GET_POD_STATUS RequestType = 4
	RequestType_GET_PODS       RequestType = 5
	RequestType_PING           RequestType = 6
)

// Enum value maps for RequestType.
var (
	RequestType_name = map[int32]string{
		0: "CREATE_POD",
		1: "UPDATE_POD",
		2: "DELETE_POD",
		3: "GET_POD",
		4: "GET_POD_STATUS",
		5: "GET_PODS",
		6: "PING",
	}
	RequestType_value = map[string]int32{
		"CREATE_POD":     0,
		"UPDATE_POD":     1,
		"DELETE_POD":     2,
		"GET_POD":        3,
		"GET_POD_STATUS": 4,
		"GET_PODS":       5,
		"PING":           6,
	}
)

func (x RequestType) Enum() *RequestType {
	p := new(RequestType)
	*p = x
	return p
}

func (x RequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_controlplane_events_shared_events_proto_enumTypes[0].Descriptor()
}

func (RequestType) Type() protoreflect.EnumType {
	return &file_controlplane_events_shared_events_proto_enumTypes[0]
}

func (x RequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestType.Descriptor instead.
func (RequestType) EnumDescriptor() ([]byte, []int) {
	return file_controlplane_events_shared_events_proto_rawDescGZIP(), []int{0}
}

// A percentage of lifecycle requests will error
type ResponseState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The request type to act upon
	Type RequestType `protobuf:"varint,1,opt,name=type,proto3,enum=apate.controlplane.events.RequestType" json:"type,omitempty"`
	// How to respond to this request
	Response scenario.Response `protobuf:"varint,2,opt,name=response,proto3,enum=apate.scenario.Response" json:"response,omitempty"`
	// The percentage of requests to handle like this
	Percentage int32 `protobuf:"varint,3,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

func (x *ResponseState) Reset() {
	*x = ResponseState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controlplane_events_shared_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseState) ProtoMessage() {}

func (x *ResponseState) ProtoReflect() protoreflect.Message {
	mi := &file_controlplane_events_shared_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseState.ProtoReflect.Descriptor instead.
func (*ResponseState) Descriptor() ([]byte, []int) {
	return file_controlplane_events_shared_events_proto_rawDescGZIP(), []int{0}
}

func (x *ResponseState) GetType() RequestType {
	if x != nil {
		return x.Type
	}
	return RequestType_CREATE_POD
}

func (x *ResponseState) GetResponse() scenario.Response {
	if x != nil {
		return x.Response
	}
	return scenario.Response_NORMAL
}

func (x *ResponseState) GetPercentage() int32 {
	if x != nil {
		return x.Percentage
	}
	return 0
}

var File_controlplane_events_shared_events_proto protoreflect.FileDescriptor

var file_controlplane_events_shared_events_proto_rawDesc = []byte{
	0x0a, 0x27, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x5f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x61, 0x70, 0x61, 0x74, 0x65,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x1a, 0x15, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x01, 0x0a, 0x0d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3a, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x61, 0x70,
	0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x34, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x61, 0x70,
	0x61, 0x74, 0x65, 0x2e, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x2a,
	0x76, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e,
	0x0a, 0x0a, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x50, 0x4f, 0x44, 0x10, 0x00, 0x12, 0x0e,
	0x0a, 0x0a, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x50, 0x4f, 0x44, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x50, 0x4f, 0x44, 0x10, 0x02, 0x12, 0x0b,
	0x0a, 0x07, 0x47, 0x45, 0x54, 0x5f, 0x50, 0x4f, 0x44, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x47,
	0x45, 0x54, 0x5f, 0x50, 0x4f, 0x44, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x04, 0x12,
	0x0c, 0x0a, 0x08, 0x47, 0x45, 0x54, 0x5f, 0x50, 0x4f, 0x44, 0x53, 0x10, 0x05, 0x12, 0x08, 0x0a,
	0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x42, 0x4f, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d,
	0x75, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controlplane_events_shared_events_proto_rawDescOnce sync.Once
	file_controlplane_events_shared_events_proto_rawDescData = file_controlplane_events_shared_events_proto_rawDesc
)

func file_controlplane_events_shared_events_proto_rawDescGZIP() []byte {
	file_controlplane_events_shared_events_proto_rawDescOnce.Do(func() {
		file_controlplane_events_shared_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_controlplane_events_shared_events_proto_rawDescData)
	})
	return file_controlplane_events_shared_events_proto_rawDescData
}

var file_controlplane_events_shared_events_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_controlplane_events_shared_events_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_controlplane_events_shared_events_proto_goTypes = []interface{}{
	(RequestType)(0),       // 0: apate.controlplane.events.RequestType
	(*ResponseState)(nil),  // 1: apate.controlplane.events.ResponseState
	(scenario.Response)(0), // 2: apate.scenario.Response
}
var file_controlplane_events_shared_events_proto_depIdxs = []int32{
	0, // 0: apate.controlplane.events.ResponseState.type:type_name -> apate.controlplane.events.RequestType
	2, // 1: apate.controlplane.events.ResponseState.response:type_name -> apate.scenario.Response
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_controlplane_events_shared_events_proto_init() }
func file_controlplane_events_shared_events_proto_init() {
	if File_controlplane_events_shared_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controlplane_events_shared_events_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseState); i {
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
			RawDescriptor: file_controlplane_events_shared_events_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_controlplane_events_shared_events_proto_goTypes,
		DependencyIndexes: file_controlplane_events_shared_events_proto_depIdxs,
		EnumInfos:         file_controlplane_events_shared_events_proto_enumTypes,
		MessageInfos:      file_controlplane_events_shared_events_proto_msgTypes,
	}.Build()
	File_controlplane_events_shared_events_proto = out.File
	file_controlplane_events_shared_events_proto_rawDesc = nil
	file_controlplane_events_shared_events_proto_goTypes = nil
	file_controlplane_events_shared_events_proto_depIdxs = nil
}