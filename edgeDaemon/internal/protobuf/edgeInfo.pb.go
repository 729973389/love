// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: edgeInfo.proto

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

type EdgeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SerialNumber string      `protobuf:"bytes,1,opt,name=serialNumber,proto3" json:"serialNumber,omitempty"`
	Data         *SystemInfo `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *EdgeInfo) Reset() {
	*x = EdgeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_edgeInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EdgeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EdgeInfo) ProtoMessage() {}

func (x *EdgeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_edgeInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EdgeInfo.ProtoReflect.Descriptor instead.
func (*EdgeInfo) Descriptor() ([]byte, []int) {
	return file_edgeInfo_proto_rawDescGZIP(), []int{0}
}

func (x *EdgeInfo) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *EdgeInfo) GetData() *SystemInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_edgeInfo_proto protoreflect.FileDescriptor

var file_edgeInfo_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x0c, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x08, 0x65, 0x64, 0x67, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_edgeInfo_proto_rawDescOnce sync.Once
	file_edgeInfo_proto_rawDescData = file_edgeInfo_proto_rawDesc
)

func file_edgeInfo_proto_rawDescGZIP() []byte {
	file_edgeInfo_proto_rawDescOnce.Do(func() {
		file_edgeInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_edgeInfo_proto_rawDescData)
	})
	return file_edgeInfo_proto_rawDescData
}

var file_edgeInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_edgeInfo_proto_goTypes = []interface{}{
	(*EdgeInfo)(nil),   // 0: protobuf.edgeInfo
	(*SystemInfo)(nil), // 1: protobuf.SystemInfo
}
var file_edgeInfo_proto_depIdxs = []int32{
	1, // 0: protobuf.edgeInfo.data:type_name -> protobuf.SystemInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_edgeInfo_proto_init() }
func file_edgeInfo_proto_init() {
	if File_edgeInfo_proto != nil {
		return
	}
	file_system_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_edgeInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EdgeInfo); i {
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
			RawDescriptor: file_edgeInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_edgeInfo_proto_goTypes,
		DependencyIndexes: file_edgeInfo_proto_depIdxs,
		MessageInfos:      file_edgeInfo_proto_msgTypes,
	}.Build()
	File_edgeInfo_proto = out.File
	file_edgeInfo_proto_rawDesc = nil
	file_edgeInfo_proto_goTypes = nil
	file_edgeInfo_proto_depIdxs = nil
}
