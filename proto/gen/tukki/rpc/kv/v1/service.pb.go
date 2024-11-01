// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: tukki/rpc/kv/v1/service.proto

package kvv1

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

type QueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *QueryRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type QueryRangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min string `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max string `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *QueryRangeRequest) Reset() {
	*x = QueryRangeRequest{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRangeRequest) ProtoMessage() {}

func (x *QueryRangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRangeRequest.ProtoReflect.Descriptor instead.
func (*QueryRangeRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *QueryRangeRequest) GetMin() string {
	if x != nil {
		return x.Min
	}
	return ""
}

func (x *QueryRangeRequest) GetMax() string {
	if x != nil {
		return x.Max
	}
	return ""
}

type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pair *KvPair `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *SetRequest) GetPair() *KvPair {
	if x != nil {
		return x.Pair
	}
	return nil
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteRangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min string `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max string `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *DeleteRangeRequest) Reset() {
	*x = DeleteRangeRequest{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRangeRequest) ProtoMessage() {}

func (x *DeleteRangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRangeRequest.ProtoReflect.Descriptor instead.
func (*DeleteRangeRequest) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRangeRequest) GetMin() string {
	if x != nil {
		return x.Min
	}
	return ""
}

func (x *DeleteRangeRequest) GetMax() string {
	if x != nil {
		return x.Max
	}
	return ""
}

type KvPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KvPair) Reset() {
	*x = KvPair{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KvPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KvPair) ProtoMessage() {}

func (x *KvPair) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KvPair.ProtoReflect.Descriptor instead.
func (*KvPair) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *KvPair) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KvPair) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type QueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*QueryResponse_Error
	//	*QueryResponse_Pair
	Value isQueryResponse_Value `protobuf_oneof:"value"`
}

func (x *QueryResponse) Reset() {
	*x = QueryResponse{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResponse) ProtoMessage() {}

func (x *QueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResponse.ProtoReflect.Descriptor instead.
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{7}
}

func (m *QueryResponse) GetValue() isQueryResponse_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *QueryResponse) GetError() *Error {
	if x, ok := x.GetValue().(*QueryResponse_Error); ok {
		return x.Error
	}
	return nil
}

func (x *QueryResponse) GetPair() *KvPair {
	if x, ok := x.GetValue().(*QueryResponse_Pair); ok {
		return x.Pair
	}
	return nil
}

type isQueryResponse_Value interface {
	isQueryResponse_Value()
}

type QueryResponse_Error struct {
	Error *Error `protobuf:"bytes,1,opt,name=error,proto3,oneof"`
}

type QueryResponse_Pair struct {
	Pair *KvPair `protobuf:"bytes,2,opt,name=pair,proto3,oneof"`
}

func (*QueryResponse_Error) isQueryResponse_Value() {}

func (*QueryResponse_Pair) isQueryResponse_Value() {}

type QueryRangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pair *KvPair `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
}

func (x *QueryRangeResponse) Reset() {
	*x = QueryRangeResponse{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRangeResponse) ProtoMessage() {}

func (x *QueryRangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRangeResponse.ProtoReflect.Descriptor instead.
func (*QueryRangeResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{8}
}

func (x *QueryRangeResponse) GetPair() *KvPair {
	if x != nil {
		return x.Pair
	}
	return nil
}

type SetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SetResponse) Reset() {
	*x = SetResponse{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResponse) ProtoMessage() {}

func (x *SetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResponse.ProtoReflect.Descriptor instead.
func (*SetResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{9}
}

func (x *SetResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type DeleteRangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error   *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Deleted uint64 `protobuf:"varint,2,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *DeleteRangeResponse) Reset() {
	*x = DeleteRangeResponse{}
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRangeResponse) ProtoMessage() {}

func (x *DeleteRangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_rpc_kv_v1_service_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRangeResponse.ProtoReflect.Descriptor instead.
func (*DeleteRangeResponse) Descriptor() ([]byte, []int) {
	return file_tukki_rpc_kv_v1_service_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteRangeResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *DeleteRangeResponse) GetDeleted() uint64 {
	if x != nil {
		return x.Deleted
	}
	return 0
}

var File_tukki_rpc_kv_v1_service_proto protoreflect.FileDescriptor

var file_tukki_rpc_kv_v1_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x6b, 0x76, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0f, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31,
	0x22, 0x20, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x22, 0x37, 0x0a, 0x11, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x22, 0x39, 0x0a, 0x0a, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x70, 0x61, 0x69,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x76, 0x50, 0x61, 0x69, 0x72,
	0x52, 0x04, 0x70, 0x61, 0x69, 0x72, 0x22, 0x21, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x38, 0x0a, 0x12, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x69,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x61, 0x78, 0x22, 0x30, 0x0a, 0x06, 0x4b, 0x76, 0x50, 0x61, 0x69, 0x72, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x77, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x2d, 0x0a, 0x04, 0x70, 0x61, 0x69,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x76, 0x50, 0x61, 0x69, 0x72,
	0x48, 0x00, 0x52, 0x04, 0x70, 0x61, 0x69, 0x72, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x41, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x70, 0x61, 0x69, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x76, 0x50, 0x61, 0x69, 0x72, 0x52, 0x04,
	0x70, 0x61, 0x69, 0x72, 0x22, 0x3b, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b,
	0x76, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x3e, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b,
	0x76, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x5d, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x32, 0x9d, 0x03, 0x0a, 0x09, 0x4b, 0x76, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48,
	0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1d, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x22, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x74, 0x75, 0x6b,
	0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x42, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x74, 0x75, 0x6b,
	0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x1e, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x12, 0x23, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0xb6, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x6b, 0x76, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x6b, 0x65, 0x6b, 0x73, 0x2f, 0x74, 0x75, 0x6b, 0x6b,
	0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x72, 0x70,
	0x63, 0x2f, 0x6b, 0x76, 0x2f, 0x76, 0x31, 0x3b, 0x6b, 0x76, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x54,
	0x52, 0x4b, 0xaa, 0x02, 0x0f, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x4b,
	0x76, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0f, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x5c, 0x52, 0x70, 0x63,
	0x5c, 0x4b, 0x76, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1b, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x5c, 0x52,
	0x70, 0x63, 0x5c, 0x4b, 0x76, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x3a, 0x3a, 0x52, 0x70,
	0x63, 0x3a, 0x3a, 0x4b, 0x76, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_tukki_rpc_kv_v1_service_proto_rawDescOnce sync.Once
	file_tukki_rpc_kv_v1_service_proto_rawDescData = file_tukki_rpc_kv_v1_service_proto_rawDesc
)

func file_tukki_rpc_kv_v1_service_proto_rawDescGZIP() []byte {
	file_tukki_rpc_kv_v1_service_proto_rawDescOnce.Do(func() {
		file_tukki_rpc_kv_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_tukki_rpc_kv_v1_service_proto_rawDescData)
	})
	return file_tukki_rpc_kv_v1_service_proto_rawDescData
}

var file_tukki_rpc_kv_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_tukki_rpc_kv_v1_service_proto_goTypes = []any{
	(*QueryRequest)(nil),        // 0: tukki.rpc.kv.v1.QueryRequest
	(*QueryRangeRequest)(nil),   // 1: tukki.rpc.kv.v1.QueryRangeRequest
	(*SetRequest)(nil),          // 2: tukki.rpc.kv.v1.SetRequest
	(*DeleteRequest)(nil),       // 3: tukki.rpc.kv.v1.DeleteRequest
	(*DeleteRangeRequest)(nil),  // 4: tukki.rpc.kv.v1.DeleteRangeRequest
	(*KvPair)(nil),              // 5: tukki.rpc.kv.v1.KvPair
	(*Error)(nil),               // 6: tukki.rpc.kv.v1.Error
	(*QueryResponse)(nil),       // 7: tukki.rpc.kv.v1.QueryResponse
	(*QueryRangeResponse)(nil),  // 8: tukki.rpc.kv.v1.QueryRangeResponse
	(*SetResponse)(nil),         // 9: tukki.rpc.kv.v1.SetResponse
	(*DeleteResponse)(nil),      // 10: tukki.rpc.kv.v1.DeleteResponse
	(*DeleteRangeResponse)(nil), // 11: tukki.rpc.kv.v1.DeleteRangeResponse
}
var file_tukki_rpc_kv_v1_service_proto_depIdxs = []int32{
	5,  // 0: tukki.rpc.kv.v1.SetRequest.pair:type_name -> tukki.rpc.kv.v1.KvPair
	6,  // 1: tukki.rpc.kv.v1.QueryResponse.error:type_name -> tukki.rpc.kv.v1.Error
	5,  // 2: tukki.rpc.kv.v1.QueryResponse.pair:type_name -> tukki.rpc.kv.v1.KvPair
	5,  // 3: tukki.rpc.kv.v1.QueryRangeResponse.pair:type_name -> tukki.rpc.kv.v1.KvPair
	6,  // 4: tukki.rpc.kv.v1.SetResponse.error:type_name -> tukki.rpc.kv.v1.Error
	6,  // 5: tukki.rpc.kv.v1.DeleteResponse.error:type_name -> tukki.rpc.kv.v1.Error
	6,  // 6: tukki.rpc.kv.v1.DeleteRangeResponse.error:type_name -> tukki.rpc.kv.v1.Error
	0,  // 7: tukki.rpc.kv.v1.KvService.Query:input_type -> tukki.rpc.kv.v1.QueryRequest
	1,  // 8: tukki.rpc.kv.v1.KvService.QueryRange:input_type -> tukki.rpc.kv.v1.QueryRangeRequest
	2,  // 9: tukki.rpc.kv.v1.KvService.Set:input_type -> tukki.rpc.kv.v1.SetRequest
	3,  // 10: tukki.rpc.kv.v1.KvService.Delete:input_type -> tukki.rpc.kv.v1.DeleteRequest
	4,  // 11: tukki.rpc.kv.v1.KvService.DeleteRange:input_type -> tukki.rpc.kv.v1.DeleteRangeRequest
	7,  // 12: tukki.rpc.kv.v1.KvService.Query:output_type -> tukki.rpc.kv.v1.QueryResponse
	8,  // 13: tukki.rpc.kv.v1.KvService.QueryRange:output_type -> tukki.rpc.kv.v1.QueryRangeResponse
	9,  // 14: tukki.rpc.kv.v1.KvService.Set:output_type -> tukki.rpc.kv.v1.SetResponse
	10, // 15: tukki.rpc.kv.v1.KvService.Delete:output_type -> tukki.rpc.kv.v1.DeleteResponse
	11, // 16: tukki.rpc.kv.v1.KvService.DeleteRange:output_type -> tukki.rpc.kv.v1.DeleteRangeResponse
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_tukki_rpc_kv_v1_service_proto_init() }
func file_tukki_rpc_kv_v1_service_proto_init() {
	if File_tukki_rpc_kv_v1_service_proto != nil {
		return
	}
	file_tukki_rpc_kv_v1_service_proto_msgTypes[7].OneofWrappers = []any{
		(*QueryResponse_Error)(nil),
		(*QueryResponse_Pair)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tukki_rpc_kv_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tukki_rpc_kv_v1_service_proto_goTypes,
		DependencyIndexes: file_tukki_rpc_kv_v1_service_proto_depIdxs,
		MessageInfos:      file_tukki_rpc_kv_v1_service_proto_msgTypes,
	}.Build()
	File_tukki_rpc_kv_v1_service_proto = out.File
	file_tukki_rpc_kv_v1_service_proto_rawDesc = nil
	file_tukki_rpc_kv_v1_service_proto_goTypes = nil
	file_tukki_rpc_kv_v1_service_proto_depIdxs = nil
}
