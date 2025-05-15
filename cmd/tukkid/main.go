package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/internal/grpc/kv"
	"github.com/jukeks/tukki/internal/grpc/sstable"
	"github.com/jukeks/tukki/internal/replica"
	"github.com/jukeks/tukki/internal/storage/journal"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/sstable/v1"
	"google.golang.org/grpc"
)

func defaultDatabaseDir() string {
	return "./tukki-db"
}

var (
	port     = flag.Int("port", 50051, "The server port")
	raftPort = flag.Int("raft-port", 50000, "The Raft server port")
	nodeId   = flag.String("node-id", "node1", "The node ID")
	dbDir    = flag.String("db-dir", defaultDatabaseDir(),
		"The directory to store the database")
	raftPeerList = flag.String("raft-peers", "",
		"The Raft peers. Only relevant when initializing")
	sstablePeerList = flag.String("sstable-peers", "",
		"The SSTable peers. Used to sync missing segments. Should have the server port.")
	inititialize = flag.Bool("init", false, "Initialize the database")
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parsePeers(peers string) ([]replica.Peer, error) {
	peerList := strings.Split(peers, ",")
	result := make([]replica.Peer, 0, len(peerList))
	for _, peer := range peerList {
		components := strings.Split(peer, "@")
		if len(components) != 2 {
			return nil, fmt.Errorf("invalid peer: %s", peer)
		}
		result = append(result, replica.Peer{Id: components[0], Addr: components[1]})
	}
	return result, nil
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("failed to create cpuprofile: %v", err)
		}
		pprof.StartCPUProfile(f)
	}

	defer pprof.StopCPUProfile()

	err := os.MkdirAll(*dbDir, 0755)
	if err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	config := db.GetDefaultConfig()
	// In memory journal is good as raft will replay logs on startup
	config.JournalMode = journal.WriteModeInMemory
	// Raft replays logs on startup, db needs to be aware
	config.ReplicaMode = true

	db, err := db.OpenDatabaseWithConfig(*dbDir, config)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sstablePeers, err := parsePeers(*sstablePeerList)
	if err != nil {
		log.Fatalf("failed to sstable parse peers: %v", err)
	}

	var raftPeers []replica.Peer
	if *raftPeerList != "" {
		raftPeers, err = parsePeers(*raftPeerList)
		if err != nil {
			log.Fatalf("failed to raft parse peers: %v", err)
		}
	}

	n := replica.New(false, sstablePeers, db, *dbDir, *dbDir+"/raft", fmt.Sprintf("localhost:%d", *raftPort))
	if err := n.Open(*nodeId, *inititialize, raftPeers); err != nil {
		log.Fatalf("failed to open node: %v", err)
	}

	kvServer := kv.NewKVServer(n)
	sstableServer := sstable.NewSstableServer(db)

	grpcServer := grpc.NewServer()
	kvv1.RegisterKvServiceServer(grpcServer, kvServer)
	sstablev1.RegisterSstableServiceServer(grpcServer, sstableServer)

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			<-sigchnl
			grpcServer.GracefulStop()
			break
		}
	}()

	if err := grpcServer.Serve(ls); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
