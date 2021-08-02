package encoder

import (
	"errors"

	"github.com/dkotik/gopar3/shard"
)

const (
	telomereEncoderBufferSize = 4096

	// PaddingByte is added to shards to make them all the same size.
	PaddingByte = byte('#')
)

// Encoder adds data resiliency to its input.
type Encoder struct {
	RequiredShards      uint8 // // TODO: private
	RedundantShards     uint8 // // TODO: private
	prototype           shard.TagPrototype
	shardSize           int // TODO: replace with shard size // int64?
	telomeresLength     int
	telomeresBufferSize int
	crossCheckFrequency uint
	errc                chan (error)
	// ChosenCheckSum  hash.Hash32
	// w                   []io.Writer
}

// NewEncoder initializes the encoder with options. Default options are used, if no options were specified.
func NewEncoder(withOptions ...Option) (e *Encoder, err error) {
	if len(withOptions) == 0 {
		withOptions = []Option{WithDefaultOptions()}
	}

	e = &Encoder{
		errc: make(chan (error)),
	}
	if err = WithOptions(withOptions...)(e); err != nil {
		return nil, err
	}

	if uint16(e.RequiredShards+e.RedundantShards) > 256 { // TODO: test for this
		return nil, errors.New("sum of data and parity shards cannot exceed 256")
	}

	return nil, nil
}
