// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: tukki/storage/segments/v1/segments.proto

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

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *Segment) Reset() {
	*x = Segment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Segment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Segment) ProtoMessage() {}

func (x *Segment) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_tukki_storage_segments_v1_segments_proto_rawDescGZIP(), []int{0}
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

type SegmentStarted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	JournalFilename string `protobuf:"bytes,2,opt,name=journal_filename,json=journalFilename,proto3" json:"journal_filename,omitempty"`
}

func (x *SegmentStarted) Reset() {
	*x = SegmentStarted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentStarted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentStarted) ProtoMessage() {}

func (x *SegmentStarted) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentStarted.ProtoReflect.Descriptor instead.
func (*SegmentStarted) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_segments_proto_rawDescGZIP(), []int{1}
}

func (x *SegmentStarted) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SegmentStarted) GetJournalFilename() string {
	if x != nil {
		return x.JournalFilename
	}
	return ""
}

type SegmentAdded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Segment *Segment `protobuf:"bytes,1,opt,name=segment,proto3" json:"segment,omitempty"`
}

func (x *SegmentAdded) Reset() {
	*x = SegmentAdded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentAdded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentAdded) ProtoMessage() {}

func (x *SegmentAdded) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentAdded.ProtoReflect.Descriptor instead.
func (*SegmentAdded) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_segments_proto_rawDescGZIP(), []int{2}
}

func (x *SegmentAdded) GetSegment() *Segment {
	if x != nil {
		return x.Segment
	}
	return nil
}

type SegmentRemoved struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Segment *Segment `protobuf:"bytes,1,opt,name=segment,proto3" json:"segment,omitempty"`
}

func (x *SegmentRemoved) Reset() {
	*x = SegmentRemoved{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentRemoved) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentRemoved) ProtoMessage() {}

func (x *SegmentRemoved) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentRemoved.ProtoReflect.Descriptor instead.
func (*SegmentRemoved) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_segments_proto_rawDescGZIP(), []int{3}
}

func (x *SegmentRemoved) GetSegment() *Segment {
	if x != nil {
		return x.Segment
	}
	return nil
}

type SegmentJournalEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Entry:
	//
	//	*SegmentJournalEntry_Started
	//	*SegmentJournalEntry_Added
	//	*SegmentJournalEntry_Removed
	Entry isSegmentJournalEntry_Entry `protobuf_oneof:"entry"`
}

func (x *SegmentJournalEntry) Reset() {
	*x = SegmentJournalEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentJournalEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentJournalEntry) ProtoMessage() {}

