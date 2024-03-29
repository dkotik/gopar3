package gopar3

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkotik/gopar3/telomeres"
	"github.com/klauspost/reedsolomon"
	"golang.org/x/sync/errgroup"
)

const (
	ShardLimit      = 1<<(TagBytesForShardOrder*8) - 1
	ShardBatchLimit = 1<<(TagBytesForShardBatch*8) - 1
	SourceSizeLimit = 1<<(TagBytesForSourceSize*8) - 1
)

// castagnoliTable sources [crc.New] with 0x82f63b78
// polynomial. It is known for superior error detection
// and use in BitTorrent and iSCSI protocols.
//
// Example of BitTorrent use:
// https://github.com/anacrolix/torrent/blob/master/bep40.go
var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

func Inflate(
	ctx context.Context,
	destination string,
	source string,
	shardQuorum uint8,
	shardParity uint8,
	shardSize int,
) (err error) {
	f, err := os.Stat(source)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return errors.New("cannot inflate a directory")
	}
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, r.Close())
	}()

	var (
		l = &BatchLoader{
			Quorum:    int(shardQuorum),
			Shards:    int(shardQuorum + shardParity),
			ShardSize: shardSize,
		}
	)
	rs, err := reedsolomon.New(
		l.Quorum, l.Shards-l.Quorum,
		reedsolomon.WithAutoGoroutines(shardSize),
	)
	if err != nil {
		return err
	}

	tag, err := NewTag(ctx, r, shardQuorum)
	if err != nil {
		return err
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	var w io.Writer
	{ // create io.Writer
		f, err = os.Stat(destination)
		if err == os.ErrNotExist {
			w, err = os.Create(destination)
			if err != nil {
				return err
			}
		} else if err == nil && f.IsDir() {
			ext := filepath.Ext(source)
			base := strings.TrimSuffix(filepath.Base(source), ext)
			w, err = os.Create(filepath.Join(
				destination,
				fmt.Sprintf(`%s%x%s.gopar3`, base, tag.SourceCRC, ext),
			))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	wg, ctx := errgroup.WithContext(ctx)
	batches := make(chan [][]byte, 4)
	wg.Go(func() (err error) {
		defer close(batches)
		var batch [][]byte
		for {
			batch, _, err = l.Load(r)
			if err != nil {
				// panic(fmt.Sprintf("%+v", batch))
				if err == io.EOF {
					return nil
				}
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case batches <- batch:
			}
		}
	})

	batchesWithParity := make(chan [][]byte, 4)
	wg.Go(func() (err error) {
		defer close(batchesWithParity)
		for batch := range batches {
			if err = rs.Reconstruct(batch); err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case batchesWithParity <- batch:
			}
		}
		return nil
	})

	wtlm, err := telomeres.NewEncoder(w, 5)
	if err != nil {
		return err
	}
	shardWriter, err := NewWriter(wtlm, NewSequentialTagger(tag, shardQuorum+shardParity))
	if err != nil {
		return err
	}
	wg.Go(func() (err error) {
		for batch := range batchesWithParity {
			// log.Printf("batch size is %d", len(batch))
			for _, shard := range batch {
				if _, err = shardWriter.Write(shard); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return wg.Wait()
}

func CastagnoliSum(ctx context.Context, r io.Reader) (uint32, error) {
	var (
		crc  = crc32.New(castagnoliTable)
		b    = make([]byte, 2*32*1024)
		n    int
		rerr error
		werr error
	)
	for rerr == nil {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		n, rerr = io.ReadFull(r, b)
		if _, werr = crc.Write(b[:n]); werr != nil {
			return 0, werr
		}
		switch rerr {
		case io.EOF, io.ErrUnexpectedEOF:
		default:
			return 0, rerr
		}
	}
	return crc.Sum32(), nil
}
