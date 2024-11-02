package db

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/segments"
)

var logger = log.New(log.Writer(), "compactor: ", log.LstdFlags)

const T1TargetSegmentSize = 160 * 1024 * 1024
const T2TargetSegmentSize = T1TargetSegmentSize * 10

func (db *Database) compactor() {
	for {
		select {
		case <-time.After(30 * time.Second):
			db.Compact()
		case <-db.compactorStop:
			logger.Printf("compactor stopped")
			return
		}
	}
}

func getSegmentFileSize(dbDir string, segment segments.SegmentMetadata) (int64, error) {
	return files.FileSize(dbDir, segment.SegmentFile)
}

func (db *Database) getSegmentsReverseSorted() []segments.SegmentMetadata {
	segments := db.getSegmentsSorted()
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Id < segments[j].Id
	})
	return segments
}

func bytesToMegaBytes(bytes int64) string {
	return fmt.Sprintf("%d MiB", bytes/1024/1024)
}

type CompactionSegment struct {
	Id   segments.SegmentId
	Size int64
}

type CompactionCriteria func(size int64) bool

func DecideMergedSegments(criteria CompactionCriteria, segmentList []CompactionSegment) []segments.SegmentId {
	for len(segmentList) >= 4 {
		size := segmentList[0].Size
		if !criteria(size) {
			segmentList = segmentList[1:]
			continue
		}

		nextSmallSegments := []segments.SegmentId{segmentList[0].Id}
		for _, sg := range segmentList[1:] {
			if !criteria(sg.Size) {
				break
			}

			nextSmallSegments = append(nextSmallSegments, sg.Id)
		}

		if len(nextSmallSegments) >= 4 {
			// latest first
			sort.Slice(nextSmallSegments, func(i, j int) bool {
				return nextSmallSegments[i] > nextSmallSegments[j]
			})

			return nextSmallSegments
		}

		break
	}

	return nil
}

type CompactionOperation struct {
	Criteria   CompactionCriteria
	TargetSize int64
}

var compactionCriterias = []CompactionOperation{
	{Criteria: t1Criteria, TargetSize: T1TargetSegmentSize},
	{Criteria: t2Criteria, TargetSize: T2TargetSegmentSize},
}

func t1Criteria(size int64) bool {
	return size < T1TargetSegmentSize
}

func t2Criteria(size int64) bool {
	return size >= T1TargetSegmentSize && size < T2TargetSegmentSize
}

func (db *Database) findSegmentsToMerge(criteria CompactionCriteria) ([]segments.SegmentId, error) {
	sgs := db.getSegmentsReverseSorted()
	compactionInfo := make([]CompactionSegment, 0, len(sgs))
	for _, sg := range sgs {
		size, err := getSegmentFileSize(db.dbDir, sg)
		if err != nil {
			return nil, fmt.Errorf("failed to get segment file size: %w", err)
		}

		compactionInfo = append(compactionInfo, CompactionSegment{
			Id:   sg.Id,
			Size: size,
		})
	}

	return DecideMergedSegments(criteria, compactionInfo), nil
}

func (db *Database) Compact() error {
	for _, criteria := range compactionCriterias {
		segmentIds, err := db.findSegmentsToMerge(criteria.Criteria)
		if err != nil {
			logger.Printf("error finding segments to merge: %v", err)
			continue
		}
		if len(segmentIds) == 0 {
			logger.Printf("no segments to merge")
			continue
		}

		err = db.CompactSegments(uint64(criteria.TargetSize), segmentIds...)
		if err != nil {
			logger.Printf("error merging segments: %v", err)
			return err
		}
	}

	return nil
}
