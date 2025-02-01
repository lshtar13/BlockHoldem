// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v3.12.4
// source: global/vers.proto

package global

import (
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

type VersReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Addr          string                 `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Version       int64                  `protobuf:"varint,2,opt,name=Version,proto3" json:"Version,omitempty"`
	BestHeight    int64                  `protobuf:"varint,3,opt,name=BestHeight,proto3" json:"BestHeight,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersReq) Reset() {
	*x = VersReq{}
	mi := &file_global_vers_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersReq) ProtoMessage() {}

func (x *VersReq) ProtoReflect() protoreflect.Message {
	mi := &file_global_vers_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersReq.ProtoReflect.Descriptor instead.
func (*VersReq) Descriptor() ([]byte, []int) {
	return file_global_vers_proto_rawDescGZIP(), []int{0}
}

func (x *VersReq) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *VersReq) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *VersReq) GetBestHeight() int64 {
	if x != nil {
		return x.BestHeight
	}
	return 0
}

type VersRet struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       int64                  `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	BestHeight    int64                  `protobuf:"varint,2,opt,name=BestHeight,proto3" json:"BestHeight,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersRet) Reset() {
	*x = VersRet{}
	mi := &file_global_vers_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersRet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersRet) ProtoMessage() {}

func (x *VersRet) ProtoReflect() protoreflect.Message {
	mi := &file_global_vers_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersRet.ProtoReflect.Descriptor instead.
func (*VersRet) Descriptor() ([]byte, []int) {
	return file_global_vers_proto_rawDescGZIP(), []int{1}
}

func (x *VersRet) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *VersRet) GetBestHeight() int64 {
	if x != nil {
		return x.BestHeight
	}
	return 0
}

var File_global_vers_proto protoreflect.FileDescriptor

var file_global_vers_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x22, 0x57, 0x0a, 0x07, 0x56,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x73, 0x74, 0x48, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x42, 0x65, 0x73, 0x74, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x22, 0x43, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x52, 0x65, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x73,
	0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x42,
	0x65, 0x73, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x32, 0x36, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x53, 0x72, 0x76, 0x12, 0x2b, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x56, 0x65, 0x72, 0x73, 0x12,
	0x0f, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71,
	0x1a, 0x0f, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x74, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6c, 0x73, 0x68, 0x74, 0x61, 0x72, 0x31, 0x33, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_global_vers_proto_rawDescOnce sync.Once
	file_global_vers_proto_rawDescData = file_global_vers_proto_rawDesc
)

func file_global_vers_proto_rawDescGZIP() []byte {
	file_global_vers_proto_rawDescOnce.Do(func() {
		file_global_vers_proto_rawDescData = protoimpl.X.CompressGZIP(file_global_vers_proto_rawDescData)
	})
	return file_global_vers_proto_rawDescData
}

var file_global_vers_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_global_vers_proto_goTypes = []any{
	(*VersReq)(nil), // 0: global.VersReq
	(*VersRet)(nil), // 1: global.VersRet
}
var file_global_vers_proto_depIdxs = []int32{
	0, // 0: global.versSrv.ReqVers:input_type -> global.VersReq
	1, // 1: global.versSrv.ReqVers:output_type -> global.VersRet
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_global_vers_proto_init() }
func file_global_vers_proto_init() {
	if File_global_vers_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_global_vers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_global_vers_proto_goTypes,
		DependencyIndexes: file_global_vers_proto_depIdxs,
		MessageInfos:      file_global_vers_proto_msgTypes,
	}.Build()
	File_global_vers_proto = out.File
	file_global_vers_proto_rawDesc = nil
	file_global_vers_proto_goTypes = nil
	file_global_vers_proto_depIdxs = nil
}
