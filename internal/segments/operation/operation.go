package operation

import segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"

type Segment struct {
	Id       uint64
	Filename string
}

type SegmentOperation interface {
	Id() uint64
	StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry
	Execute() error
	CompletedJournalEntry() *segmentsv1.SegmentOperationJournalEntry
}
