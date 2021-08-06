package decoder

import (
	"context"
	"fmt"
	"hash"
	"io"
)

type shardFilter func([]byte) bool

// NewDecoder constructs a decoder.
func NewDecoder(withOptions ...Option) (d *Decoder, err error) {
	d = &Decoder{}
	withOptions = append(withOptions, WithDefaultOptions())
	if err = WithOptions(withOptions...)(d); err != nil {
		return nil, err
	}
	return d, nil
}

// Decoder restores the original data from a set of streams.
type Decoder struct {
	batchSize       int
	sniffDepth      uint16
	maxShardSize    int64
	checksumFactory func() hash.Hash32
	shardFilter     shardFilter
	errc            chan (error)
}

// Decode reads the streams, orders shards, recovers data, and writes it out to destination writer.
func (d *Decoder) Decode(ctx context.Context, w io.Writer, streams []io.Reader) (err error) {
	in := make(chan ([]byte), d.sniffDepth)
	err = func() error {
		// move all of this inside sniff $$$$$$$$$$$$$$$$$$$$$$$$$$$$$
		stack, common, err := d.Sniff(in)
		if err != nil {
			return err
		}
		d.shardFilter = func() shardFilter {
			func(a []byte) bool {
			if len(a) != len(common) {
				return false
			}
			// TODO: compare to common
			return true
		}
		for _, shard := range stack {
			if d.shardFilter(shard) {
				in <- shard
			}
		}
		return nil
	}()
	if err != nil {
		return
	}
	for _, r := range streams {
		d.StartReading(r, in)
	}
	// pipe
	// read streams concurrently while rejecting bad chunks
	// discard chunks with less than tag's length +1

	// order chunks and try to restore

	// collect restored data and write it out

	return
}

func (d *Decoder) String() string {
	return fmt.Sprintf("[Decoder readLimit=%.2fMb sniffDepth=%d]",
		float64(d.maxShardSize)/float64(2<<20),
		d.sniffDepth)
}
