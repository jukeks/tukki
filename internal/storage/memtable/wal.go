package memtable

import (
	"io"
	"log"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/journal"
	walv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/wal/v1"
)

type Wal struct {
	journal *journal.Journal
}

func OpenWal(dbDir string, journalName files.Filename, mt Memtable) (*Wal, error) {
	handle := func(r *journal.JournalReader) error {
		return readJournal(r, mt)
	}
	if journalName == "" {
		log.Fatalf("journal name is empty")
	}
	j, err := journal.OpenJournal(dbDir, journalName, journal.WriteModeAsync, handle)
	if err != nil {
		return nil, err
	}

	return &Wal{j}, nil
}

func (mtj *Wal) Set(key, value string) error {
	return mtj.journal.Writer.Write(&walv1.WalEntry{
		Key:     key,
		Value:   value,
		Deleted: false,
	})
}

func (mtj *Wal) Delete(key string) error {
	return mtj.journal.Writer.Write(&walv1.WalEntry{
		Key:     key,
		Deleted: true,
	})
}

func (mtj *Wal) Close() error {
	return mtj.journal.Close()
}

func (mtj *Wal) Snapshot() []byte {
	return mtj.journal.Writer.Snapshot()
}

func readJournal(journalReader *journal.JournalReader, mt Memtable) error {
	for {
		journalEntry := &walv1.WalEntry{}
		err := journalReader.Read(journalEntry)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if journalEntry.Deleted {
			mt.Delete(journalEntry.Key)
		} else {
			mt.Insert(journalEntry.Key, journalEntry.Value)
		}
	}
}

func (mtj *Wal) Size() uint64 {
	stat, err := mtj.journal.File.Stat()
	if err != nil {
		log.Fatalf("failed to get journal file size: %v", err)
	}

	return uint64(stat.Size())
}
