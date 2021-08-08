package gopar3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash"
	"io"

	"github.com/dkotik/gopar3/shard"
)

var (
	ErrShardTooSmall = errors.New("the shard is too small to contain a tag")
	ErrShardBroken   = errors.New("the shard is broken, hashes do not match")
)

// func NewReader(r io.Reader, withOptions ...ReaderOption) (*Reader, error) {
// 	reader := &Reader{
// 		source: telomeres.NewDecoder(r, 4, 24, 2<<10),
// 	}
// 	withOptions = append(withOptions, WithReaderDefaultOptions())
// 	if err := WithReaderOptions(withOptions...)(reader); err != nil {
// 		return nil, err
// 	}
// 	return reader, nil
// }

type Scanner struct {
	MaxBytesBeforeGivingUp int64
	ChecksumFactory        func() hash.Hash32
	ErrorHandler           func(error) bool
	// Validator              func([]byte) error // TODO: not really needed here
}

// Pipe transfers valid shards to the out channel. Checksum is truncated.
func (s *Scanner) Pipe(ctx context.Context, r io.Reader, out chan<- ([]byte)) {
	go func() {
		var i uint64 = 1
		for {
			shard, err := s.NextShard(r)
			if err != nil {
				if err != io.EOF {
					if !s.ErrorHandler(
						fmt.Errorf("could not accept shard №%d: %w", i, err)) {
						break
					}
				}
				break // stop on end of file
			}
			select {
			case <-ctx.Done():
				break
			case out <- shard:
				i++
			}
		}
	}()
}

func (s *Scanner) NextShard(r io.Reader) ([]byte, error) {
	buffer := &bytes.Buffer{}
	n, err := io.CopyN(buffer, r, s.MaxBytesBeforeGivingUp)
	if err != nil {
		return nil, err
	}
	if n <= shard.TagSize {
		return nil, ErrShardTooSmall
	}

	b := buffer.Bytes()
	checksumPosition := buffer.Len() - 4
	cs := s.ChecksumFactory()
	if bytes.Compare(cs.Sum(b[:checksumPosition]), b[checksumPosition:]) != 0 {
		return nil, ErrShardBroken
	}
	// if err = s.Validator(b); err != nil {
	// 	return nil, err
	// }
	return b[:checksumPosition], nil
}
