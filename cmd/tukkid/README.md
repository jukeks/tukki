# Multi node tukki

tukkid is fault tolerant and distributed tukki. It has automatic leader failover and single leader replication throught Raft.

Currently only voting raft members can be specified with the configuration so the tukkid should be run with either 3 or 5 nodes. Raft requires `n/2 + 1` members to vote, so with 3 nodes it can tolerate losing 1 node and with 5 nodes it can tolerate losing 2 nodes.

## Configuration

The most important configuration options are:
```bash
$ ./bin/tukkid --help
Usage of ./bin/tukkid:
  -node-id string
    	The node ID (default "node1")
  -db-dir string
    	The directory to store the database (default "./tukki-db")
  -port int
    	The server port (default 50051)
  -raft-port int
    	The Raft server port (default 50000)
  -raft-peers string
    	The Raft peers. Only relevant when initializing.
  -sstable-peers string
    	The SSTable peers. Used to sync segments on restore.
  -init
    	Initialize the database
```

When initializing a tukki cluster, a node should be started with:
```bash
./bin/tukkid \
    -node-id 1 \
    -db-dir ./node-1 \
    -port 50010 \
    -raft-port 50011 \
    -raft-peers 2@localhost:50021,3@localhost:50031 \
    -sstable-peers 2@localhost:50020,3@localhost:50030 \
    -init
```

And rest two nodes with

```bash
./bin/tukkid \
    -node-id 2 \
    -db-dir ./node-2 \
    -port 50020 \
    -raft-port 50021 \
    -sstable-peers 1@localhost:50010,3@localhost:50030
```
and
```bash
./bin/tukkid \
    -node-id 3 \
    -db-dir ./node-3 \
    -port 50030 \
    -raft-port 50031 \
    -sstable-peers 2@localhost:50020,1@localhost:50010
```
