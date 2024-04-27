package segments

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type SegmentJournal struct {
	journal *journal.Journal
}

func OpenSegmentJournal(dbDir string) (*SegmentJournal, *CurrentSegments, error) {
	var currentSegments *CurrentSegments
	handle := func(r *journal.JournalReader) error {
		var err error
		currentSegments, err = readSegmentJournal(r)
		return err
	}

	j, err := journal.OpenJournal(dbDir, "segments.journal", handle)
	if err != nil {
		return nil, nil, err
	}

	return &SegmentJournal{j}, currentSegments, nil
}

type OngoingSegment struct {
	Id              SegmentId
	JournalFilename storage.Filename
}

type CurrentSegments struct {
	Ongoing  OngoingSegment
	Segments map[SegmentId]Segment
}

func readSegmentJournal(r *journal.JournalReader) (*CurrentSegments, error) {
	var ongoing OngoingSegment
	segmentMap := make(map[SegmentId]Segment)
	for {
		journalEntry := &segmentsv1.SegmentJournalEntry{}
		err := r.Read(journalEntry)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch journalEntry.Entry.(type) {
		case *segmentsv1.SegmentJournalEntry_Started:
			started := journalEntry.GetStarted()
			ongoing = OngoingSegment{
				Id:              SegmentId(started.Id),
				JournalFilename: storage.Filename(started.JournalFilename),
			}
		case *segmentsv1.SegmentJournalEntry_Added:
			added := journalEntry.GetAdded()
			segment := Segment{
				Id:       SegmentId(added.Segment.Id),
				Filename: storage.Filename(added.Segment.Filename),
			}
			segmentMap[segment.Id] = segment

		case *segmentsv1.SegmentJournalEntry_Removed:
			removed := journalEntry.GetRemoved()
			delete(segmentMap, SegmentId(removed.Segment.Id))
		}
	}

	return &CurrentSegments{
		Ongoing:  ongoing,
		Segments: segmentMap,
	}, nil
}

func (sj *SegmentJournal) StartSegment(id SegmentId, journalFilename storage.Filename) error {
	entry := &segmentsv1.SegmentJournalEntry{
		Entry: &segmentsv1.SegmentJournalEntry_Started{
			Started: &segmentsv1.SegmentStarted{
				Id:              uint64(id),
				JournalFilename: string(journalFilename),
			},
		},
	}

	return sj.journal.Writer.Write(entry)
}

func (sj *SegmentJournal) AddSegment(segment Segment) error {
	entry := &segmentsv1.SegmentJournalEntry{
		Entry: &segmentsv1.SegmentJournalEntry_Added{
			Added: &segmentsv1.SegmentAdded{
				Segment: &segmentsv1.Segment{
					Id:       uint64(segment.Id),
					Filename: string(segment.Filename),
				},
			},
		},
	}

	return sj.journal.Writer.Write(entry)
}

func (sj *SegmentJournal) RemoveSegment(id SegmentId) error {
	entry := &segmentsv1.SegmentJournalEntry{
		Entry: &segmentsv1.SegmentJournalEntry_Removed{
			Removed: &segmentsv1.SegmentRemoved{
				Segment: &segmentsv1.Segment{
					Id: uint64(id),
				},
			},
		},
	}

	return sj.journal.Writer.Write(entry)
}

func (sj *SegmentJournal) Close() error {
	return sj.journal.Close()
}
