package decoder

import (
	"context"
	"fmt"
	"hash"
	"io"
)

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
	requiredShards  uint8
	redundantShards uint8
	shardSize       int
	maxShardSize    int64
	checksumFactory func() hash.Hash32
	shardFilter     func([]byte) bool // shard filter rejects shards that do not match the most popular sniff set tag signature
	errc            chan (error)
}

// Decode reads the streams, orders shards, recovers data, and writes it out to destination writer.
func (d *Decoder) Decode(ctx context.Context, w io.Writer, streams []io.Reader) (err error) {
	in := make(chan ([]byte), d.sniffDepth)
	defer close(in)

	// TODO: launch consumers of in channel here, so they are ready for shards of the sniffer
	batches := d.orderAndGroup(in)
	complete := d.CompleteWithReedSolomon(batches)

	if err = d.SniffAndSetupFilter(in, streams); err != nil {
		return err
	}
	for _, r := range streams {
		d.StartReading(r, in)
	}

	return d.WriteAll(w, complete) // TODO: add context?
}

func (d *Decoder) String() string {
	return fmt.Sprintf("[Decoder readLimit=%.2fMb sniffDepth=%d]",
		float64(d.maxShardSize)/float64(2<<20),
		d.sniffDepth)
}
