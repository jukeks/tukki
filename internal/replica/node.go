package replica

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"

	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/internal/grpc/kv"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
	maxMdbSize          = 64 * 1024 * 1024 * 1024
)

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Min   string `json:"min,omitempty"`
	Max   string `json:"max,omitempty"`
}

type Peer struct {
	Id   string
	Addr string
}

type Node struct {
	RaftDir  string
	RaftBind string
	inmem    bool
	peers    []Peer

	mu    sync.RWMutex
	db    *db.Database
	dbDir string

	raft *raft.Raft // The consensus mechanism

	logger *log.Logger
}

// New returns a new Store.
func New(inmem bool, peers []Peer, db *db.Database, dbDir, raftDir, raftBind string) *Node {
	return &Node{
		inmem:    inmem,
		peers:    peers,
		logger:   log.New(os.Stderr, "[store] ", log.LstdFlags),
		db:       db,
		dbDir:    dbDir,
		RaftDir:  raftDir,
		RaftBind: raftBind,
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (n *Node) Open(localID string, enableSingle bool, peers []Peer) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", n.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(n.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(n.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	if n.inmem {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		raftDb, err := db.OpenDatabase(n.RaftDir)
		if err != nil {
			return fmt.Errorf("failed to open raft database, err: %w", err)
		}
		raftukki := &Raftukki{db: raftDb}
		stableStore = raftukki
		logStore = raftukki
	}

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(n), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %w", err)
	}
	n.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		for _, peer := range peers {
			configuration.Servers = append(configuration.Servers, raft.Server{
				ID:      raft.ServerID(peer.Id),
				Address: raft.ServerAddress(peer.Addr),
			})
		}
		ra.BootstrapCluster(configuration)
	}

	return nil
}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (n *Node) Join(nodeID, addr string) error {
	n.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := n.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		n.logger.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				n.logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := n.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	f := n.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	n.logger.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}

// Get returns the value for the given key.
func (n *Node) Get(key string) (string, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.db.Get(key)
}

// GetRange implements kv.DB.
func (n *Node) GetRange(min string, max string) ([]kv.Pair, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	dbPairs, err := n.db.GetRange(min, max)
	if err != nil {
		return nil, err
	}

	var pairs []kv.Pair
	for _, dbPair := range dbPairs {
		pairs = append(pairs, kv.Pair{
			Key:   dbPair.Key,
			Value: dbPair.Value,
		})
	}

	return pairs, nil
}

// Set sets the value for the given key.
func (n *Node) Set(key, value string) error {
	if n.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := n.raft.Apply(b, raftTimeout)
	if f.Error() != nil {
		return f.Error()
	}
	err, _ = f.Response().(error)
	return err
}

// Delete deletes the given key.
func (n *Node) Delete(key string) error {
	if n.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &command{
		Op:  "delete",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := n.raft.Apply(b, raftTimeout)
	if f.Error() != nil {
		return f.Error()
	}
	err, _ = f.Response().(error)
	return err
}

// DeleteRange implements kv.DB.
func (n *Node) DeleteRange(min string, max string) (uint64, error) {
	if n.raft.State() != raft.Leader {
		return 0, fmt.Errorf("not leader")
	}

	c := &command{
		Op:  "deleteRange",
		Min: min,
		Max: max,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return 0, err
	}

	f := n.raft.Apply(b, raftTimeout)
	if f.Error() != nil {
		return 0, f.Error()
	}

	if f.Error() != nil {
		return 0, f.Error()
	}
	err, _ = f.Response().(error)
	return 0, err
}
