package segments

import "github.com/jukeks/tukki/internal/segments/operation"

type Segment struct {
	Id       uint64
	Filename string
}

type SegmentManager struct {
	segmentJournal   *SegmentJournal
	operationJournal *operation.SegmentOperationJournal

	segments   map[uint64]Segment
	operations map[uint64]operation.SegmentOperation
}

func NewSegmentManager(dbDir string) (*SegmentManager, error) {
	segmentJournal, segments, err := OpenSegmentJournal(dbDir)
	if err != nil {
		return nil, err
	}

	operationJournal, operations, err := operation.OpenSegmentOperationJournal(dbDir)
	if err != nil {
		return nil, err
	}

	return &SegmentManager{
		segmentJournal:   segmentJournal,
		operationJournal: operationJournal,
		segments:         segments,
		operations:       operations,
	}, nil
}
