package decoder

import (
	"context"
	"hash"
	"io"
)

// Decoder restores the original data from a set of streams.
type Decoder struct {
	checksumFactory func() hash.Hash32
}

// Decode reads the streams, orders shards, recovers data, and writes it out to destination writer.
func (d *Decoder) Decode(ctx context.Context, w io.Writer, streams []io.Reader) (err error) {
	// pipe
	// read streams concurrently while rejecting bad chunks

	// order chunks and try to restore

	// collect restored data and write it out

	return
}
