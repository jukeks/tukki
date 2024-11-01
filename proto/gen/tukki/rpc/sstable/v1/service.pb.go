// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: tukki/rpc/sstable/v1/service.proto

package sstablev1

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

type GetSstableRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetSstableRequest) Reset() {
	*x = GetSstableRequest{}
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSstableRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSstableRequest) ProtoMessage() {}

func (x *GetSstableRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSstableRequest.ProtoReflect.Descriptor instead.
func (*GetSstableRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_sstable_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetSstableRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type SSTableRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value   string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Deleted bool   `protobuf:"varint,3,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *SSTableRecord) Reset() {
	*x = SSTableRecord{}
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SSTableRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SSTableRecord) ProtoMessage() {}

func (x *SSTableRecord) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SSTableRecord.ProtoReflect.Descriptor instead.
func (*SSTableRecord) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_sstable_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *SSTableRecord) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SSTableRecord) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SSTableRecord) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type GetSstableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record *SSTableRecord `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *GetSstableResponse) Reset() {
	*x = GetSstableResponse{}
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSstableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSstableResponse) ProtoMessage() {}

func (x *GetSstableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_sstable_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSstableResponse.ProtoReflect.Descriptor instead.
func (*GetSstableResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_sstable_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetSstableResponse) GetRecord() *SSTableRecord {
	if x != nil {
		return x.Record
	}
	return nil
}

var File_tukki_rpc_sstable_v1_service_proto protoreflect.FileDescriptor

var file_tukki_rpc_sstable_v1_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x73, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x23, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x51, 0x0a, 0x0d, 0x53, 0x53, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x22, 0x51, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x53, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x06, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x32, 0x75, 0x0a, 0x0e, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53, 0x73,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x27, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x73, 0x73, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0xd9, 0x01, 0x0a,
	0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x73,
	0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x6b, 0x65, 0x6b, 0x73, 0x2f, 0x74, 0x75, 0x6b,
	0x6b, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x73,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x54, 0x52, 0x53, 0xaa, 0x02, 0x14,
	0x54, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x14, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x5c, 0x52, 0x70, 0x63,
	0x5c, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x20, 0x54, 0x75,
	0x6b, 0x6b, 0x69, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x17, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x53, 0x73, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tukki_rpc_sstable_v1_service_proto_rawDescOnce sync.Once
	file_tukki_rpc_sstable_v1_service_proto_rawDescData = file_tukki_rpc_sstable_v1_service_proto_rawDesc
)

func file_tukki_rpc_sstable_v1_service_proto_rawDescGZIP() []byte {
	file_tukki_rpc_sstable_v1_service_proto_rawDescOnce.Do(func() {
		file_tukki_rpc_sstable_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_tukki_rpc_sstable_v1_service_proto_rawDescData)
	})
	return file_tukki_rpc_sstable_v1_service_proto_rawDescData
}

var file_tukki_rpc_sstable_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tukki_rpc_sstable_v1_service_proto_goTypes = []any{
	(*GetSstableRequest)(nil),  // 0: tukki.rpc.sstable.v1.GetSstableRequest
	(*SSTableRecord)(nil),      // 1: tukki.rpc.sstable.v1.SSTableRecord
	(*GetSstableResponse)(nil), // 2: tukki.rpc.sstable.v1.GetSstableResponse
}
var file_tukki_rpc_sstable_v1_service_proto_depIdxs = []int32{
	1, // 0: tukki.rpc.sstable.v1.GetSstableResponse.record:type_name -> tukki.rpc.sstable.v1.SSTableRecord
	0, // 1: tukki.rpc.sstable.v1.SstableService.GetSstable:input_type -> tukki.rpc.sstable.v1.GetSstableRequest
	2, // 2: tukki.rpc.sstable.v1.SstableService.GetSstable:output_type -> tukki.rpc.sstable.v1.GetSstableResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tukki_rpc_sstable_v1_service_proto_init() }
func file_tukki_rpc_sstable_v1_service_proto_init() {
	if File_tukki_rpc_sstable_v1_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tukki_rpc_sstable_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tukki_rpc_sstable_v1_service_proto_goTypes,
		DependencyIndexes: file_tukki_rpc_sstable_v1_service_proto_depIdxs,
		MessageInfos:      file_tukki_rpc_sstable_v1_service_proto_msgTypes,
	}.Build()
	File_tukki_rpc_sstable_v1_service_proto = out.File
	file_tukki_rpc_sstable_v1_service_proto_rawDesc = nil
	file_tukki_rpc_sstable_v1_service_proto_goTypes = nil
	file_tukki_rpc_sstable_v1_service_proto_depIdxs = nil
}
