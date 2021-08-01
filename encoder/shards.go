package encoder

import (
	"bytes"
	"fmt"

	"github.com/klauspost/reedsolomon"
)

func (e *Encoder) createShards(block uint32, base []*bytes.Buffer) error {
	req, red := int(e.RequiredShards), int(e.RedundantShards)
	total := req + red
	if l := len(base); l != req {
		return fmt.Errorf("need %d pieces, but only got %d", req, l)
	}
	enc, err := reedsolomon.New(req, red)
	if err != nil {
		return err
	}

	// padding, err := e.AddPadding(base)
	// if err != nil {
	// 	return err
	// }

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

	// tag := e.prototype // tag every shard
	// tag.SetPadding(padding)
	// tag.SetBatchSequence(block)
	// for i := 0; i < total; i++ {
	// 	tag.SetShardSequence(uint8(i))
	// 	n, err := base[i].Write(tag[:])
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if n < shard.TagSize {
	// 		return io.ErrShortBuffer
	// 	}
	// }

	return nil
}
