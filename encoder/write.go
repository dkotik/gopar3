package encoder

import (
	"bytes"
	"errors"
	"hash"
	"io"

	"github.com/dkotik/gopar3/shard"
	"github.com/dkotik/gopar3/telomeres"
)

// BatchWriter consumes batches and commits them to disk.
type BatchWriter func(<-chan (*Batch)) error

var (
	// ErrShardDidNotFit triggers when there are not enough bytes to write the entire shard.
	ErrShardDidNotFit = errors.New("failed to write the entire shard")
	// ErrShardTagDidNotFit triggers when there are not enough bytes to write the entire shard tag.
	ErrShardTagDidNotFit = errors.New("failed to write the entire shard tag")
)

type encodingSettings struct { // everything needed to complete encoding
	tag                shard.Tag
	telomereLength     int
	telomereBufferSize int
	checksum           hash.Hash32
}

// NewSingleDestinationWriter puts all the batches sequentially into the same Writer.
func NewSingleDestinationWriter(w io.Writer, s *encodingSettings) BatchWriter {
	return func(bb <-chan (*Batch)) (err error) {
		telw := telomeres.NewEncoder(w, s.telomereLength, s.telomereBufferSize)
		chkw := &checkSumWriter{telw, s.checksum}
		proto := s.tag.Prototype()

		if _, err = telw.Cut(); err != nil {
			return
		}

		ec := make(chan (error))
		go func() {
			var err error
			defer func() {
				ec <- err // capture error
			}()

			var i uint32
			var j uint8
			var n int64

			for b := range bb {
				j = 0
				proto.SetPadding(b.padding)
				for _, shard := range b.shards {
					if shard == nil {
						break // channel was closed
					}
					n, err = io.Copy(chkw, bytes.NewReader(shard))
					if err != nil {
						return
					}
					if n != int64(len(shard)) {
						err = ErrShardDidNotFit
						return
					}
					proto.SetBatchSequence(i)
					proto.SetShardSequence(j)
					j++

					n, err = io.Copy(chkw, bytes.NewReader(proto[:]))
					if err != nil {
						return
					}
					if n != int64(len(shard)) {
						err = ErrShardTagDidNotFit
						return
					}
					if err = chkw.Cut(); err != nil {
						return
					}
					if _, err = telw.Cut(); err != nil {
						return
					}
				}
				i++
			}

			_, err = telw.Cut()
			return
		}()
		return <-ec // wait for error
	}
}
