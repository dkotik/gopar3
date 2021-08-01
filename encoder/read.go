package encoder

import (
	"bytes"
	"io"
)

type batch struct {
	shards  [][]byte
	padding uint32
}

func (e *Encoder) readBatchOfShards(r io.Reader) (*batch, error) {
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
	return &batch{stack, uint32(padding)}, nil
}
