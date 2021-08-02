package encoder

import (
	"context"
	"io"
)

// Encode uses a pipe pattern to split the contents of the reader into writers while encoding each block.
func (e *Encoder) Encode(ctx context.Context, r io.Reader) (err error) {
	// prepare writers

	// Stage 1: read chunks
	inChannel := e.batchStream(r)

	// Stage 2: complete Reed-Solomon shards
	rsChannel := e.CompleteWithReedSolomon(inChannel)

	// Stage 3: commit the shards to output
	writer := e.NewSingleDestinationWriter("/tmp/test.gopar3")
	writer(rsChannel, e.errc)

	return <-e.errc
}
