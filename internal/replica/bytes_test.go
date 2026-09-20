package replica

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	"github.com/jukeks/tukki/internal/db"
	"google.golang.org/protobuf/proto"
)

func TestByteCommandEncoding(t *testing.T) {
	original := command{Version: 1, Op: "set", Key: []byte{255, 0}, Value: []byte{128, 0}, Min: []byte{0}, Max: []byte{255}}
	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded command
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Op != original.Op || !bytes.Equal(decoded.Key, original.Key) || !bytes.Equal(decoded.Value, original.Value) || !bytes.Equal(decoded.Min, original.Min) || !bytes.Equal(decoded.Max, original.Max) {
		t.Fatalf("command changed: %+v", decoded)
	}
	// These old text values are also valid base64, but must remain text.
	if err := json.Unmarshal([]byte(`{"op":"set","key":"YWJj","value":"ZGVm"}`), &decoded); err != nil {
		t.Fatal(err)
	}
	if string(decoded.Key) != "YWJj" || string(decoded.Value) != "ZGVm" {
		t.Fatalf("legacy text changed: %+v", decoded)
	}
	if err := json.Unmarshal([]byte(`{"version":2}`), &decoded); err == nil {
		t.Fatal("accepted an unknown command version")
	}
}

func TestBinaryRaftStorage(t *testing.T) {
	dir := t.TempDir()
	d, err := db.OpenDatabase(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if d != nil {
			d.Close()
		}
	})
	store := NewRaftukki(d)
	stableKey, stableValue := []byte{255, 0}, []byte{128, 0, 255}
	if err := store.Set(stableKey, stableValue); err != nil {
		t.Fatal(err)
	}
	logs := []*raft.Log{
		{Index: 1, Term: 2, Type: raft.LogCommand, Data: []byte{255, 0, 128}, Extensions: []byte{254, 0}, AppendedAt: time.Unix(123, 0)},
		{Index: 2, Term: 2, Type: raft.LogCommand, Data: []byte{0, 255}},
	}
	if err := store.StoreLogs(logs); err != nil {
		t.Fatal(err)
	}
	raw, err := d.Get(logKey(1))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := proto.Marshal(logFromRaft(logs[0]))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw, expected) {
		t.Fatal("Raft log is not stored as raw protobuf")
	}
	// Retain access to logs written by earlier versions.
	if err := d.Set(logKey(3), []byte(base64.StdEncoding.EncodeToString(expected))); err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	d = nil
	d, err = db.OpenDatabase(dir)
	if err != nil {
		t.Fatal(err)
	}
	store = NewRaftukki(d)
	got, err := store.Get(stableKey)
	if err != nil || !bytes.Equal(got, stableValue) {
		t.Fatalf("stable value = %x, %v", got, err)
	}
	for _, index := range []uint64{1, 3} {
		var got raft.Log
		if err := store.GetLog(index, &got); err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(logFromRaft(&got), logFromRaft(logs[0])) {
			t.Fatalf("log %d changed: %+v", index, got)
		}
	}
}
