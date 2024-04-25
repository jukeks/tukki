package segments

import (
	"io"
)

type SegmentJournal struct {
	writer io.Writer
}

func NewSegmentJournal(writer io.Writer) *SegmentJournal {
	return &SegmentJournal{
		writer: writer,
	}
}
