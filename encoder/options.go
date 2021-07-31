package encoder

import (
	"errors"
	"math"
)

// Option modifies the encoder.
type Option func(e *Encoder) error

// WithOptions aggregates multiple options into one.
func WithOptions(options ...Option) Option {
	return func(e *Encoder) error {
		for _, option := range options {
			if err := option(e); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithDefaultOptions are used when no options are provided.
func WithDefaultOptions() Option {
	return WithOptions(
		WithShardSize(512),
		WithCrossCheckFrequency(4),
	)
}

// WithGrowthFactor adjusts required and redundant shards as a percentage of total. Growth factor represents the approximate increase in size of the total output compared to the input file.
func WithGrowthFactor(totalFragments uint, g float64) Option {
	return func(e *Encoder) error {
		if totalFragments < 3 {
			return errors.New("cannot use less than three fragments")
		}
		if g <= 1 {
			return errors.New("growth factor must be great than 1")
		}
		e.RedundantShards = uint8(math.Ceil(float64(totalFragments) / g))
		e.RequiredShards = uint8(totalFragments) - e.RedundantShards
		return nil
	}
}

// WithShardSize sets the size of created chunks. Smaller chunks make the output more resilient at the cost of disk space and recovery speed.
func WithShardSize(inbytes int) Option {
	// if padding > int(^uint16(0)) { // see const shard.MaxPadding
	// 	// TODO: cannot have more than 65535 padding,
	// 	// which means 65535 / 256 = 255 maxShardsize
	// 	// which is VERY limiting
	// 	// this needs to be accounted for in the option
	// 	return 0, errors.New("padding value is overflowing")
	// }
	return func(e *Encoder) error {
		if inbytes < 1 {
			return errors.New("cannot encode using 0 or negative shard size")
		}
		e.shardSize = inbytes
		return nil
	}
}

// WithCrossCheckFrequency sets the number of blocks, after which a cross-check hash is written to the output writer.
func WithCrossCheckFrequency(f uint8) Option {
	return func(e *Encoder) error {
		if f == 0 {
			return errors.New("cross check frequency must be greater than zero")
		}
		e.crossCheckFrequency = uint(f)
		return nil
	}
}
