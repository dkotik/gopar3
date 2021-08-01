package encoder

import (
	"bytes"
	"io"
)

// Batch holds a group of required shards with added redundant shards.
// The redundant shards are set initially to <nil> to be filled with
// Reed-Solomon values. If there are not enough required shards,
// additional shards are repetitions of the last. Padding represents
// the number of bytes to discard when decoding the batch.
type Batch struct {
	shards  [][]byte
	padding uint32
}

func (e *Encoder) readBatchOfShards(r io.Reader) (*Batch, error) {
	stack := make([][]byte, e.RequiredShards+e.RedundantShards)
	var (
		i           uint8
		padding     int64
		morePadding = int64(e.shardSize)
		err         error
	)
	for ; i < e.RequiredShards; i++ {
		b := &bytes.Buffer{}
		if padding, err = io.CopyN(b, r, int64(e.shardSize)); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		stack[i] = b.Bytes()
	}

	padding = morePadding - padding // turn written n into padding number
	for j := i; j < e.RequiredShards; j++ {
		// fill in any missing shards by copies of the last
		stack[j] = stack[i]
		padding += morePadding
	}
	return &Batch{stack, uint32(padding)}, nil
}
