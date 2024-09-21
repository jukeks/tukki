package segments

import segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"

type OperationId uint64

type SegmentOperation interface {
	Id() OperationId
	StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry
	Execute() error
	CompletedJournalEntry() *segmentsv1.SegmentOperationJournalEntry
}
