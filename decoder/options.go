package decoder

import (
	"errors"
	"hash"
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

// WithBatchSize determines how many pieces are worked on at a time.
func WithBatchSize(size int) Option {
	return func(d *Decoder) error {
		if size < 9 {
			return errors.New("cannot work with less than 9 shards at a time")
		}
		d.batchSize = size
		return nil
	}
}

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
