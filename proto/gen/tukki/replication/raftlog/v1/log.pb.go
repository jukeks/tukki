// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: tukki/replication/raftlog/v1/log.proto

package raftlogv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index      uint64                 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Term       uint64                 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Type       uint32                 `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Data       []byte                 `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Extensions []byte                 `protobuf:"bytes,5,opt,name=extensions,proto3" json:"extensions,omitempty"`
	AppendedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=appended_at,json=appendedAt,proto3" json:"appended_at,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	mi := &file_tukki_replication_raftlog_v1_log_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_replication_raftlog_v1_log_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_tukki_replication_raftlog_v1_log_proto_rawDescGZIP(), []int{0}
}

func (x *Log) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Log) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Log) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Log) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Log) GetExtensions() []byte {
	if x != nil {
		return x.Extensions
	}
	return nil
}

func (x *Log) GetAppendedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AppendedAt
	}
	return nil
}

var File_tukki_replication_raftlog_v1_log_proto protoreflect.FileDescriptor

var file_tukki_replication_raftlog_v1_log_proto_rawDesc = []byte{
	0x0a, 0x26, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6c,
	0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x61, 0x66, 0x74,
	0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb4, 0x01, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x42, 0x85,
	0x02, 0x0a, 0x20, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x6f, 0x67,
	0x2e, 0x76, 0x31, 0x42, 0x08, 0x4c, 0x6f, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x6b, 0x65,
	0x6b, 0x73, 0x2f, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x61, 0x66, 0x74,
	0x6c, 0x6f, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x54, 0x52, 0x52, 0xaa, 0x02, 0x1c, 0x54, 0x75,
	0x6b, 0x6b, 0x69, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x52, 0x61, 0x66, 0x74, 0x6c, 0x6f, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1c, 0x54, 0x75, 0x6b,
	0x6b, 0x69, 0x5c, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x52,
	0x61, 0x66, 0x74, 0x6c, 0x6f, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x28, 0x54, 0x75, 0x6b, 0x6b,
	0x69, 0x5c, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x52, 0x61,
	0x66, 0x74, 0x6c, 0x6f, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1f, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x3a, 0x3a, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x52, 0x61, 0x66, 0x74, 0x6c,
	0x6f, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tukki_replication_raftlog_v1_log_proto_rawDescOnce sync.Once
	file_tukki_replication_raftlog_v1_log_proto_rawDescData = file_tukki_replication_raftlog_v1_log_proto_rawDesc
)

func file_tukki_replication_raftlog_v1_log_proto_rawDescGZIP() []byte {
	file_tukki_replication_raftlog_v1_log_proto_rawDescOnce.Do(func() {
		file_tukki_replication_raftlog_v1_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_tukki_replication_raftlog_v1_log_proto_rawDescData)
	})
	return file_tukki_replication_raftlog_v1_log_proto_rawDescData
}

var file_tukki_replication_raftlog_v1_log_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tukki_replication_raftlog_v1_log_proto_goTypes = []any{
	(*Log)(nil),                   // 0: tukki.replication.raftlog.v1.Log
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_tukki_replication_raftlog_v1_log_proto_depIdxs = []int32{
	1, // 0: tukki.replication.raftlog.v1.Log.appended_at:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tukki_replication_raftlog_v1_log_proto_init() }
func file_tukki_replication_raftlog_v1_log_proto_init() {
	if File_tukki_replication_raftlog_v1_log_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tukki_replication_raftlog_v1_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tukki_replication_raftlog_v1_log_proto_goTypes,
		DependencyIndexes: file_tukki_replication_raftlog_v1_log_proto_depIdxs,
		MessageInfos:      file_tukki_replication_raftlog_v1_log_proto_msgTypes,
	}.Build()
	File_tukki_replication_raftlog_v1_log_proto = out.File
	file_tukki_replication_raftlog_v1_log_proto_rawDesc = nil
	file_tukki_replication_raftlog_v1_log_proto_goTypes = nil
	file_tukki_replication_raftlog_v1_log_proto_depIdxs = nil
}