# Multi node tukki

```bash
./bin/node -port 50010 -node-id 1 -raft-port 50011 -db-dir ./node-1 -peers 2@localhost:50021,3@localhost:50031

./bin/node -port 50020 -node-id 2 -raft-port 50021 -db-dir ./node-2

./bin/node -port 50030 -node-id 3 -raft-port 50031 -db-dir ./node-3
```