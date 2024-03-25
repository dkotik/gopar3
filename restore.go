package gopar3

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"slices"

	"github.com/klauspost/reedsolomon"
	"golang.org/x/sync/errgroup"
)

// Restore writes recovered contents of a file using shards
// of a normalized [Index].
func Restore(ctx context.Context, w io.Writer, f *File) (err error) {
	if f.Error != "" {
		return errors.New(f.Error)
	}

	batches := make([][]*Shard, f.Batches)
	lastBatch := int(f.Batches - 1)
	batch := 0
	for _, shard := range f.Shards {
		if shard.Error != "" {
			continue
		}
		batch = int(shard.Tag.ShardBatch)
		fmt.Println("batch:", batch)
		if batch > lastBatch {
			// TODO: fix
			// return fmt.Errorf("batch %d out of maximum range of %d", batch, lastBatch)
			break
		}
		batches[batch] = append(batches[batch], shard)
	}
	// panic(fmt.Sprintf("%+v", batches))

	quorum := int(f.Quorum)
	available := 0
	mostShards := 0
	for i, batch := range batches {
		available = len(batch)
		if available < quorum {
			return fmt.Errorf("cannot recover batch #%d, because there are only %d shards available out of %d required", i, available, quorum)
		}
		slices.SortFunc(batch, func(a, b *Shard) int {
			// return a negative number when a < b,
			// a positive number when a > b,
			// zero when a == b
			if a.Tag.ShardOrder < b.Tag.ShardOrder {
				return -1
			} else if a.Tag.ShardOrder > b.Tag.ShardOrder {
				return 1
			}
			return 0
		})
		if available > mostShards {
			mostShards = available
		}
	}

	forReconstruction := make(chan [][]byte, 4)
	forWriting := make(chan [][]byte, 4)
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() (err error) {
		for _, batch := range batches {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			shards := make([][]byte, len(batch))
			for i, shard := range batch {
				shards[i], err = shard.Load()
				if err != nil {
					return err
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case forReconstruction <- shards:
			}
		}
		close(forReconstruction)
		return nil
	})

	wg.Go(func() (err error) {
		rs, err := reedsolomon.New(quorum, mostShards-quorum)
		if err != nil {
			return err
		}
		for shards := range forReconstruction {
			rs.Reconstruct(shards)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case forWriting <- shards[:quorum]:
			}
		}
		close(forWriting)
		return nil
	})

	wg.Go(func() (err error) {
		var (
			written int64
			n       int
			crc     = crc32.New(castagnoliTable)
		)
		// panic("w")
		for shards := range forWriting {
			for _, shard := range shards {
				// io.Copy(w, bytes.NewReader(shard)) ?
				n, err = w.Write(shard)
				if err != nil {
					return err
				}
				if _, err = crc.Write(shard); err != nil {
					return err
				}
				written += int64(n)
			}
		}

		if written != int64(f.Size) {
			return fmt.Errorf("the number of written bytes %d does not match expected file size %d", written, f.Size)
		}
		if crc.Sum32() != f.CastagnoliSum {
			return errors.New("circular redundancy check does not match the expected value; the file is corrupt and cannot be recovered")
		}
		return nil
	})

	return wg.Wait()
}
