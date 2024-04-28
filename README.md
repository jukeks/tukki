<img width="150" height="150" align="left" style="float: left; margin: 0 10px 0 0;" alt="MS-DOS logo" src="./logo.png">   

# tukki

A toy key-value store built to explore log-structured merge-tree (LSM tree) 
concepts.

## Operations

tukki supports only three operations:

* Get key
* Set key
* Delete key

## Disk format

tukki uses length prefixed protobuf messages as on disk format.

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

## Memtable

Memtable is an in-memory representation of the current live segment. Writes are
buffered there until written to disk but also newest values are read from it.
Memtable needs to be a an ordered datastructure as the subsequent disk format is
sorted by keys.

tukki uses a red black tree for its memtable.

### Write-Ahead Log

Database writes are first written to a write-ahead log (WAL).

```
        +-----------+        +----------+
 Write  |           | Write  |          |
------->|    WAL    +------->| Memtable |
        |           |        |          |
        +-----------+        +----------+
```

Memtable can be constructed from WAL on restart or crash recovery.

WAL messages are of format

```
message WalEntry {
    string key = 1;
    string value = 2;
    bool deleted = 3;
}
```

## SSTable

Sorted-String Table (SSTable) is the on-disk format of the segments. SSTables 
are immutable files but can be merged and compacted.

SSTable records look like

```
message SSTableRecord {
    string key = 1;
    string value = 2;
    bool deleted = 3;
}
```

### Segment journal

Segment journal records creation of new segments and merging existing segments. 
Journal is read at startup to find all segments and possible incomplete
operations.


## Indexes

Not implemented yet.