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

const TargetSegmentSize = 160 * 1024 * 1024

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

func DecideMergedSegments(targetSize int64, segmentList []CompactionSegment) []segments.SegmentId {
	for len(segmentList) >= 4 {
		size := segmentList[0].Size
		if size > targetSize {
			log.Printf("segment %d is already larger than target size: %s",
				segmentList[0].Id, bytesToMegaBytes(size))
			segmentList = segmentList[1:]
			log.Printf("segmentList: %v", segmentList)
			continue
		}

		nextSmallSegments := []segments.SegmentId{segmentList[0].Id}
		for _, sg := range segmentList[1:] {
			if sg.Size > targetSize {
				log.Printf("segment %d too large: %s", sg.Id, bytesToMegaBytes(sg.Size))
				break
			}

			log.Printf("segment %d is small enough: %s", sg.Id, bytesToMegaBytes(sg.Size))
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

func (db *Database) findSegmentsToMerge(targetSize int64) ([]segments.SegmentId, error) {
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
	return DecideMergedSegments(targetSize, compactionInfo), nil
}

func (db *Database) Compact() {
	segmentIds, err := db.findSegmentsToMerge(TargetSegmentSize)
	if err != nil {
		logger.Printf("error finding segments to merge: %v", err)
		return
	}
	if len(segmentIds) == 0 {
		logger.Printf("no segments to merge")
		return
	}

	logger.Printf("merging segments: %v", segmentIds)
	err = db.CompactSegments(TargetSegmentSize, segmentIds...)
	if err != nil {
		logger.Printf("error merging segments: %v", err)
	}
}
