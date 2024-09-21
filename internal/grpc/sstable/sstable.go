package sstable

import (
	"io"

	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/internal/storage/segments"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/sstable/v1"
)

type SSTableServer struct {
	sstablev1.UnimplementedSstableServiceServer
	db *db.Database
}

func NewSstableServer(db *db.Database) *SSTableServer {
	return &SSTableServer{
		db: db,
	}
}

func (s *SSTableServer) GetSstable(req *sstablev1.GetSstableRequest, stream sstablev1.SstableService_GetSstableServer) error {
	reader, cleanup, err := s.db.GetSSTableReader(segments.SegmentId(req.Id))
	if err != nil {
		return err
	}
	defer cleanup()

	for {
		entry, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		msg := &sstablev1.GetSstableResponse{
			Record: &sstablev1.SSTableRecord{
				Key:     entry.Key,
				Value:   entry.Value,
				Deleted: entry.Deleted,
			},
		}

		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
