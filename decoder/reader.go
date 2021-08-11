package decoder

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/dkotik/gopar3"
)

// TODO: moving all of this to gopar3 Reader

var (
	ErrShardTooSmall = errors.New("the shard is too small to contain a tag")
	ErrShardBroken   = errors.New("the shard is broken, hashes do not match")
)

func (d *Decoder) StartReading(r io.Reader, out chan<- ([]byte)) {
	go func() {
		var i uint64
		for {
			shard, err := d.ReadShard(r)
			if err != nil {
				if err != io.EOF {
					d.errc <- fmt.Errorf("could not accept shard №%d: %w", i, err)
				}
				break
			}
			if !d.shardFilter(shard) {
				continue
			}
			out <- shard
		}
	}()
}

func (d *Decoder) ReadShard(r io.Reader) ([]byte, error) {
	buffer := &bytes.Buffer{}
	n, err := io.CopyN(buffer, r, d.maxShardSize)
	if err != nil {
		return nil, err
	}
	if n < gopar3.TagSize {
		return nil, ErrShardTooSmall
	}

	b := buffer.Bytes()
	checksumPosition := buffer.Len() - 4
	cs := d.checksumFactory()
	if bytes.Compare(cs.Sum(b[:checksumPosition]), b[checksumPosition:]) != 0 {
		return nil, ErrShardBroken
	}
	return b[:checksumPosition], nil
}