func (x *SegmentJournalEntry) ProtoReflect() protoreflect.Message {
	mi := &file_tukki_storage_segments_v1_segments_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentJournalEntry.ProtoReflect.Descriptor instead.
func (*SegmentJournalEntry) Descriptor() ([]byte, []int) {
	return file_tukki_storage_segments_v1_segments_proto_rawDescGZIP(), []int{4}
}

func (m *SegmentJournalEntry) GetEntry() isSegmentJournalEntry_Entry {
	if m != nil {
		return m.Entry
	}
	return nil
}

func (x *SegmentJournalEntry) GetStarted() *SegmentStarted {
	if x, ok := x.GetEntry().(*SegmentJournalEntry_Started); ok {
		return x.Started
	}
	return nil
}

func (x *SegmentJournalEntry) GetAdded() *SegmentAdded {
	if x, ok := x.GetEntry().(*SegmentJournalEntry_Added); ok {
		return x.Added
	}
	return nil
}

func (x *SegmentJournalEntry) GetRemoved() *SegmentRemoved {
	if x, ok := x.GetEntry().(*SegmentJournalEntry_Removed); ok {
		return x.Removed
	}
	return nil
}

type isSegmentJournalEntry_Entry interface {
	isSegmentJournalEntry_Entry()
}

type SegmentJournalEntry_Started struct {
	Started *SegmentStarted `protobuf:"bytes,1,opt,name=started,proto3,oneof"`
}

type SegmentJournalEntry_Added struct {
	Added *SegmentAdded `protobuf:"bytes,2,opt,name=added,proto3,oneof"`
}

type SegmentJournalEntry_Removed struct {
	Removed *SegmentRemoved `protobuf:"bytes,3,opt,name=removed,proto3,oneof"`
}

func (*SegmentJournalEntry_Started) isSegmentJournalEntry_Entry() {}

func (*SegmentJournalEntry_Added) isSegmentJournalEntry_Entry() {}

func (*SegmentJournalEntry_Removed) isSegmentJournalEntry_Entry() {}

var File_tukki_storage_segments_v1_segments_proto protoreflect.FileDescriptor

var file_tukki_storage_segments_v1_segments_proto_rawDesc = []byte{
	0x0a, 0x28, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f,
	0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x74, 0x75, 0x6b, 0x6b,
	0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x35, 0x0a, 0x07, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4b, 0x0a, 0x0e,
	0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x29,
	0x0a, 0x10, 0x6a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6a, 0x6f, 0x75, 0x72, 0x6e, 0x61,
	0x6c, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x0c, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x65, 0x64, 0x12, 0x3c, 0x0a, 0x07, 0x73, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x74, 0x75, 0x6b,
	0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07,
	0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x4e, 0x0a, 0x0e, 0x53, 0x65, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x3c, 0x0a, 0x07, 0x73, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x74, 0x75, 0x6b,
	0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07,
	0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0xed, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x67, 0x6d,
	0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x45, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x29, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x48, 0x00, 0x52, 0x07, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x3f, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x65, 0x64, 0x48, 0x00,
	0x52, 0x05, 0x61, 0x64, 0x64, 0x65, 0x64, 0x12, 0x45, 0x0a, 0x07, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x74, 0x75, 0x6b, 0x6b, 0x69,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x64, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x07,
	0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x42, 0xf9, 0x01, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e,
	0x74, 0x75, 0x6b, 0x6b, 0x69, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65,
	0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x53, 0x65, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68,
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
	file_tukki_storage_segments_v1_segments_proto_rawDescOnce sync.Once
	file_tukki_storage_segments_v1_segments_proto_rawDescData = file_tukki_storage_segments_v1_segments_proto_rawDesc
)

func file_tukki_storage_segments_v1_segments_proto_rawDescGZIP() []byte {
	file_tukki_storage_segments_v1_segments_proto_rawDescOnce.Do(func() {
		file_tukki_storage_segments_v1_segments_proto_rawDescData = protoimpl.X.CompressGZIP(file_tukki_storage_segments_v1_segments_proto_rawDescData)
	})
	return file_tukki_storage_segments_v1_segments_proto_rawDescData
}

var file_tukki_storage_segments_v1_segments_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_tukki_storage_segments_v1_segments_proto_goTypes = []interface{}{
	(*Segment)(nil),             // 0: tukki.storage.segments.v1.Segment
	(*SegmentStarted)(nil),      // 1: tukki.storage.segments.v1.SegmentStarted
	(*SegmentAdded)(nil),        // 2: tukki.storage.segments.v1.SegmentAdded
	(*SegmentRemoved)(nil),      // 3: tukki.storage.segments.v1.SegmentRemoved
	(*SegmentJournalEntry)(nil), // 4: tukki.storage.segments.v1.SegmentJournalEntry
}
var file_tukki_storage_segments_v1_segments_proto_depIdxs = []int32{
	0, // 0: tukki.storage.segments.v1.SegmentAdded.segment:type_name -> tukki.storage.segments.v1.Segment
	0, // 1: tukki.storage.segments.v1.SegmentRemoved.segment:type_name -> tukki.storage.segments.v1.Segment
	1, // 2: tukki.storage.segments.v1.SegmentJournalEntry.started:type_name -> tukki.storage.segments.v1.SegmentStarted
	2, // 3: tukki.storage.segments.v1.SegmentJournalEntry.added:type_name -> tukki.storage.segments.v1.SegmentAdded
	3, // 4: tukki.storage.segments.v1.SegmentJournalEntry.removed:type_name -> tukki.storage.segments.v1.SegmentRemoved
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_tukki_storage_segments_v1_segments_proto_init() }
func file_tukki_storage_segments_v1_segments_proto_init() {
	if File_tukki_storage_segments_v1_segments_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tukki_storage_segments_v1_segments_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Segment); i {
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
		file_tukki_storage_segments_v1_segments_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentStarted); i {
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
		file_tukki_storage_segments_v1_segments_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentAdded); i {
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
		file_tukki_storage_segments_v1_segments_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentRemoved); i {
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
		file_tukki_storage_segments_v1_segments_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentJournalEntry); i {
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
	file_tukki_storage_segments_v1_segments_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*SegmentJournalEntry_Started)(nil),
		(*SegmentJournalEntry_Added)(nil),
		(*SegmentJournalEntry_Removed)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tukki_storage_segments_v1_segments_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tukki_storage_segments_v1_segments_proto_goTypes,
		DependencyIndexes: file_tukki_storage_segments_v1_segments_proto_depIdxs,
		MessageInfos:      file_tukki_storage_segments_v1_segments_proto_msgTypes,
	}.Build()
	File_tukki_storage_segments_v1_segments_proto = out.File
	file_tukki_storage_segments_v1_segments_proto_rawDesc = nil
	file_tukki_storage_segments_v1_segments_proto_goTypes = nil
	file_tukki_storage_segments_v1_segments_proto_depIdxs = nil
}
