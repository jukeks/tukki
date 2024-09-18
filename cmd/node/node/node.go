package node

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"

	"github.com/jukeks/tukki/internal/db"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Node struct {
	RaftDir  string
	RaftBind string
	inmem    bool

	mu sync.RWMutex
	db *db.Database

	raft *raft.Raft // The consensus mechanism

	logger *log.Logger
}

// New returns a new Store.
func New(inmem bool, db *db.Database, raftDir string, raftBind string) *Node {
	return &Node{
		inmem:    inmem,
		logger:   log.New(os.Stderr, "[store] ", log.LstdFlags),
		db:       db,
		RaftDir:  raftDir,
		RaftBind: raftBind,
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (n *Node) Open(enableSingle bool, localID string) error {
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
		boltDB, err := raftboltdb.New(raftboltdb.Options{
			Path: filepath.Join(n.RaftDir, "raft.db"),
		})
		if err != nil {
			return fmt.Errorf("new bbolt store: %s", err)
		}
		logStore = boltDB
		stableStore = boltDB
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
		ra.BootstrapCluster(configuration)
	}

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

type fsm Node

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// TODO
	return &fsmSnapshot{}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	// TODO
	return nil
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.db.Set(key, value); err != nil {
		f.logger.Printf("failed to apply set key %s: %s", key, err)
		return err
	}
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.db.Delete(key); err != nil {
		f.logger.Printf("failed to apply delete key %s: %s", key, err)
		return err
	}
	return nil
}

type fsmSnapshot struct {
	store map[string]string
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
