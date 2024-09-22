# Multi node tukki


```bash
./bin/tukkid -port 50010 -node-id 1 -raft-port 50011 -db-dir ./node-1 -raft-peers 2@localhost:50021,3@localhost:50031 -sstable-peers 2@localhost:50020,3@localhost:50030 -init

./bin/tukkid -port 50020 -node-id 2 -raft-port 50021 -db-dir ./node-2 -sstable-peers 1@localhost:50010,3@localhost:50030

./bin/tukkid -port 50030 -node-id 3 -raft-port 50031 -db-dir ./node-3 -sstable-peers 2@localhost:50020,1@localhost:50010
```