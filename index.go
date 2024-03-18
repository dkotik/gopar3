package gopar3

import (
	"context"
	"errors"
	"os"
	"sync"
)

// Shard is a piece of a data. A certain number of shards are
// required to recover one file block.
type Shard struct {
	Source        string
	CursorStart   int64
	CursorEnd     int64
	CursorGap     int64
	IsRecoverable bool
	WasRecovered  bool
	Tag           Tag
}

// Index is a map of known shards arranged by [Tag.BlockDifferentiator]
// gathered from a list of files that could contain recovery data
// for any number of files. Index can be saved to complete
// recovery operations in more than one execution.
//
// All methods are safe for concurrent use.
type Index struct {
	shards map[[DifferentiatorSize]byte][]Shard
	mu     *sync.Mutex
}

func (i *Index) AddShard(s Shard) error {
	block := s.Tag.BlockDifferentiator
	i.mu.Lock()
	shards, _ := i.shards[block]
	i.shards[block] = append(shards, s)
	i.mu.Unlock()
	return nil
}

func (i *Index) AddFile(
	ctx context.Context,
	source string,
	progress func(context.Context, *Index, *Shard) error,
) (err error) {
	if progress == nil {
		return errors.New("cannot use a <nil> progress function")
	}

	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	// r := telomeres.NewDecoder(f, 4, 100, 1024) // TODO: tweak
	// b := &bytes.Buffer{}
	// for {
	//
	// }
	return err
}
