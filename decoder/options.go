package decoder

import (
	"errors"
	"hash"
)

// Option configures the decoder.
type Option func(d *Decoder) error

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
