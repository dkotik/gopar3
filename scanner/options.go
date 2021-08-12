package scanner

import (
	"errors"
	"hash"
)

type Option func(s *Scanner) error
type ChecksumFactory func() hash.Hash32

func WithOptions(options ...Option) Option {
	return func(s *Scanner) (err error) {
		for _, option := range options {
			if err = option(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithDefaultOptions() Option {
	return func(s *Scanner) error {
		defaults := make([]Option, 0)
		// TODO: fill out
		return WithOptions(defaults...)(s)
	}
}

func WithErrorHandler(f func(error) bool) Option {
	return func(s *Scanner) error {
		if f == nil {
			return errors.New("cannot use an empty error handler")
		}
		s.errorHandler = f
		return nil
	}
}

func WithMaxBytesPerShard(max int64) Option {
	return func(s *Scanner) error {
		if max < 5 {
			return errors.New("shards cannot contain less than 5 bytes")
		}
		s.maxBytesPerShard = max
		return nil
	}
}

func WithChecksumFactory(f ChecksumFactory) Option {
	return func(s *Scanner) error {
		if f == nil {
			return errors.New("cannot use an empty checksum factory")
		}
		s.checksumFactory = f
		return nil
	}
}
