package db

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/segments"
)

var logger = log.New(log.Writer(), "compactor: ", log.LstdFlags)

func (db *Database) compactor() {
	for {
		select {
		case <-time.After(30 * time.Second):
			db.compact()
		case <-db.compactorStop:
			logger.Printf("compactor stopped")
			return
		}
	}
}

func getSegmentFileSize(dbDir string, segment segments.SegmentMetadata) (int64, error) {
	return files.FileSize(dbDir, segment.SegmentFile)
}

func getOptimalSegmentCountFor(totalSize int64) int {
	return int(math.Log2(float64(totalSize))) + 1
}

func getTotalSize(dbDir string, segments []segments.SegmentMetadata) (int64, error) {
	totalSize := int64(0)
	for _, segment := range segments {
		size, err := getSegmentFileSize(dbDir, segment)
		if err != nil {
			return 0, err
		}
		totalSize += size
	}
	return totalSize, nil
}

func (db *Database) findSegmentsToMerge() ([]segments.SegmentId, error) {
	segmentsToConsider := db.getSegmentsSorted()
	segmentCount := len(segmentsToConsider)

	if segmentCount < 2 {
		logger.Printf("not enough segments to merge")
		return nil, nil
	}

	// if the segments are too small, don't merge
	totalSize, err := getTotalSize(db.dbDir, segmentsToConsider)
	if err != nil {
		return nil, fmt.Errorf("failed to get total size of segments: %w", err)
	}

	optimalCount := getOptimalSegmentCountFor(totalSize)
	if optimalCount > segmentCount {
		logger.Printf("no need to merge yet: %d segments, total size %d MiB",
			segmentCount, totalSize/1024/1024)
		return nil, nil
	}

	// consider last two
	toMerge := segmentsToConsider[len(segmentsToConsider)-2:]

	ids := []segments.SegmentId{toMerge[0].Id, toMerge[1].Id}
	// sort so that we return the smallest id first
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	return ids, nil
}

func (db *Database) compact() {
}
