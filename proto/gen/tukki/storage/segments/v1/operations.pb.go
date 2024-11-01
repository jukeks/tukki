// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: tukki/storage/segments/v1/operations.proto

package segmentsv1

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

type Segment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Filename        string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	MembersFilename string `protobuf:"bytes,3,opt,name=members_filename,json=membersFilename,proto3" json:"members_filename,omitempty"`
	IndexFilename   string `protobuf:"bytes,4,opt,name=index_filename,json=indexFilename,proto3" json:"index_filename,omitempty"`
}

func (x *Segment) Reset() {
	*x = Segment{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Segment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Segment) ProtoMessage() {}

func (x *Segment) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Segment.ProtoReflect.Descriptor instead.
func (*Segment) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{0}
}

func (x *Segment) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Segment) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Segment) GetMembersFilename() string {
	if x != nil {
		return x.MembersFilename
	}
	return ""
}

func (x *Segment) GetIndexFilename() string {
	if x != nil {
		return x.IndexFilename
	}
	return ""
}

type LiveSegment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Segment     *Segment `protobuf:"bytes,1,opt,name=segment,proto3" json:"segment,omitempty"`
	WalFilename string   `protobuf:"bytes,2,opt,name=wal_filename,json=walFilename,proto3" json:"wal_filename,omitempty"`
}

func (x *LiveSegment) Reset() {
	*x = LiveSegment{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LiveSegment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LiveSegment) ProtoMessage() {}

func (x *LiveSegment) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LiveSegment.ProtoReflect.Descriptor instead.
func (*LiveSegment) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{1}
}

func (x *LiveSegment) GetSegment() *Segment {
	if x != nil {
		return x.Segment
	}
	return nil
}

func (x *LiveSegment) GetWalFilename() string {
	if x != nil {
		return x.WalFilename
	}
	return ""
}

type AddSegment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CompletingSegment *LiveSegment `protobuf:"bytes,1,opt,name=completing_segment,json=completingSegment,proto3,oneof" json:"completing_segment,omitempty"`
	NewSegment        *LiveSegment `protobuf:"bytes,2,opt,name=new_segment,json=newSegment,proto3" json:"new_segment,omitempty"`
}

func (x *AddSegment) Reset() {
	*x = AddSegment{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddSegment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddSegment) ProtoMessage() {}

func (x *AddSegment) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddSegment.ProtoReflect.Descriptor instead.
func (*AddSegment) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{2}
}

func (x *AddSegment) GetCompletingSegment() *LiveSegment {
	if x != nil {
		return x.CompletingSegment
	}
	return nil
}

func (x *AddSegment) GetNewSegment() *LiveSegment {
	if x != nil {
		return x.NewSegment
	}
	return nil
}

type MergeSegments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewSegment      *Segment   `protobuf:"bytes,1,opt,name=new_segment,json=newSegment,proto3" json:"new_segment,omitempty"`
	SegmentsToMerge []*Segment `protobuf:"bytes,2,rep,name=segments_to_merge,json=segmentsToMerge,proto3" json:"segments_to_merge,omitempty"`
}

func (x *MergeSegments) Reset() {
	*x = MergeSegments{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MergeSegments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeSegments) ProtoMessage() {}

func (x *MergeSegments) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeSegments.ProtoReflect.Descriptor instead.
func (*MergeSegments) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{3}
}

func (x *MergeSegments) GetNewSegment() *Segment {
	if x != nil {
		return x.NewSegment
	}
	return nil
}

func (x *MergeSegments) GetSegmentsToMerge() []*Segment {
	if x != nil {
		return x.SegmentsToMerge
	}
	return nil
}

type SegmentOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Operation:
	//
	//	*SegmentOperation_Add
	//	*SegmentOperation_Merge
	Operation isSegmentOperation_Operation `protobuf_oneof:"operation"`
}

func (x *SegmentOperation) Reset() {
	*x = SegmentOperation{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SegmentOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentOperation) ProtoMessage() {}

func (x *SegmentOperation) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentOperation.ProtoReflect.Descriptor instead.
func (*SegmentOperation) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{4}
}

func (x *SegmentOperation) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (m *SegmentOperation) GetOperation() isSegmentOperation_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (x *SegmentOperation) GetAdd() *AddSegment {
	if x, ok := x.GetOperation().(*SegmentOperation_Add); ok {
		return x.Add
	}
	return nil
}

