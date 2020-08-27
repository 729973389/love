// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: diskInfo.proto

package pcp

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

type DiskInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Used string `protobuf:"bytes,1,opt,name=Used,proto3" json:"Used,omitempty"`
	Free string `protobuf:"bytes,2,opt,name=Free,proto3" json:"Free,omitempty"`
	Rate string `protobuf:"bytes,3,opt,name=Rate,proto3" json:"Rate,omitempty"`
}

func (x *DiskInfo) Reset() {
	*x = DiskInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diskInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiskInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiskInfo) ProtoMessage() {}

func (x *DiskInfo) ProtoReflect() protoreflect.Message {
	mi := &file_diskInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiskInfo.ProtoReflect.Descriptor instead.
func (*DiskInfo) Descriptor() ([]byte, []int) {
	return file_diskInfo_proto_rawDescGZIP(), []int{0}
}

func (x *DiskInfo) GetUsed() string {
	if x != nil {
		return x.Used
	}
	return ""
}

func (x *DiskInfo) GetFree() string {
	if x != nil {
		return x.Free
	}
	return ""
}

func (x *DiskInfo) GetRate() string {
	if x != nil {
		return x.Rate
	}
	return ""
}

var File_diskInfo_proto protoreflect.FileDescriptor

var file_diskInfo_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x64, 0x69, 0x73, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x46, 0x0a, 0x08, 0x44, 0x69,
	0x73, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x73, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72,
	0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x65, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x52, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x52, 0x61,
	0x74, 0x65, 0x42, 0x02, 0x5a, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_diskInfo_proto_rawDescOnce sync.Once
	file_diskInfo_proto_rawDescData = file_diskInfo_proto_rawDesc
)

func file_diskInfo_proto_rawDescGZIP() []byte {
	file_diskInfo_proto_rawDescOnce.Do(func() {
		file_diskInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_diskInfo_proto_rawDescData)
	})
	return file_diskInfo_proto_rawDescData
}

var file_diskInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_diskInfo_proto_goTypes = []interface{}{
	(*DiskInfo)(nil), // 0: protobuf.DiskInfo
}
var file_diskInfo_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_diskInfo_proto_init() }
func file_diskInfo_proto_init() {
	if File_diskInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_diskInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiskInfo); i {
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
			RawDescriptor: file_diskInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_diskInfo_proto_goTypes,
		DependencyIndexes: file_diskInfo_proto_depIdxs,
		MessageInfos:      file_diskInfo_proto_msgTypes,
	}.Build()
	File_diskInfo_proto = out.File
	file_diskInfo_proto_rawDesc = nil
	file_diskInfo_proto_goTypes = nil
	file_diskInfo_proto_depIdxs = nil
}
