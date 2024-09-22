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
	raftboltdb "github.com/hashicorp/raft-mdb"

	"github.com/jukeks/tukki/internal/db"
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
		mdbStore, err := raftboltdb.NewMDBStoreWithSize(n.RaftDir, maxMdbSize)
		if err != nil {
			return fmt.Errorf("failed to create mdb store: %w", err)
		}
		logStore = mdbStore
		stableStore = mdbStore
	}

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(n), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
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
	return f.Error()
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
	return f.Error()
}