func (x *SegmentOperation) GetMerge() *MergeSegments {
	if x, ok := x.GetOperation().(*SegmentOperation_Merge); ok {
		return x.Merge
	}
	return nil
}

type isSegmentOperation_Operation interface {
	isSegmentOperation_Operation()
}

type SegmentOperation_Add struct {
	Add *AddSegment `protobuf:"bytes,2,opt,name=add,proto3,oneof"`
}

type SegmentOperation_Merge struct {
	Merge *MergeSegments `protobuf:"bytes,3,opt,name=merge,proto3,oneof"`
}

func (*SegmentOperation_Add) isSegmentOperation_Operation() {}

func (*SegmentOperation_Merge) isSegmentOperation_Operation() {}

type SegmentOperationJournalEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Entry:
	//
	//	*SegmentOperationJournalEntry_Started
	//	*SegmentOperationJournalEntry_Completed
	Entry isSegmentOperationJournalEntry_Entry `protobuf_oneof:"entry"`
}

func (x *SegmentOperationJournalEntry) Reset() {
	*x = SegmentOperationJournalEntry{}
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SegmentOperationJournalEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentOperationJournalEntry) ProtoMessage() {}

func (x *SegmentOperationJournalEntry) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_operations_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentOperationJournalEntry.ProtoReflect.Descriptor instead.
func (*SegmentOperationJournalEntry) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_operations_proto_rawDescGZIP(), []int{5}
}

func (m *SegmentOperationJournalEntry) GetEntry() isSegmentOperationJournalEntry_Entry {
	if m != nil {
		return m.Entry
	}
	return nil
}

func (x *SegmentOperationJournalEntry) GetStarted() *SegmentOperation {
	if x, ok := x.GetEntry().(*SegmentOperationJournalEntry_Started); ok {
		return x.Started
	}
	return nil
}

func (x *SegmentOperationJournalEntry) GetCompleted() uint64 {
	if x, ok := x.GetEntry().(*SegmentOperationJournalEntry_Completed); ok {
		return x.Completed
	}
	return 0
}

type isSegmentOperationJournalEntry_Entry interface {
	isSegmentOperationJournalEntry_Entry()
}

type SegmentOperationJournalEntry_Started struct {
	Started *SegmentOperation `protobuf:"bytes,1,opt,name=started,proto3,oneof"`
}

type SegmentOperationJournalEntry_Completed struct {
	Completed uint64 `protobuf:"varint,2,opt,name=completed,proto3,oneof"`
}

func (*SegmentOperationJournalEntry_Started) isSegmentOperationJournalEntry_Entry() {}

func (*SegmentOperationJournalEntry_Completed) isSegmentOperationJournalEntry_Entry() {}

var File_tukki_storage_segments_v1_operations_proto protoreflect.FileDescriptor

var file_tukki_storage_segments_v1_operations_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f,
	0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x74, 0x75,
	0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x87, 0x01, 0x0a, 0x07, 0x53, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x29, 0x0a, 0x10, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x6e, 0x0a, 0x0b, 0x4c, 0x69, 0x76, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x3c, 0x0a, 0x07, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x77, 0x61, 0x6c, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x61, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0xc8, 0x01, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x5a, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x74,
	0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x76, 0x65, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x11, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x47, 0x0a, 0x0b,
	0x6e, 0x65, 0x77, 0x5f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x76, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x53, 0x65,
	0x67, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x15, 0x0a, 0x13, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0xa4, 0x01, 0x0a,
	0x0d, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x43,
	0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x53, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x4e, 0x0a, 0x11, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f,
	0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x0f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x54, 0x6f, 0x4d, 0x65,
	0x72, 0x67, 0x65, 0x22, 0xac, 0x01, 0x0a, 0x10, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x03,
	0x61, 0x64, 0x64, 0x12, 0x40, 0x0a, 0x05, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x28, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x65, 0x72, 0x67, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x48, 0x00, 0x52, 0x05,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x90, 0x01, 0x0a, 0x1c, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x47, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x48, 0x00, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x09,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48,
	0x00, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x42, 0x07, 0x0a, 0x05,
	0x65, 0x6e, 0x74, 0x72, 0x79, 0x42, 0xfb, 0x01, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x75,
	0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x6b, 0x65, 0x6b, 0x73, 0x2f, 0x74, 0x75,
	0x6b, 0x6b, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x2f, 0x76, 0x31, 0x3b, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x54, 0x53, 0x53, 0xaa, 0x02, 0x19, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x19, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x5c, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x5c, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x25, 0x54,
	0x75, 0x6b, 0x6b, 0x69, 0x5c, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5c, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c, 0x54, 0x75, 0x6b, 0x6b, 0x69, 0x3a, 0x3a, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x3a, 0x3a, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tukki_storage_segments_v1_operations_proto_rawDescOnce sync.Once
	file_tukki_storage_segments_v1_operations_proto_rawDescData = file_tukki_storage_segments_v1_operations_proto_rawDesc
)

