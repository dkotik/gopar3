package decoder

import (
	"bytes"
	"context"
	"errors"
	"hash"
	"io"

	"github.com/dkotik/gopar3/swap"
)

// Scanner locates valid shards.
type Scanner struct {
	r         io.Reader
	readLimit int64
	checksum  hash.Hash32
	swap      *swap.Swap
}

func (s *Scanner) Scan(ctx context.Context, valid chan<- (swap.SwapReference)) (err error) {
	var (
		n   int64
		ref swap.SwapReference
	)

	isValid := NewValidator(s.checksum)
	for {
		b := &bytes.Buffer{}
		n, err = io.CopyN(b, s.r, s.readLimit)
		if err != nil {
			return
		}
		if n >= s.readLimit {
			return errors.New("reached the read limit")
		}
		if !isValid(b.Bytes()) {
			continue // corrupted shard
		}

		ref, err = s.swap.Reserve(b)
		if err != nil {
			return
		}
		valid <- ref
	}

	return
}

// NewValidator sets up a function that validates bytes that end with a 32bit checksum.
func NewValidator(c hash.Hash32) func(b []byte) bool {
	return func(b []byte) bool {
		// TODO: make sure this function does not cause race conditions due to using the table?
		length := len(b) - 4 // roll back
		if length <= 0 {
			return false // not enough bytes
		}
		c.Reset()
		n, err := c.Write(b[:length])
		if err != nil {
			return false
		}
		if n != length {
			return false
		}
		return 0 == bytes.Compare(
			b[length:], c.Sum(b[:length]))
	}
}
