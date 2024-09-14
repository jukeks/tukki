<img width="150" height="150" align="left" style="float: left; margin: 0 10px 0 0;" alt="tukki logo" src="./logo.png">   

# tukki

A toy key-value store built to explore log-structured merge-tree (LSM tree) 
concepts.

<br />

## Operations

tukki supports only three operations:

* Get key
* Set key
* Delete key

Keys and values are both UTF-8 strings.

## Database structures

### Memtable

Memtable is an in-memory representation of the current live segment. Writes are
buffered there until written to disk but also newest values are read from it.
Memtable needs to be a an ordered datastructure as the subsequent disk format is
sorted by keys.

tukki uses a red black tree for its memtable.

### Write-Ahead Log

Database writes are first written to a write-ahead log (WAL). [In the style of 
MongoDB](https://www.mongodb.com/docs/manual/reference/configuration-options/#mongodb-setting-storage.journal.commitIntervalMs), the WAL writes are batched and flushed to disk every 100 ms.

```
         +-----------+         +----------+
  Write  |           |  Write  |          |
-------->|    WAL    +-------->| Memtable |
         |           |         |          |
         +-----------+         +----------+
```

Memtable can be constructed from WAL on restart or crash recovery.

WAL messages are of format

```proto
message WalEntry {
    string key = 1;
    string value = 2;
    bool deleted = 3;
}
```

### SSTable

Sorted-String Table (SSTable) is the on-disk format of the database segments.
SSTables are immutable files but can be merged and compacted.

SSTable records look like

```proto
message SSTableRecord {
    string key = 1;
    string value = 2;
    bool deleted = 3;
}
```

### Bloom filters

tukki uses bloom filters for fast key memberships tests to find relevant segments.

### Indexes

Tukki has per segment primary key indexes. The indexes are stored as series of

```proto
message IndexEntry {
    string key = 1;
    uint64 offset = 2;
}
```

i.e. each IndexEntry tells at which offset a key's record is located in the
segment file.

### Segment journal

Segment journal records creation of new segments and merging existing segments. 
Segment journal is read at startup to find all segments and possible incomplete
operations.


### Disk format

tukki uses length prefixed protobuf messages as on disk format for all its types.

```
 0                   1                   2                   3   
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         Message Length                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Protobuf message                       |
+                                                               +
```

Length is a 32-bit unsigned integer which represents how many following bytes 
the message takes.

The messages are specificed at [`tukki.storage` protobuf package](proto/tukki/storage/).
