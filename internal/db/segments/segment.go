package segments

import "github.com/jukeks/tukki/internal/storage"

type SegmentId uint64

type SegmentMetadata struct {
	Id          SegmentId
	SegmentFile storage.Filename
	MembersFile storage.Filename
	IndexFile   storage.Filename
}
