// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: oneOf.proto

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
	//	*Message_DeviceCommand
	//	*Message_DeviceGister
	//	*Message_DeviceMap
	//	*Message_EdgeProperties
	Switch isMessage_Switch `protobuf_oneof:"switch"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneOf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_oneOf_proto_msgTypes[0]
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
	return file_oneOf_proto_rawDescGZIP(), []int{0}
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

func (x *Message) GetEdgeInfo() *InterfaceEdge {
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

func (x *Message) GetDeviceCommand() *DeviceCommand {
	if x, ok := x.GetSwitch().(*Message_DeviceCommand); ok {
		return x.DeviceCommand
	}
	return nil
}

func (x *Message) GetDeviceGister() *DeviceGister {
	if x, ok := x.GetSwitch().(*Message_DeviceGister); ok {
		return x.DeviceGister
	}
	return nil
}

func (x *Message) GetDeviceMap() *DeviceMap {
	if x, ok := x.GetSwitch().(*Message_DeviceMap); ok {
		return x.DeviceMap
	}
	return nil
}

func (x *Message) GetEdgeProperties() *EdgeProperties {
	if x, ok := x.GetSwitch().(*Message_EdgeProperties); ok {
		return x.EdgeProperties
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
	EdgeInfo *InterfaceEdge `protobuf:"bytes,2,opt,name=edgeInfo,proto3,oneof"`
}

type Message_DeviceInfo struct {
	DeviceInfo *DeviceInfo `protobuf:"bytes,3,opt,name=deviceInfo,proto3,oneof"`
}

type Message_DeviceCommand struct {
	DeviceCommand *DeviceCommand `protobuf:"bytes,4,opt,name=deviceCommand,proto3,oneof"`
}

type Message_DeviceGister struct {
	DeviceGister *DeviceGister `protobuf:"bytes,5,opt,name=deviceGister,proto3,oneof"`
}

type Message_DeviceMap struct {
	DeviceMap *DeviceMap `protobuf:"bytes,6,opt,name=deviceMap,proto3,oneof"`
}

type Message_EdgeProperties struct {
	EdgeProperties *EdgeProperties `protobuf:"bytes,7,opt,name=edgeProperties,proto3,oneof"`
}

func (*Message_Author) isMessage_Switch() {}

func (*Message_EdgeInfo) isMessage_Switch() {}

func (*Message_DeviceInfo) isMessage_Switch() {}

func (*Message_DeviceCommand) isMessage_Switch() {}

func (*Message_DeviceGister) isMessage_Switch() {}

func (*Message_DeviceMap) isMessage_Switch() {}

func (*Message_EdgeProperties) isMessage_Switch() {}

var File_oneOf_proto protoreflect.FileDescriptor

var file_oneOf_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x6e, 0x65, 0x4f, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x0c, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x45, 0x64, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x12, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x47, 0x69, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x61, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x65, 0x64, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x70,
	0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa6, 0x03, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x06, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x45, 0x64, 0x67, 0x65, 0x48,
	0x00, 0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x36, 0x0a, 0x0a, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x3f, 0x0a, 0x0d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x0d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x3c, 0x0a, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x47, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x47, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x47, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x33, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x61, 0x70, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x61, 0x70, 0x48, 0x00, 0x52, 0x09, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x4d, 0x61, 0x70, 0x12, 0x42, 0x0a, 0x0e, 0x65, 0x64, 0x67, 0x65, 0x50,
	0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x50,
	0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x48, 0x00, 0x52, 0x0e, 0x65, 0x64, 0x67,
	0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x73,
	0x77, 0x69, 0x74, 0x63, 0x68, 0x42, 0x02, 0x5a, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_oneOf_proto_rawDescOnce sync.Once
	file_oneOf_proto_rawDescData = file_oneOf_proto_rawDesc
)

func file_oneOf_proto_rawDescGZIP() []byte {
	file_oneOf_proto_rawDescOnce.Do(func() {
		file_oneOf_proto_rawDescData = protoimpl.X.CompressGZIP(file_oneOf_proto_rawDescData)
	})
	return file_oneOf_proto_rawDescData
}

var file_oneOf_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_oneOf_proto_goTypes = []interface{}{
	(*Message)(nil),        // 0: protobuf.message
	(*Author)(nil),         // 1: protobuf.author
	(*InterfaceEdge)(nil),  // 2: protobuf.InterfaceEdge
	(*DeviceInfo)(nil),     // 3: protobuf.deviceInfo
	(*DeviceCommand)(nil),  // 4: protobuf.deviceCommand
	(*DeviceGister)(nil),   // 5: protobuf.deviceGister
	(*DeviceMap)(nil),      // 6: protobuf.deviceMap
	(*EdgeProperties)(nil), // 7: protobuf.edgeProperties
}
var file_oneOf_proto_depIdxs = []int32{
	1, // 0: protobuf.message.author:type_name -> protobuf.author
	2, // 1: protobuf.message.edgeInfo:type_name -> protobuf.InterfaceEdge
	3, // 2: protobuf.message.deviceInfo:type_name -> protobuf.deviceInfo
	4, // 3: protobuf.message.deviceCommand:type_name -> protobuf.deviceCommand
	5, // 4: protobuf.message.deviceGister:type_name -> protobuf.deviceGister
	6, // 5: protobuf.message.deviceMap:type_name -> protobuf.deviceMap
	7, // 6: protobuf.message.edgeProperties:type_name -> protobuf.edgeProperties
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_oneOf_proto_init() }
func file_oneOf_proto_init() {
	if File_oneOf_proto != nil {
		return
	}
	file_author_proto_init()
	file_interfaceEdge_proto_init()
	file_deviceInfo_proto_init()
	file_deviceCommand_proto_init()
	file_deviceGister_proto_init()
	file_deviceMap_proto_init()
	file_edgeProperties_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_oneOf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	file_oneOf_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Author)(nil),
		(*Message_EdgeInfo)(nil),
		(*Message_DeviceInfo)(nil),
		(*Message_DeviceCommand)(nil),
		(*Message_DeviceGister)(nil),
		(*Message_DeviceMap)(nil),
		(*Message_EdgeProperties)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_oneOf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_oneOf_proto_goTypes,
		DependencyIndexes: file_oneOf_proto_depIdxs,
		MessageInfos:      file_oneOf_proto_msgTypes,
	}.Build()
	File_oneOf_proto = out.File
	file_oneOf_proto_rawDesc = nil
	file_oneOf_proto_goTypes = nil
	file_oneOf_proto_depIdxs = nil
}
