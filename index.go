package gopar3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"slices"
	"sync"

	"github.com/dkotik/gopar3/telomeres"
	"golang.org/x/sync/errgroup"
)

// Shard is a piece of a data. A certain number of shards are
// required to recover one file block.
type Shard struct {
	Source        string
	Size          int64
	FirstByte     int64
	LastByte      int64
	CastagnoliSum uint32
	Error         string
	Tag
}

// Differentiator returns hexadecimal representation of
// source Castagnoli sum, source size, and quorum
// with appended decimal shard size. The combination
// should be unique for different sources.
func (s *Shard) Differentiator() string {
	return fmt.Sprintf(
		"%x_%db",
		s.Tag.Bytes()[:DifferentiatorSize],
		s.Size,
	)
}

// Load reads associated data from disk into bytes
func (s *Shard) Load(ctx context.Context) (_ []byte, err error) {
	// TODO: this function should be run in a map[source]*Reader
	f, err := os.Open(s.Source)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	if _, err = f.Seek(s.FirstByte, io.SeekStart); err != nil {
		return nil, err
	}
	r := telomeres.NewDecoder(f)
	// if err = r.SeekChunk(ctx); err != nil {
	// 	return nil, err
	// }
	b := &bytes.Buffer{}
	if _, err = r.StreamChunk(ctx, b); err != nil {
		return nil, err
	}
	return b.Bytes()[TagBytesForCRC+TagSize:], nil
}

type File struct {
	Shards        []*Shard
	Quorum        uint8
	Size          uint64
	Padding       uint64
	Batches       uint16
	CastagnoliSum uint32
	Error         string
}

// Index is a map of known shards arranged by [Tag.BlockDifferentiator]
// gathered from a list of files that could contain recovery data
// for any number of files. Index can be saved to complete
// recovery operations in more than one execution.
type Index map[string]*File

func (i Index) Normalize() (err error) {
	if len(i) == 0 {
		return errors.New("no data shards were detected in input files")
	}
	var shardSize int64
	for _, f := range i {
		for _, shard := range f.Shards {
			if shard.Error != "" {
				continue // do not consider data from corrupt shards
			}
			// there is no need for statisticalMeanOfSortedSlice
			// because shards are already grouped by differentiator
			// as the Index key
			f.CastagnoliSum = shard.Tag.SourceCRC
			f.Size = shard.Tag.SourceSize
			f.Quorum = shard.Tag.ShardQuorum
			shardSize = shard.Size
			break // found one recoverable
		}
		if f.CastagnoliSum == 0 {
			f.Error = "there are no recoverable shards"
			continue
		}

		slices.SortFunc(f.Shards, func(a, b *Shard) int {
			// return a negative number when a < b,
			// a positive number when a > b,
			// zero when a == b
			if a.Tag.ShardBatch < b.Tag.ShardBatch {
				return -1
			} else if a.Tag.ShardBatch > b.Tag.ShardBatch {
				return 1
			}
			if a.Tag.ShardOrder < b.Tag.ShardOrder {
				return -1
			} else if a.Tag.ShardOrder > b.Tag.ShardOrder {
				return 1
			}
			return 0
		})

		f.Batches = uint16(math.Ceil(
			float64(f.Size)/float64(shardSize*int64(f.Quorum)),
		)) + 1
		f.Padding = uint64(f.Batches)*uint64(f.Quorum)*uint64(shardSize) - f.Size

		// validate file
		batch := make(map[uint8]uint32)
		currentBatch := uint16(0)
		quorum := int(f.Quorum)
		knownSum := uint32(0)
		ok := false
		for _, shard := range f.Shards {
			if shard.Error != "" {
				continue // do not consider data from corrupt shards
			}
			if shard.Tag.ShardBatch != currentBatch {
				if len(batch) < quorum {
					f.Error = fmt.Sprintf("batch %d has %d recoverable shards instead of %d required", currentBatch, len(batch), quorum)
					break
				}
				currentBatch++
				if shard.Tag.ShardBatch != currentBatch {
					f.Error = fmt.Sprintf("there are no recoverable shards for batch %d", currentBatch)
					break
				}
				batch = make(map[uint8]uint32) // reset
			}
			if knownSum, ok = batch[shard.Tag.ShardOrder]; ok {
				if shard.CastagnoliSum == knownSum {
					shard.Error = "duplicate shard"
				} else {
					shard.Error = "duplicate shard with corrupt CRC"
				}
			} else {
				batch[shard.Tag.ShardOrder] = shard.CastagnoliSum
			}
		}
		if currentBatch+1 < f.Batches {
			f.Error = fmt.Sprintf("there are only %d recoverable batches out of %d required for restoration", currentBatch+1, f.Batches)
		}
		// panic(f.Batches)
	}
	return nil
}

// NewIndex scans files for shards and recovers as much information
// about them as possible to assess the presence and possibility
// of data recovery in those shards.
func NewIndex(ctx context.Context, files ...string) (index Index, err error) {
	for _, source := range files {
		info, err := os.Stat(source)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			// TODO: queue files within the folders
			// dir, err := os.ReadDir(source)
			// if err != nil {
			//   return nil, err
			// }
			// for _, file := range dir {
			//   if file.IsDir() {
			//     continue
			//   }
			//   files = append(files, file)
			// }
			return nil, fmt.Errorf("cannot read a directory: %s", source)
		}
	}

	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(runtime.NumCPU())
	index = make(Index)
	mu := &sync.Mutex{}

	for _, file := range files {
		wg.Go(func() (err error) {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer func() {
				if err == io.EOF {
					err = nil
				}
				err = errors.Join(err, f.Close())
			}()

			r := NewReader(file, f)
			for {
				// b := &bytes.Buffer{}
				// shard, err := r.NextShard(ctx, b)
				// log.Fatalf("%s", b.String())
				shard, err := r.NextShard(ctx, io.Discard)
				if err != nil {
					break
				}
				differentiator := shard.Differentiator()
				mu.Lock()
				file, ok := index[differentiator]
				if !ok {
					file = &File{}
					index[differentiator] = file
				}
				file.Shards = append(file.Shards, shard)
				mu.Unlock()
			}
			return err
		})
	}

	return index, errors.Join(wg.Wait(), index.Normalize())
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

func statisticalMeanOfSortedSlice[T any](s []T) T {
	switch count := len(s); count {
	case 0:
		var zero T
		return zero
	default:
		return s[count/2]
	}
}
