// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: deviceInfo.proto

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

type DeviceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceType string `protobuf:"bytes,1,opt,name=deviceType,proto3" json:"deviceType,omitempty"`
	Data       string `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Uuid       string `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *DeviceInfo) Reset() {
	*x = DeviceInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deviceInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceInfo) ProtoMessage() {}

func (x *DeviceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_deviceInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceInfo.ProtoReflect.Descriptor instead.
func (*DeviceInfo) Descriptor() ([]byte, []int) {
	return file_deviceInfo_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceInfo) GetDeviceType() string {
	if x != nil {
		return x.DeviceType
	}
	return ""
}

func (x *DeviceInfo) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *DeviceInfo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

var File_deviceInfo_proto protoreflect.FileDescriptor

var file_deviceInfo_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x54, 0x0a, 0x0a,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x42, 0x02, 0x5a, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deviceInfo_proto_rawDescOnce sync.Once
	file_deviceInfo_proto_rawDescData = file_deviceInfo_proto_rawDesc
)

func file_deviceInfo_proto_rawDescGZIP() []byte {
	file_deviceInfo_proto_rawDescOnce.Do(func() {
		file_deviceInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_deviceInfo_proto_rawDescData)
	})
	return file_deviceInfo_proto_rawDescData
}

var file_deviceInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_deviceInfo_proto_goTypes = []interface{}{
	(*DeviceInfo)(nil), // 0: protobuf.deviceInfo
}
var file_deviceInfo_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_deviceInfo_proto_init() }
func file_deviceInfo_proto_init() {
	if File_deviceInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_deviceInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceInfo); i {
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
			RawDescriptor: file_deviceInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_deviceInfo_proto_goTypes,
		DependencyIndexes: file_deviceInfo_proto_depIdxs,
		MessageInfos:      file_deviceInfo_proto_msgTypes,
	}.Build()
	File_deviceInfo_proto = out.File
	file_deviceInfo_proto_rawDesc = nil
	file_deviceInfo_proto_goTypes = nil
	file_deviceInfo_proto_depIdxs = nil
}