package gopar3

import (
	"context"
	"io"
)

// PaddingByte is a shard filler for making them equal size.
const PaddingByte = '?'

// BatchLoader loads a batch of padded shards from an [io.Reader]
// for [reedsolomon.Encoder] reconstruction. The parity
// shards are `nil`.
type BatchLoader struct {
	Quorum    int
	Shards    int
	ShardSize int
}

// Stream sends batches into a channel until context cancellation
// or reader error. Closes the channel when the reader stops.
func (g *BatchLoader) Stream(
	ctx context.Context,
	c chan<- [][]byte,
	r io.Reader,
) (err error) {
	var batch [][]byte
	for {
		batch, _, err = g.Load(r)
		switch err {
		case nil:
			c <- batch
		case io.EOF:
			c <- batch
			close(c)
			return nil
		default:
			close(c)
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

// Seek repositions the cursor at the expected boundary of the
// next shard batch. Useful for running multiple batch loaders
// in parallel.
func (g *BatchLoader) Seek(r io.Seeker, batch int) (n int64, err error) {
	return r.Seek(int64(batch*g.Quorum*g.ShardSize), io.SeekStart)
}

// Load returns a batch of padded shards and the the number
// of data bytes that were read. The parity
// shards are `nil`. The final batch returns io.EOF as error.
func (g *BatchLoader) Load(r io.Reader) (batch [][]byte, loaded int, err error) {
	batch = make([][]byte, g.Shards)
	var (
		n int
		i int
		j int
	)

	for i = range batch[:g.Quorum] {
		shard := make([]byte, g.ShardSize)
		n, err = io.ReadFull(r, shard)
		loaded += n
		switch err {
		case nil:
			batch[i] = shard
		case io.ErrUnexpectedEOF: // short write
			for j = range shard[n:] {
				// fill the rest of the shard with padding chars
				shard[n+j] = PaddingByte
			}
			batch[i] = shard
			i++
			fallthrough
		case io.EOF:
			if i < g.Quorum {
				goto padRemaining
			}
			return batch, loaded, io.EOF
		}
		batch[i] = shard
	}
	return batch, loaded, nil

padRemaining:
	shard := make([]byte, g.ShardSize)
	for j = range g.ShardSize {
		shard[j] = PaddingByte
	}
	for j = range batch[i:g.Quorum] {
		batch[i+j] = shard
	}
	return batch, loaded, io.EOF
}