func file_tukki_storage_segments_v1_operations_proto_rawDescGZIP() []byte {
	file_tukki_storage_segments_v1_operations_proto_rawDescOnce.Do(func() {
		file_tukki_storage_segments_v1_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_tukki_storage_segments_v1_operations_proto_rawDescData)
	})
	return file_tukki_storage_segments_v1_operations_proto_rawDescData
}

var file_tukki_storage_segments_v1_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tukki_storage_segments_v1_operations_proto_goTypes = []any{
	(*Segment)(nil),                      // 0: tukki.storage.segments.v1.Segment
	(*LiveSegment)(nil),                  // 1: tukki.storage.segments.v1.LiveSegment
	(*AddSegment)(nil),                   // 2: tukki.storage.segments.v1.AddSegment
	(*MergeSegments)(nil),                // 3: tukki.storage.segments.v1.MergeSegments
	(*SegmentOperation)(nil),             // 4: tukki.storage.segments.v1.SegmentOperation
	(*SegmentOperationJournalEntry)(nil), // 5: tukki.storage.segments.v1.SegmentOperationJournalEntry
}
var file_tukki_storage_segments_v1_operations_proto_depIdxs = []int32{
	0, // 0: tukki.storage.segments.v1.LiveSegment.segment:type_name -> tukki.storage.segments.v1.Segment
	1, // 1: tukki.storage.segments.v1.AddSegment.completing_segment:type_name -> tukki.storage.segments.v1.LiveSegment
	1, // 2: tukki.storage.segments.v1.AddSegment.new_segment:type_name -> tukki.storage.segments.v1.LiveSegment
	0, // 3: tukki.storage.segments.v1.MergeSegments.new_segment:type_name -> tukki.storage.segments.v1.Segment
	0, // 4: tukki.storage.segments.v1.MergeSegments.segments_to_merge:type_name -> tukki.storage.segments.v1.Segment
	2, // 5: tukki.storage.segments.v1.SegmentOperation.add:type_name -> tukki.storage.segments.v1.AddSegment
	3, // 6: tukki.storage.segments.v1.SegmentOperation.merge:type_name -> tukki.storage.segments.v1.MergeSegments
	4, // 7: tukki.storage.segments.v1.SegmentOperationJournalEntry.started:type_name -> tukki.storage.segments.v1.SegmentOperation
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_tukki_storage_segments_v1_operations_proto_init() }
func file_tukki_storage_segments_v1_operations_proto_init() {
	if File_tukki_storage_segments_v1_operations_proto != nil {
		return
	}
	file_tukki_storage_segments_v1_operations_proto_msgTypes[2].OneofWrappers = []any{}
	file_tukki_storage_segments_v1_operations_proto_msgTypes[4].OneofWrappers = []any{
		(*SegmentOperation_Add)(nil),
		(*SegmentOperation_Merge)(nil),
	}
	file_tukki_storage_segments_v1_operations_proto_msgTypes[5].OneofWrappers = []any{
		(*SegmentOperationJournalEntry_Started)(nil),
		(*SegmentOperationJournalEntry_Completed)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tukki_storage_segments_v1_operations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tukki_storage_segments_v1_operations_proto_goTypes,
		DependencyIndexes: file_tukki_storage_segments_v1_operations_proto_depIdxs,
		MessageInfos:      file_tukki_storage_segments_v1_operations_proto_msgTypes,
	}.Build()
	File_tukki_storage_segments_v1_operations_proto = out.File
	file_tukki_storage_segments_v1_operations_proto_rawDesc = nil
	file_tukki_storage_segments_v1_operations_proto_goTypes = nil
	file_tukki_storage_segments_v1_operations_proto_depIdxs = nil
}
