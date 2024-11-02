package keyvalue

import (
	"fmt"
	"io"
	"log"
)

type SubIterator interface {
	Get() (IteratorEntry, error)
	Seek(key string) error
	Progress()
	Close() error
}

type Iterator struct {
	min              string
	max              string
	returnTombstones bool
	iterators        []SubIterator
	usableIterators  []SubIterator
}

func NewIterator(min, max string, returnTombstones bool, iterators ...SubIterator) (*Iterator, error) {
	usableIterators := make([]SubIterator, 0, len(iterators))
	usableIterators = append(usableIterators, iterators...)
	i := &Iterator{
		min:              min,
		max:              max,
		returnTombstones: returnTombstones,
		iterators:        iterators,
		usableIterators:  usableIterators,
	}

	if min != "" {
		i.usableIterators = make([]SubIterator, 0)
		for _, iter := range iterators {
			if err := iter.Seek(min); err != nil {
				if err == io.EOF {
					continue
				}
				return nil, fmt.Errorf("failed to seek to min key: %w", err)
			}

			i.usableIterators = append(i.usableIterators, iter)
		}
	}

	return i, nil
}

func (i *Iterator) Close() {
	for _, iter := range i.iterators {
		if err := iter.Close(); err != nil {
			log.Printf("failed to close iterator: %v", err)
		}
	}
}

func (i *Iterator) Next() (IteratorEntry, error) {
	for {
		// find first error free segment that is not past the end
		result := i.usableIterators[0]
		canProceed := false
		for _, iter := range i.usableIterators {
			current, err := iter.Get()
			if err == nil && (i.max == "" || current.Key <= i.max) {
				result = iter
				canProceed = true
				break
			}
		}
		if !canProceed {
			return IteratorEntry{}, io.EOF
		}

		for _, iter := range i.usableIterators {
			current, err := iter.Get()
			if err != nil {
				if err != io.EOF {
					return IteratorEntry{}, fmt.Errorf("failed to read next entry from %+v: %w", iter, err)
				}
				continue
			}

			if i.max != "" && current.Key > i.max {
				continue
			}

			currentResult, _ := result.Get()
			if current.Key < currentResult.Key {
				result = iter
			}
		}

		ret, _ := result.Get()

		// there can be earlier entries in other segments
		// we need to advance them to the next entry
		// to avoid exposing old value
		for _, iter := range i.usableIterators {
			current, _ := iter.Get()
			if current.Key == ret.Key {
				iter.Progress()
			}
		}

		if !i.returnTombstones && ret.Deleted {
			continue
		}

		return ret, nil
	}
}
