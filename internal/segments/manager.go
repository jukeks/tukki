package segments

type SegmentManager struct {
	segmentJournal   *SegmentJournal
	operationJournal *SegmentOperationJournal

	segments   map[SegmentId]Segment
	operations map[OperationId]SegmentOperation
}

func NewSegmentManager(dbDir string) (*SegmentManager, error) {
	segmentJournal, segments, err := OpenSegmentJournal(dbDir)
	if err != nil {
		return nil, err
	}

	operationJournal, operations, err := OpenSegmentOperationJournal(dbDir)
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
