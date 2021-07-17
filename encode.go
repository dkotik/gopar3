package gopar3

import (
	"hash"
	"io"
)

// Encoder adds data resiliency to its input.
type Encoder struct {
	RequiredShards  uint8
	RedundantShards uint8
	ChosenCheckSum  hash.Hash32

	w []io.Writer
}

// Encode uses a pipe pattern to split the contents of the reader into writers while encoding each block.
func (*Encoder) Encode(r io.Reader) error {
	// Stage 1: read chunks

	// Stage 2: generate Reed-Solomon shards

	// Stage 3: commit the shards to output

	return nil
}

func commit(w io.Writer, stream <-chan (SwapReference)) error {
	// wcross := &checkSumWriter{w, nil}
	// wshard := &checkSumWriter{wcross, nil}
	//
	// for ref := range stream {
	//
	// }
	return nil
}
