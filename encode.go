package gopar3

import (
	"bytes"
	"context"
	"errors"
	"hash"
	"io"

	"github.com/dkotik/gopar3/telomeres"
	"github.com/klauspost/reedsolomon"
)

// Encoder adds data resiliency to its input.
type Encoder struct {
	RequiredShards  uint8
	RedundantShards uint8
	ChosenCheckSum  hash.Hash32

	blockSize int
	swap      *Swap
	out       []*telomeres.Encoder
	w         []io.Writer
}

func NewEncoder(blockSize int) (*Encoder, error) {
	e := &Encoder{
		blockSize: blockSize,
	}
	if uint16(e.RequiredShards+e.RedundantShards) > 256 { // TODO: test for this
		return nil, errors.New("sum of data and parity shards cannot exceed 256")
	}

	return nil, nil
}

// Use context.Context instead
// func (e *Encoder) IsDone() bool {
// 	return false
// }

// Encode uses a pipe pattern to split the contents of the reader into writers while encoding each block.
func (e *Encoder) Encode(ctx context.Context, r io.Reader) (err error) {
	// Stage 1: read chunks
	// read blocksize*RequiredShards of bytes
	// pad to the required length

	// Stage 2: generate Reed-Solomon shards
	enc, err := reedsolomon.New(int(e.RequiredShards), int(e.RedundantShards))
	if err != nil {
		return err
	}
	data := make([][]byte, e.RequiredShards+e.RedundantShards)
	// refs := make([]SwapReference, e.RequiredShards+e.RedundantShards)
	var i uint8
	for ; i < e.RequiredShards; i++ {
		// data[i] = <-stream
		// check length
		// what if stream is closed? fill with padded data, create a buffer for each // pad at the reader
	}
	// set padding for tag
	if err = enc.Encode(data); err != nil { // this fills up redundant shards
		return err
	}
	for i = 0; i < e.RequiredShards; i++ {
		// add tags to data
		// feed data refs to disk writers
	}
	for i = 0; i < e.RedundantShards; i++ {
		// add tags to data
		// feed data refs to disk writers
	}

	// Stage 3: commit the shards to output

	return nil
}

func (e *Encoder) chunk(
	ctx context.Context,
	r io.Reader,
	s *Swap,
	stream chan<- (SwapReference),
) (err error) {
	// read blocksize*RequiredShards of bytes
	// pad to the required length
	limit := int64(e.blockSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		b := &bytes.Buffer{}
		// b.Grow(blockSize * e.RequiredShards) ???
		ref, err := s.Reserve(b)
		if err != nil {
			return err
		}
		_, err = io.CopyN(b, r, limit)
		if err != nil {
			return err
		}
		stream <- ref
	}
	return
}

func commit(t *telomeres.Encoder, crossCheckFrequency uint, s *Swap, stream <-chan (SwapReference)) (err error) {
	wcross := &checkSumWriter{t, nil}
	wshard := &checkSumWriter{wcross, nil}

	var i uint
	for ref := range stream {
		swap, err := s.Retrieve(ref)
		if err != nil {
			return err
		}
		if _, err = io.Copy(wshard, swap); err != nil {
			return err
		}
		if err = wshard.Cut(); err != nil {
			return err
		}
		if _, err = t.Cut(); err != nil {
			return err
		}

		i++
		if i%crossCheckFrequency == 0 {
			if err = wcross.Cut(); err != nil {
				return err
			}
			if _, err = t.Cut(); err != nil {
				return err
			}
		}
	}
	// final cross-check
	if err = wcross.Cut(); err != nil {
		return err
	}
	if _, err = t.Cut(); err != nil {
		return err
	}
	return nil
}
