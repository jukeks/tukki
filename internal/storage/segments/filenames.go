package segments

import (
	"fmt"

	"github.com/jukeks/tukki/internal/storage/files"
)

func GetWalFilename(id SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("wal-%d.journal", id))
}

func GetSegmentFilename(id SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("segment-%d", id))
}

func GetMergedSegmentFilename(a, b SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("segment-%d-%d", a, b))
}

func GetCompactedSegmentFilename(a SegmentId, sequence OperationId) files.Filename {
	return files.Filename(fmt.Sprintf("segment-%d-%d", a, sequence))
}

func GetMembersFilename(id SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("members-%d", id))
}

func GetMergedMembersFilename(a, b SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("members-%d-%d", a, b))
}

func GetCompactedMembersFilename(a SegmentId, sequence OperationId) files.Filename {
	return files.Filename(fmt.Sprintf("members-%d-%d", a, sequence))
}

func GetIndexFilename(id SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("index-%d", id))
}

func GetMergedIndexFilename(a, b SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("index-%d-%d", a, b))
}

func GetCompactedIndexFilename(a SegmentId, sequence OperationId) files.Filename {
	return files.Filename(fmt.Sprintf("index-%d-%d", a, sequence))
}
