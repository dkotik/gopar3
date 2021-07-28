package encoder

import (
	"bytes"
	"context"
	"errors"
	"hash"
	"io"

	"github.com/dkotik/gopar3/shard"
	"github.com/dkotik/gopar3/telomeres"
)

const (
	telomereEncoderBufferSize = 4096

	// PaddingByte is added to shards to make them all the same size.
	PaddingByte = byte('#')
)

// Encoder adds data resiliency to its input.
type Encoder struct {
	RequiredShards  uint8
	RedundantShards uint8
	ChosenCheckSum  hash.Hash32

	prototype           shard.TagPrototype
	shardSize           int // TODO: replace with shard size
	crossCheckFrequency uint
	telomeresLength     int
	w                   []io.Writer
	// out                 []*telomeres.Encoder
	// swap                *swap.Swap
}

// NewEncoder initializes the encoder with options. Default options are used, if no options were specified.
func NewEncoder(withOptions ...Option) (e *Encoder, err error) {
	if len(withOptions) == 0 {
		withOptions = []Option{WithDefaultOptions()}
	}

	e = &Encoder{}
	if err = WithOptions(withOptions...)(e); err != nil {
		return nil, err
	}

	if uint16(e.RequiredShards+e.RedundantShards) > 256 { // TODO: test for this
		return nil, errors.New("sum of data and parity shards cannot exceed 256")
	}

	return nil, nil
}

// Encode uses a pipe pattern to split the contents of the reader into writers while encoding each block.
func (e *Encoder) Encode(ctx context.Context, r io.Reader) (err error) {
	// prepare writers

	// Stage 1: read chunks
	// read blocksize*RequiredShards of bytes
	// pad to the required length

	// Stage 2: create Reed-Solomon shards

	// Stage 3: commit the shards to output

	return nil
}

func (e *Encoder) commit(ctx context.Context, w io.Writer, stream <-chan (*bytes.Buffer)) (err error) {
	t := telomeres.NewEncoder(w, e.telomeresLength, telomereEncoderBufferSize)
	wcross := &checkSumWriter{t, nil}
	wshard := &checkSumWriter{wcross, nil}

	var i uint
	for b := range stream {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if b == nil {
			return nil // finished
		}
		if _, err = io.Copy(wshard, b); err != nil {
			return err
		}
		if err = wshard.Cut(); err != nil {
			return err
		}
		if _, err = t.Cut(); err != nil {
			return err
		}

		i++
		if i%e.crossCheckFrequency == 0 {
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
