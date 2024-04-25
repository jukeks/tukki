package segments

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type SegmentJournal struct {
	journal *journal.Journal
}

func OpenSegmentJournal(dbDir string) (*SegmentJournal, map[uint64]Segment, error) {
	var segments map[uint64]Segment
	handle := func(r *journal.JournalReader) error {
		var err error
		segments, err = readJournal(r)
		return err
	}

	j, err := journal.OpenJournal(dbDir, "segments.journal", handle)
	if err != nil {
		return nil, nil, err
	}

	return &SegmentJournal{j}, segments, nil
}

func readJournal(r *journal.JournalReader) (map[uint64]Segment, error) {
	segmentMap := make(map[uint64]Segment)
	for {
		journalEntry := &segmentsv1.SegmentJournalEntry{}
		err := r.Read(journalEntry)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		added := journalEntry.GetAdded()
		if added != nil {
			segment := Segment{
				Id:       added.Segment.Id,
				Filename: added.Segment.Filename,
			}
			segmentMap[segment.Id] = segment
		}
		removed := journalEntry.GetRemoved()
		if removed != nil {
			delete(segmentMap, removed.Segment.Id)
		}
	}

	return segmentMap, nil
}
