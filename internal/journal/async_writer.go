package journal

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/jukeks/tukki/internal/storage"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type AsynchronousJournalWriter struct {
	w WriteSyncer
	b *bufio.Writer

	writeBuff chan protoreflect.ProtoMessage
	errors    chan error
	closed    chan bool
	close     chan bool
	err       error
}

func NewAsynchronousJournalWriter(w WriteSyncer) *AsynchronousJournalWriter {
	writer := &AsynchronousJournalWriter{
		w:         w,
		b:         bufio.NewWriterSize(w, 128*1024),
		writeBuff: make(chan protoreflect.ProtoMessage, 1000),
		close:     make(chan bool, 1),
		errors:    make(chan error, 1),
		closed:    make(chan bool, 1),
	}
	go writer.writer()

	return writer
}

func (j *AsynchronousJournalWriter) Close() error {
	j.close <- true
	<-j.closed
	return nil
}

func (j *AsynchronousJournalWriter) Write(journalEntry protoreflect.ProtoMessage) error {
	err := j.checkError()
	if err != nil {
		return err
	}

	j.writeBuff <- journalEntry
	return nil
}

func (j *AsynchronousJournalWriter) checkError() error {
	if j.err != nil {
		return j.err
	}

	select {
	case j.err = <-j.errors:
		return j.err
	default:
		return nil
	}
}

func (j *AsynchronousJournalWriter) writer() {
	defer func() {
		j.closed <- true
	}()

	for {
		err := j.processBatch()
		if err != nil {
			log.Printf("failed to process batch: %v", err)
			j.errors <- fmt.Errorf("failed to process batch: %w", err)
			return
		}

		select {
		case <-j.close:
			err := j.processBatch()
			if err != nil {
				log.Printf("failed to process batch: %v", err)
				j.errors <- fmt.Errorf("failed to process batch: %w", err)
				return
			}
			return
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}

func (j *AsynchronousJournalWriter) processBatch() error {
	written := false
messagesAvailable:
	for {
		select {
		case msg := <-j.writeBuff:
			_, err := storage.WriteLengthPrefixedProtobufMessage(j.b, msg)
			if err != nil {
				return fmt.Errorf("failed to write journal entry: %w", err)
			}

			written = true
		default:
			break messagesAvailable
		}
	}

	if written {
		err := j.b.Flush()
		if err != nil {
			return fmt.Errorf("failed to flush: %w", err)
		}

		err = j.w.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync: %w", err)
		}
	}

	return nil
}
