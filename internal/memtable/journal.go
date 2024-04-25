package memtable

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	journalv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/journal/v1"
)

type MembtableJournal struct {
	journal *journal.Journal
}

func NewJournal(dbDir string, mt Memtable) (*MembtableJournal, error) {
	handle := func(r *journal.JournalReader) error {
		return readJournal(r, mt)
	}
	j, err := journal.OpenJournal(dbDir, "journal", handle)
	if err != nil {
		return nil, err
	}

	return &MembtableJournal{j}, nil
}

func (mtj *MembtableJournal) Set(key, value string) error {
	return mtj.journal.Writer.Write(&journalv1.JournalEntry{
		Key:     key,
		Value:   value,
		Deleted: false,
	})
}

func (mtj *MembtableJournal) Delete(key string) error {
	return mtj.journal.Writer.Write(&journalv1.JournalEntry{
		Key:     key,
		Deleted: true,
	})
}

func (mtj *MembtableJournal) Close() error {
	return mtj.journal.File.Close()
}

func readJournal(journalReader *journal.JournalReader, mt Memtable) error {
	for {
		journalEntry := &journalv1.JournalEntry{}
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
