// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: scenario/events.proto

package scenario

import (
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

type Response int32

const (
	Response_NORMAL  Response = 0
	Response_TIMEOUT Response = 1
	Response_ERROR   Response = 2
)

// Enum value maps for Response.
var (
	Response_name = map[int32]string{
		0: "NORMAL",
		1: "TIMEOUT",
		2: "ERROR",
	}
	Response_value = map[string]int32{
		"NORMAL":  0,
		"TIMEOUT": 1,
		"ERROR":   2,
	}
)

func (x Response) Enum() *Response {
	p := new(Response)
	*p = x
	return p
}

func (x Response) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Response) Descriptor() protoreflect.EnumDescriptor {
	return file_scenario_events_proto_enumTypes[0].Descriptor()
}

func (Response) Type() protoreflect.EnumType {
	return &file_scenario_events_proto_enumTypes[0]
}

func (x Response) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Response.Descriptor instead.
func (Response) EnumDescriptor() ([]byte, []int) {
	return file_scenario_events_proto_rawDescGZIP(), []int{0}
}

var File_scenario_events_proto protoreflect.FileDescriptor

var file_scenario_events_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x70, 0x61, 0x74, 0x65, 0x2e, 0x73,
	0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x2a, 0x2e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x02, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x2d, 0x72, 0x65,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x63, 0x2d, 0x65, 0x6d,
	0x75, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scenario_events_proto_rawDescOnce sync.Once
	file_scenario_events_proto_rawDescData = file_scenario_events_proto_rawDesc
)

func file_scenario_events_proto_rawDescGZIP() []byte {
	file_scenario_events_proto_rawDescOnce.Do(func() {
		file_scenario_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_scenario_events_proto_rawDescData)
	})
	return file_scenario_events_proto_rawDescData
}

var file_scenario_events_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_scenario_events_proto_goTypes = []interface{}{
	(Response)(0), // 0: apate.scenario.Response
}
var file_scenario_events_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_scenario_events_proto_init() }
func file_scenario_events_proto_init() {
	if File_scenario_events_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_scenario_events_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_scenario_events_proto_goTypes,
		DependencyIndexes: file_scenario_events_proto_depIdxs,
		EnumInfos:         file_scenario_events_proto_enumTypes,
	}.Build()
	File_scenario_events_proto = out.File
	file_scenario_events_proto_rawDesc = nil
	file_scenario_events_proto_goTypes = nil
	file_scenario_events_proto_depIdxs = nil
}
