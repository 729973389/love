// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: oneof.proto

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Switch:
	//	*Message_Author
	//	*Message_EdgeInfo
	//	*Message_DeviceInfo
	Switch isMessage_Switch `protobuf_oneof:"switch"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneof_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_oneof_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_oneof_proto_rawDescGZIP(), []int{0}
}

func (m *Message) GetSwitch() isMessage_Switch {
	if m != nil {
		return m.Switch
	}
	return nil
}

func (x *Message) GetAuthor() *Author {
	if x, ok := x.GetSwitch().(*Message_Author); ok {
		return x.Author
	}
	return nil
}

func (x *Message) GetEdgeInfo() *EdgeInfo {
	if x, ok := x.GetSwitch().(*Message_EdgeInfo); ok {
		return x.EdgeInfo
	}
	return nil
}

func (x *Message) GetDeviceInfo() *DeviceInfo {
	if x, ok := x.GetSwitch().(*Message_DeviceInfo); ok {
		return x.DeviceInfo
	}
	return nil
}

type isMessage_Switch interface {
	isMessage_Switch()
}

type Message_Author struct {
	Author *Author `protobuf:"bytes,1,opt,name=author,proto3,oneof"`
}

type Message_EdgeInfo struct {
	EdgeInfo *EdgeInfo `protobuf:"bytes,2,opt,name=edgeInfo,proto3,oneof"`
}

type Message_DeviceInfo struct {
	DeviceInfo *DeviceInfo `protobuf:"bytes,3,opt,name=deviceInfo,proto3,oneof"`
}

func (*Message_Author) isMessage_Switch() {}

func (*Message_EdgeInfo) isMessage_Switch() {}

func (*Message_DeviceInfo) isMessage_Switch() {}

var File_oneof_proto protoreflect.FileDescriptor

var file_oneof_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x0c, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x01, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12,
	0x30, 0x0a, 0x08, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x65, 0x64, 0x67,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x36, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x0a, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x08, 0x0a, 0x06, 0x73, 0x77, 0x69,
	0x74, 0x63, 0x68, 0x42, 0x02, 0x5a, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_oneof_proto_rawDescOnce sync.Once
	file_oneof_proto_rawDescData = file_oneof_proto_rawDesc
)

func file_oneof_proto_rawDescGZIP() []byte {
	file_oneof_proto_rawDescOnce.Do(func() {
		file_oneof_proto_rawDescData = protoimpl.X.CompressGZIP(file_oneof_proto_rawDescData)
	})
	return file_oneof_proto_rawDescData
}

var file_oneof_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_oneof_proto_goTypes = []interface{}{
	(*Message)(nil),    // 0: protobuf.message
	(*Author)(nil),     // 1: protobuf.author
	(*EdgeInfo)(nil),   // 2: protobuf.edgeInfo
	(*DeviceInfo)(nil), // 3: protobuf.deviceInfo
}
var file_oneof_proto_depIdxs = []int32{
	1, // 0: protobuf.message.author:type_name -> protobuf.author
	2, // 1: protobuf.message.edgeInfo:type_name -> protobuf.edgeInfo
	3, // 2: protobuf.message.deviceInfo:type_name -> protobuf.deviceInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_oneof_proto_init() }
func file_oneof_proto_init() {
	if File_oneof_proto != nil {
		return
	}
	file_author_proto_init()
	file_edgeInfo_proto_init()
	file_deviceInfo_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_oneof_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_oneof_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Author)(nil),
		(*Message_EdgeInfo)(nil),
		(*Message_DeviceInfo)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_oneof_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_oneof_proto_goTypes,
		DependencyIndexes: file_oneof_proto_depIdxs,
		MessageInfos:      file_oneof_proto_msgTypes,
	}.Build()
	File_oneof_proto = out.File
	file_oneof_proto_rawDesc = nil
	file_oneof_proto_goTypes = nil
	file_oneof_proto_depIdxs = nil
}
