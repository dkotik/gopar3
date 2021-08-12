package scanner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/dkotik/gopar3/telomeres"
)

var (
	ErrShardTooSmall = errors.New("the shard is too small to contain a checksum")
	ErrShardBroken   = errors.New("the shard is broken, hashes do not match")
)

func NewScanner(r *telomeres.Decoder, withOptions ...Option) (*Scanner, error) {
	s := &Scanner{
		telomeresDecoder: r, // telomeres.NewDecoder(r, 4, 24, 2<<8),
	}
	withOptions = append(withOptions, WithDefaultOptions())
	if err := WithOptions(withOptions...)(s); err != nil {
		return nil, err
	}
	return s, nil
}

type Scanner struct {
	telomeresDecoder *telomeres.Decoder
	maxBytesPerShard int64
	checksumFactory  ChecksumFactory
	errorHandler     func(error) bool
	sequence         uint64
	// Validator              func([]byte) error // TODO: not really needed here
}

// Pipe transfers valid shards to the out channel. Checksum is truncated.
func (s *Scanner) Pipe(ctx context.Context, out chan<- ([]byte)) {
	go func() {
		for {
			shard, err := s.NextShard()
			// spew.Dump(err, string(shard))
			switch err {
			case ErrShardTooSmall:
				continue
			case telomeres.ErrEndReached:
				break
			case nil:
				// ignore
			default:
				if !s.errorHandler(
					fmt.Errorf("could not accept shard №%d: %w", s.sequence, err)) {
					break
				}
			}
			select {
			case <-ctx.Done():
				break
			case out <- shard:
				// continue
			}
		}
	}()
}

func (s *Scanner) NextShard() ([]byte, error) {
	buffer := &bytes.Buffer{}
	n, err := io.CopyN(buffer, s.telomeresDecoder, s.maxBytesPerShard)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	s.sequence++
	if n < 5 { // length of checksum, plus one byte
		return nil, ErrShardTooSmall
	}

	b := buffer.Bytes()
	checksumPosition := buffer.Len() - 4
	cs := s.checksumFactory()
	cs.Write(b[:checksumPosition])
	if bytes.Compare(cs.Sum(nil), b[checksumPosition:]) != 0 {
		return nil, ErrShardBroken
	}
	// if err = s.Validator(b); err != nil {
	// 	return nil, err
	// }
	return b[:checksumPosition], nil
}
