package encoder

import (
	"errors"
	"math"
)

const (
	defaultBufferSize = 4096
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
	return func(e *Encoder) error {
		defaults := make([]Option, 0)
		if e.requiredShards == 0 {
			defaults = append(defaults, WithRequiredShards(6))
		}
		if e.redundantShards == 0 {
			defaults = append(defaults, WithRedundantShards(3))
		}
		if e.shardSize == 0 {
			defaults = append(defaults, WithShardSize(512))
		}
		if e.telomeresBufferSize == 0 { // telomeres were not set at all
			defaults = append(defaults, WithTelomeres(9))
		}
		if e.crossCheckFrequency == 0 {
			defaults = append(defaults, WithCrossCheckFrequency(4))
		}
		return WithOptions(defaults...)(e)
	}
}

// WithRequiredShards sets the miminal number of shards neccessary for reconstruction.
func WithRequiredShards(n uint8) Option {
	return func(e *Encoder) error {
		if n == 0 {
			return errors.New("cannot use 0 required shards")
		}
		e.requiredShards = n
		return nil
	}
}

// WithRedundantShards sets the number of shards that can be lost without preventing data reconstruction.
func WithRedundantShards(n uint8) Option {
	return func(e *Encoder) error {
		if n == 0 {
			return errors.New("cannot use 0 redundant shards")
		}
		e.redundantShards = n
		return nil
	}
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
		redundant := uint8(math.Ceil(float64(totalFragments) / g))
		return WithOptions(
			WithRequiredShards(uint8(totalFragments)-redundant),
			WithRedundantShards(redundant),
		)(e)
	}
}

// WithShardSize sets the size of created chunks. Smaller chunks make the output more resilient at the cost of disk space and recovery speed.
func WithShardSize(inbytes int) Option {
	return func(e *Encoder) error {
		if inbytes < 1 {
			return errors.New("cannot encode using 0 or negative shard size")
		}
		e.shardSize = inbytes
		return nil
	}
}

// WithTelomeres sets the number of telomere characters inserted between shards and cross-checks.
func WithTelomeres(n uint8) Option {
	return func(e *Encoder) error {
		e.telomeresLength = int(n)
		if e.telomeresBufferSize == 0 {
			return WithTelomeresBufferSize(defaultBufferSize)(e)
		}
		return nil
	}
}

// WithTelomeresBufferSize sets the buffer size.
func WithTelomeresBufferSize(n int) Option {
	return func(e *Encoder) error {
		e.telomeresBufferSize = n
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
