package segments

import "github.com/jukeks/tukki/internal/storage/files"

type SegmentId uint64

type SegmentMetadata struct {
	Id          SegmentId
	SegmentFile files.Filename
	MembersFile files.Filename
	IndexFile   files.Filename
}
