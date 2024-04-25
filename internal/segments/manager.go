package segments

import "github.com/jukeks/tukki/internal/segments/operation"

type Segment struct {
	Id       uint64
	Filename string
}

type SegmentManager struct {
	segments   []Segment
	operations []operation.SegmentOperation
}
