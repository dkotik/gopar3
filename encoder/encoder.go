package encoder

import (
	"errors"

	"github.com/dkotik/gopar3"
)

// Encoder adds data resiliency to its input.
type Encoder struct {
	requiredShards      uint8
	redundantShards     uint8
	prototype           gopar3.TagPrototype
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
	e = &Encoder{
		errc: make(chan (error)),
	}

	withOptions = append(withOptions, WithDefaultOptions())
	if err = WithOptions(withOptions...)(e); err != nil {
		return nil, err
	}

	if uint16(e.requiredShards+e.redundantShards) > 256 { // TODO: test for this
		return nil, errors.New("sum of data and parity shards cannot exceed 256")
	}

	return nil, nil
}
