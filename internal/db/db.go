package db

import (
	"log"

	"github.com/jukeks/tukki/internal/segments"
)

type Database struct {
	ongoing        *segments.LiveSegment
	segmentManager *segments.SegmentManager
}

func OpenDatabase(dbDir string) *Database {
	segmentsManager, err := segments.OpenDatabase(dbDir)
	if err != nil {
		log.Fatalf("failed to open segments manager: %v", err)
	}

	ongoing := segmentsManager.GetOnGoingSegment()
	return &Database{
		ongoing:        ongoing,
		segmentManager: segmentsManager,
	}
}

func (d *Database) Set(key, value string) error {
	return d.ongoing.Set(key, value)
}

func (d *Database) Get(key string) (string, error) {
	value, err := d.ongoing.Get(key)
	if err == nil {
		return value, nil
	}

	if err != segments.ErrKeyNotFound {
		val, err := d.segmentManager.Get(key)
		if err != nil {
			return "", err
		}
		return val, nil
	}

	return "", err
}

func (d *Database) Delete(key string) error {
	err := d.ongoing.Delete(key)
	if err != nil {
		return err
	}

	d.ongoing.Delete(key)
	return nil
}

func (d *Database) Close() error {
	err := d.ongoing.Close()
	if err != nil {
		return err
	}

	return d.segmentManager.Close()
}
