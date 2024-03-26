package gopar3

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
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
		if batch > lastBatch {
			// TODO: fix
			return fmt.Errorf("batch %d out of maximum range of %d", batch, lastBatch)
			// break
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
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() (err error) {
		defer close(forReconstruction)
		for _, batch := range batches {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			shards := make([][]byte, len(batch))
			for i, shard := range batch {
				shards[i], err = shard.Load(ctx)
				if err != nil {
					return err
				}

				// fmt.Printf("%d -------------------------------\n", len(shards[i]))
				// fmt.Println(string(shards[i]))
				// if len(shards[i]) == 0 {
				// 	// od -j 1024 -N 16
				// 	fmt.Printf("check range: od -j %d -N %d", shard.CursorStart, shard.CursorEnd-shard.CursorStart)
				// 	panic("loaded an empty shard without an error")
				// }
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case forReconstruction <- shards:
			}
		}
		return nil
	})

	forWriting := make(chan [][]byte, 4)
	wg.Go(func() (err error) {
		defer close(forWriting)
		rs, err := reedsolomon.New(quorum, mostShards-quorum)
		if err != nil {
			return err
		}
		for shards := range forReconstruction {
			if err = rs.ReconstructData(shards); err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case forWriting <- shards[:quorum]:
			}
		}
		return nil
	})

	wg.Go(func() (err error) {
		var (
			written    int64
			writeLimit = int64(f.Size)
			padding    int
			n          int
			crc        = crc32.New(castagnoliTable)
		)
		for shards := range forWriting {
			// padding calculations assume that all shards are the same size
			n = len(shards[0]) // shard size here for determining padding
			if padding = int(written) + (len(shards) * n) - int(writeLimit); padding > 0 {
				shards = shards[:len(shards)-padding/n]
				if cutLast := padding % n; cutLast > 0 {
					shards[len(shards)-1] = shards[len(shards)-1][:n-cutLast]
				}
			}

			for _, shard := range shards {
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

		if written != writeLimit {
			return fmt.Errorf("the number of written bytes %d does not match expected file size %d", written, f.Size)
		}
		if crc.Sum32() != f.CastagnoliSum {
			log.Print(crc.Sum32(), f.CastagnoliSum)
			return errors.New("circular redundancy check does not match the expected value; the file is corrupt and cannot be recovered")
		}
		return nil
	})

	return wg.Wait()
}
