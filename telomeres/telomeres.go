/*
Package telomeres provides an encoder and a decoder that write
chunks of data surrounded by repeated uniform byte sequences
that designate chunk boundaries.

Telomeres amplify data resiliency by preventing two data chunks
from being corrupted by having an error on their boundary.

In addition, telomeres provide a more reliable mechanism
of chunk detection that does not depend on counting bytes
associated with chunk length. The file can suffer more
damage while its chunks remain recognizable.
*/
package telomeres

import (
	"errors"
	"fmt"
)

const (
	telomereMarkByte   = ':'
	telomereEscapeByte = '\\'
)

// Telomeres is a factory that provides [Encoder]s and [Decoder]s
// with the same parameters.
type Telomeres struct {
	mark       byte
	escape     byte
	minimum    int
	blockLimit int64
	bufferSize int
}

// New creates a valid [Telomeres] factory.
func New(withOptions ...Option) (_ *Telomeres, err error) {
	o := &telomeresOptions{}
	for _, option := range append(
		withOptions,
		withDefaultMarkByte(),
		withDefaultEscapeByte(),
		withDefaultMinimumCount(),
		withDefaultBufferSize(),
		withDefaultBlockReadLimit(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create a telomeres codec factory: %w", err)
		}
	}
	return &Telomeres{
		mark:       *o.mark,
		escape:     *o.escape,
		minimum:    o.minimum,
		blockLimit: o.blockLimit,
		bufferSize: o.bufferSize,
	}, nil
}

type telomeresOptions struct {
	mark       *byte
	escape     *byte
	minimum    int
	blockLimit int64
	bufferSize int
}

// Option configures [Telomeres] factory.
type Option func(*telomeresOptions) error

// WithMarkByte specifies the repeated telomere byte that separates
// data chunks. Colon `:` character is the default.
func WithMarkByte(c byte) Option {
	return func(o *telomeresOptions) error {
		if o.mark != nil {
			return errors.New("mark byte is already set")
		}
		o.mark = &c
		return nil
	}
}

func withDefaultMarkByte() Option {
	return func(o *telomeresOptions) error {
		if o.mark != nil {
			return nil
		}
		return WithMarkByte(':')(o)
	}
}

// WithEscapeByte specifies the repeated telomere byte that separates
// data chunks. Colon `\` character is the default.
func WithEscapeByte(c byte) Option {
	return func(o *telomeresOptions) error {
		if o.escape != nil {
			return errors.New("escape byte is already set")
		}
		o.escape = &c
		return nil
	}
}

func withDefaultEscapeByte() Option {
	return func(o *telomeresOptions) error {
		if o.escape != nil {
			return nil
		}
		return WithEscapeByte('\\')(o)
	}
}

// WithMinimumCount specifies the minimum number of repetions
// of a telomere byte to recognize data chunk boundary. Eight
// is the default.
func WithMinimumCount(n int) Option {
	return func(o *telomeresOptions) error {
		if n < 4 {
			return errors.New("telomere count of less than four is not reliable")
		}
		if o.minimum != 0 {
			return errors.New("telomere minimum count is already set")
		}
		o.minimum = n
		return nil
	}
}

func withDefaultMinimumCount() Option {
	return func(o *telomeresOptions) error {
		if o.minimum != 0 {
			return nil
		}
		return WithMinimumCount(8)(o)
	}
}

// WithBufferSize specifies the buffer size for either encoding or decoding.
func WithBufferSize(n int) Option {
	return func(o *telomeresOptions) error {
		if n < 64 {
			return errors.New("buffer size cannot be less than 64")
		}
		if o.bufferSize != 0 {
			return errors.New("buffer size is already set")
		}
		o.bufferSize = n
		return nil
	}
}

func withDefaultBufferSize() Option {
	return func(o *telomeresOptions) error {
		if o.bufferSize != 0 {
			if o.bufferSize < o.minimum*2 {
				return errors.New("buffer size should be at least double the minimum telomere count")
			}
			return nil
		}
		return WithBufferSize(o.minimum * 8 * 1024)(o)
	}
}

// WithBlockReadLimit specifies the maximum number of bytes
// that can be expected within a data chunk. Once exceeded,
// the remainig bytes are discarded until the end of the
// next telomere boundary as defined by [WithMarkByte]
// and [WithMinimumCount].
func WithBlockReadLimit(n int64) Option {
	return func(o *telomeresOptions) error {
		if n < 64 {
			return errors.New("block read limit cannot be less than 64")
		}
		if o.blockLimit != 0 {
			return errors.New("block read limit is already set")
		}
		o.blockLimit = n
		return nil
	}
}

func withDefaultBlockReadLimit() Option {
	return func(o *telomeresOptions) error {
		if o.blockLimit != 0 {
			return nil
		}
		return WithBlockReadLimit(int64(o.bufferSize) * 512)(o)
	}
}
