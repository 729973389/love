// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: edgeProperties.proto

package protobuf

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

type EdgeProperties struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OS       string `protobuf:"bytes,1,opt,name=OS,proto3" json:"OS,omitempty"`
	ARC      string `protobuf:"bytes,2,opt,name=ARC,proto3" json:"ARC,omitempty"`
	CPUNum   int32  `protobuf:"varint,3,opt,name=CPUNum,proto3" json:"CPUNum,omitempty"`
	HostName string `protobuf:"bytes,4,opt,name=hostName,proto3" json:"hostName,omitempty"`
}

func (x *EdgeProperties) Reset() {
	*x = EdgeProperties{}
	if protoimpl.UnsafeEnabled {
		mi := &file_edgeProperties_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EdgeProperties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EdgeProperties) ProtoMessage() {}

func (x *EdgeProperties) ProtoReflect() protoreflect.Message {
	mi := &file_edgeProperties_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EdgeProperties.ProtoReflect.Descriptor instead.
func (*EdgeProperties) Descriptor() ([]byte, []int) {
	return file_edgeProperties_proto_rawDescGZIP(), []int{0}
}

func (x *EdgeProperties) GetOS() string {
	if x != nil {
		return x.OS
	}
	return ""
}

func (x *EdgeProperties) GetARC() string {
	if x != nil {
		return x.ARC
	}
	return ""
}

func (x *EdgeProperties) GetCPUNum() int32 {
	if x != nil {
		return x.CPUNum
	}
	return 0
}

func (x *EdgeProperties) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

var File_edgeProperties_proto protoreflect.FileDescriptor

var file_edgeProperties_proto_rawDesc = []byte{
	0x0a, 0x14, 0x65, 0x64, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x22, 0x66, 0x0a, 0x0e, 0x65, 0x64, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69,
	0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x53, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x4f, 0x53, 0x12, 0x10, 0x0a, 0x03, 0x41, 0x52, 0x43, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x41, 0x52, 0x43, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x50, 0x55, 0x4e, 0x75, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x43, 0x50, 0x55, 0x4e, 0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08,
	0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x02, 0x5a, 0x00, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_edgeProperties_proto_rawDescOnce sync.Once
	file_edgeProperties_proto_rawDescData = file_edgeProperties_proto_rawDesc
)

func file_edgeProperties_proto_rawDescGZIP() []byte {
	file_edgeProperties_proto_rawDescOnce.Do(func() {
		file_edgeProperties_proto_rawDescData = protoimpl.X.CompressGZIP(file_edgeProperties_proto_rawDescData)
	})
	return file_edgeProperties_proto_rawDescData
}

var file_edgeProperties_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_edgeProperties_proto_goTypes = []interface{}{
	(*EdgeProperties)(nil), // 0: protobuf.edgeProperties
}
var file_edgeProperties_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_edgeProperties_proto_init() }
func file_edgeProperties_proto_init() {
	if File_edgeProperties_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_edgeProperties_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EdgeProperties); i {
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
			RawDescriptor: file_edgeProperties_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_edgeProperties_proto_goTypes,
		DependencyIndexes: file_edgeProperties_proto_depIdxs,
		MessageInfos:      file_edgeProperties_proto_msgTypes,
	}.Build()
	File_edgeProperties_proto = out.File
	file_edgeProperties_proto_rawDesc = nil
	file_edgeProperties_proto_goTypes = nil
	file_edgeProperties_proto_depIdxs = nil
}
