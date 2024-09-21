package sstable

import (
	"context"
	"testing"

	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/testutil"
	"google.golang.org/grpc"

	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/sstable/v1"
)

func TestSSTable(t *testing.T) {
	dbDir := t.TempDir()
	db, err := db.OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err := db.Set("key", "value"); err != nil {
		t.Fatalf("failed to set key: %v", err)
	}
	if _, err := db.SealCurrentSegment(); err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	server := NewSstableServer(db)
	conn, cleanup, err := testutil.RunServicer(func(s *grpc.Server) {
		sstablev1.RegisterSstableServiceServer(s, server)
	})
	if err != nil {
		t.Fatalf("failed to run servicer: %v", err)
	}
	defer cleanup()

	client := sstablev1.NewSstableServiceClient(conn)
	ctx := context.Background()

	stream, err := client.GetSstable(ctx, &sstablev1.GetSstableRequest{
		Id: 0,
	})
	if err != nil {
		t.Fatalf("failed to get sstable: %v", err)
	}

	records := []*sstablev1.SSTableRecord{}
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		records = append(records, resp.Record)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
}
