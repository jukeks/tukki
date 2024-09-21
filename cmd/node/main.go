package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jukeks/tukki/cmd/node/replica"
	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/internal/grpc/kv"
	"github.com/jukeks/tukki/internal/grpc/sstable"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/sstable/v1"
	"google.golang.org/grpc"
)

func defaultDatabaseDir() string {
	return "./tukki-db"
}

var (
	port            = flag.Int("port", 50051, "The server port")
	raftPort        = flag.Int("raft-port", 50000, "The Raft server port")
	nodeId          = flag.String("node-id", "node1", "The node ID")
	dbDir           = flag.String("db-dir", defaultDatabaseDir(), "The directory to store the database")
	raftPeerList    = flag.String("raft-peers", "", "The Raft peers")
	sstablePeerList = flag.String("sstable-peers", "", "The SSTable peers")
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

	sstablePeers, err := parsePeers(*sstablePeerList)
	if err != nil {
		log.Fatalf("failed to sstable parse peers: %v", err)
	}
	/*
		var raftPeers []node.Peer
		if *raftPeerList != "" {
			raftPeers, err = parsePeers(*raftPeerList)
			if err != nil {
				log.Fatalf("failed to raft parse peers: %v", err)
			}
		}*/

	n := replica.New(false, sstablePeers, db, *dbDir, *dbDir, fmt.Sprintf("localhost:%d", *raftPort))
	if err := n.Open(*raftPeerList != "", *nodeId); err != nil {
		log.Fatalf("failed to open node: %v", err)
	}

	go func() {
		/*
			time.Sleep(5 * time.Second)
			if *raftPeerList != "" {
				for _, peer := range raftPeers {
					if err := n.Join(peer.Id, peer.Addr); err != nil {
						log.Fatalf("failed to join peer: %v", err)
					}
				}
			}*/
	}()

	kvServer := kv.NewKVServer(n)
	sstableServer := sstable.NewSstableServer(db)

	grpcServer := grpc.NewServer()
	kvv1.RegisterKvServiceServer(grpcServer, kvServer)
	sstablev1.RegisterSstableServiceServer(grpcServer, sstableServer)

	grpcServer.Serve(ls)
}
