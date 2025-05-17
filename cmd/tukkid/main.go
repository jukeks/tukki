package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/goccy/go-yaml"
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
	config     = flag.String("config", "", "config file path")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Config struct {
	NodeID     string `yaml:"node-id"`
	DBDir      string `yaml:"db-dir"`
	PublicPort int    `yaml:"public-port"`
	Cluster    struct {
		Port  int    `yaml:"port"`
		Init  bool   `yaml:"init"`
		Peers []Peer `yaml:"peers"`
	} `yaml:"cluster"`
}

type Peer struct {
	ID          string `yaml:"id"`
	RaftAddr    string `yaml:"raft-addr"`
	SSTableAddr string `yaml:"sstable-addr"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getRaftPeers(peers []Peer) []replica.Peer {
	sstablePeers := make([]replica.Peer, len(peers))
	for i, peer := range peers {
		sstablePeers[i] = replica.Peer{
			Id:   peer.ID,
			Addr: peer.RaftAddr,
		}
	}

	return sstablePeers
}

func getSSTablePeers(peers []Peer) []replica.Peer {
	sstablePeers := make([]replica.Peer, len(peers))
	for i, peer := range peers {
		sstablePeers[i] = replica.Peer{
			Id:   peer.ID,
			Addr: peer.SSTableAddr,
		}
	}

	return sstablePeers
}

func main() {
	flag.Parse()
	cfg, err := loadConfig(*config)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("failed to create cpuprofile: %v", err)
		}
		pprof.StartCPUProfile(f)
	}

	defer pprof.StopCPUProfile()

	err = os.MkdirAll(cfg.DBDir, 0755)
	if err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	config := db.GetDefaultConfig()
	// In memory journal is good as raft will replay logs on startup
	config.JournalMode = journal.WriteModeInMemory
	// Raft replays logs on startup, db needs to be aware
	config.ReplicaMode = true

	db, err := db.OpenDatabaseWithConfig(cfg.DBDir, config)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.PublicPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sstablePeers := getSSTablePeers(cfg.Cluster.Peers)
	raftPeers := getRaftPeers(cfg.Cluster.Peers)

	n := replica.New(false, sstablePeers, db, cfg.DBDir, cfg.DBDir+"/raft", fmt.Sprintf("localhost:%d", cfg.Cluster.Port))
	if err := n.Open(cfg.NodeID, cfg.Cluster.Init, raftPeers); err != nil {
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
