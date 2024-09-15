package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/jukeks/tukki/internal/db"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	"google.golang.org/grpc"
)

func defaultDatabaseDir() string {
	return "./tukki-db"
}

var (
	port  = flag.Int("port", 50051, "The server port")
	dbDir = flag.String("db-dir", defaultDatabaseDir(), "The directory to store the database")
)

func main() {
	flag.Parse()

	err := os.MkdirAll(*dbDir, 0755)
	if err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	db, err := db.OpenDatabase(*dbDir)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	kvServer := NewKvServer(db)

	grpcServer := grpc.NewServer()
	kvv1.RegisterKvServiceServer(grpcServer, kvServer)
	grpcServer.Serve(ls)
}

type kvServer struct {
	kvv1.UnimplementedKvServiceServer
	lock sync.RWMutex
	db   *db.Database
}

func NewKvServer(db *db.Database) *kvServer {
	return &kvServer{db: db}
}

func (s *kvServer) Query(ctx context.Context, req *kvv1.QueryRequest) (*kvv1.QueryResponse, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, err := s.db.Get(req.Key)
	if err != nil {
		return nil, err
	}

	return &kvv1.QueryResponse{Value: &kvv1.QueryResponse_Pair{Pair: &kvv1.KvPair{
		Key:   req.Key,
		Value: value,
	}}}, nil
}

func (s *kvServer) Set(ctx context.Context, req *kvv1.SetRequest) (*kvv1.SetResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := s.db.Set(req.Pair.Key, req.Pair.Value)
	if err != nil {
		return nil, err
	}

	return &kvv1.SetResponse{}, nil
}

func (s *kvServer) Delete(ctx context.Context, req *kvv1.DeleteRequest) (*kvv1.DeleteResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := s.db.Delete(req.Key)
	if err != nil {
		return nil, err
	}

	return &kvv1.DeleteResponse{}, nil
}
