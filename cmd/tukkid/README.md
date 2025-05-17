# tukkid or Multi node tukki

tukkid is fault tolerant and distributed tukki. It has automatic leader failover and single leader replication through Raft.

Currently only voting raft members can be specified with the configuration so the tukkid should be run with either 3 or 5 nodes. Raft requires `n/2 + 1` members to vote, so with 3 nodes it can tolerate losing 1 node and with 5 nodes it can tolerate losing 2 nodes.

## Configuration

```bash
$ ./bin/tukkid --help
Usage of ./bin/tukkid:
  -config string
        config file path
```

When initializing a tukki cluster, a node should be started with `cluster.init: true` like in the example config:

```yaml
node-id: "1"
db-dir: "./node-1"
public-port: 50010
cluster:
  port: 50011
  peers: 
    - id: "2"
      raft-addr: "localhost:50021"
      sstable-addr: "localhost:50020"
    - id: "3"
      raft-addr: "localhost:50031"
      sstable-addr: "localhost:50030"
  init: true
```

and run with
```bash
./bin/tukkid -config ./example/node1.yaml
```
(See [node1.yaml](./example/node1.yaml))

And the other two nodes with

```bash
./bin/tukkid -config ./example/node2.yaml
```
(See [node2.yaml](./example/node2.yaml))

and
```bash
./bin/tukkid -config ./example/node3.yaml
```
(See [node3.yaml](./example/node3.yaml))