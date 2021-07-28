package encoder

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/dkotik/gopar3/shard"
	"github.com/klauspost/reedsolomon"
)

// AddPadding fill the last buffer to full shard size. If the number of given shards is less than required for Reed-Solomon computation, additional buffers filled with padding bytes are appended to the shard list.
func (e *Encoder) AddPadding(shards []*bytes.Buffer) (uint16, error) {
	req := int(e.RequiredShards) // TODO: move conversions to constructor?
	padding := e.shardSize - shards[req].Len()
	if padding > 0 {
		n, err := shards[req].Write(bytes.Repeat([]byte{PaddingByte}, padding))
		if err != nil {
			return 0, err
		}
		if n != padding {
			return 0, io.ErrShortBuffer
		}
	}

	if have := len(shards); have < req {
		empty := bytes.Repeat([]byte{PaddingByte}, e.shardSize)
		for ; have < req; have++ {
			b := &bytes.Buffer{}
			b.Grow(e.shardSize + shard.TagSize)
			n, err := b.Write(empty)
			if err != nil {
				return 0, err
			}
			if n != padding {
				return 0, io.ErrShortBuffer
			}
			shards = append(shards, b)
			padding += n
		}
	}

	if padding > int(^uint16(0)) {
		// TODO: cannot have more than 65535 padding,
		// which means 65535 / 256 = 255 maxShardsize
		// which is VERY limiting
		// this needs to be accounted for in the option
		return 0, errors.New("padding value is overflowing")
	}

	return uint16(padding), nil
}

func (e *Encoder) createShards(block uint64, base []*bytes.Buffer) error {
	req, red := int(e.RequiredShards), int(e.RedundantShards)
	total := req + red
	if l := len(base); l != req {
		return fmt.Errorf("need %d pieces, but only got %d", req, l)
	}
	enc, err := reedsolomon.New(req, red)
	if err != nil {
		return err
	}

	padding, err := e.AddPadding(base)
	if err != nil {
		return err
	}

	data := make([][]byte, total)
	for i, b := range base {
		data[i] = b.Bytes()
	}
	for i := req; i < total; i++ {
		data[i] = nil // those will be filled
	}
	if err = enc.Encode(data); err != nil { // fill
		return err
	}

	for i := req; i < total; i++ { // add new shards
		base = append(base, bytes.NewBuffer(data[i]))
	}

	tag := e.prototype // tag every shard
	tag.SetPadding(padding)
	tag.SetBlockSequence(block)
	for i := 0; i < total; i++ {
		tag.SetShardSequence(uint8(i))
		n, err := base[i].Write(tag[:])
		if err != nil {
			return err
		}
		if n < shard.TagSize {
			return io.ErrShortBuffer
		}
	}

	return nil
}
