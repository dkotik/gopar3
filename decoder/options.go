package decoder

import (
	"errors"
	"hash"

	"github.com/dkotik/gopar3"
	"github.com/dkotik/gopar3/scanner"
)

// Option configures the decoder.
type Option func(d *Decoder) error

// WithOptions aggregates multiple options into one.
func WithOptions(options ...Option) Option {
	return func(e *Decoder) error {
		for _, option := range options {
			if err := option(e); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithDefaultOptions makes sure that all the missing values are set.
func WithDefaultOptions() Option {
	return func(d *Decoder) error {
		defaults := make([]Option, 0)
		if d.maxShardSize == 0 {
			defaults = append(defaults, WithMaxShardSize(2<<20*16))
		}
		if d.sniffDepth == 0 {
			defaults = append(defaults, WithSniffDepth(36))
		}
		if d.checksumFactory == nil {
			defaults = append(defaults, WithChecksumFactory(scanner.NewChecksum))
		}
		return WithOptions(defaults...)(d)
	}
}

// WithSniffDepth determines how many pieces are worked on at a time.
func WithSniffDepth(limit uint16) Option {
	return func(d *Decoder) error {
		if limit < 9 {
			return errors.New("cannot work with less than 9 shards at a time")
		}
		d.sniffDepth = limit
		return nil
	}
}

// WithSniffDepth determines how many pieces are worked on at a time.
func WithMaxShardSize(limit int64) Option {
	return func(d *Decoder) error {
		if limit < gopar3.TagSize+1 {
			return errors.New("cannot use less shard bytes than required for a tag")
		}
		d.maxShardSize = limit
		return nil
	}
}

// TODO: the batch is read from the tag
// // WithBatchSize determines how many pieces are worked on at a time.
// func WithBatchSize(size int) Option {
// 	return func(d *Decoder) error {
// 		if size < 9 {
// 			return errors.New("cannot work with less than 9 shards at a time")
// 		}
// 		d.batchSize = size
// 		return nil
// 	}
// }

// WithChecksumFactory provides checksums for validating shards.
func WithChecksumFactory(factory func() hash.Hash32) Option {
	return func(d *Decoder) error {
		if factory == nil {
			return errors.New("cannot use an empty factory")
		}
		d.checksumFactory = factory
		return nil
	}
}

// WithValidator?
