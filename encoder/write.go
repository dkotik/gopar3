package encoder

import (
	"bytes"
	"errors"
	"io"

	"github.com/dkotik/gopar3/shard"
)

// BatchWriter consumes batches and commits them to disk.
type BatchWriter interface {
	WriteAllBatches(<-chan (*batch)) error
}

type singleDestinationBatchWriter struct {
	w     io.Writer // replace with telomereWriter
	proto shard.TagPrototype
}

func (s *singleDestinationBatchWriter) WriteAllBatches(bb <-chan (*batch)) error {
	// TODO: write telomere

	var i uint32
	var j uint8
	for b := range bb {
		j = 0
		s.proto.SetPadding(b.padding)
		for _, shard := range b.shards {
			if shard == nil {
				break // channel was closed
			}
			n, err := io.Copy(s.w, bytes.NewReader(shard))
			if err != nil {
				return err
			}
			if n != int64(len(shard)) {
				return errors.New("failed to write the entire shard")
			}
			s.proto.SetBatchSequence(i)
			s.proto.SetShardSequence(j)
			j++

			n, err = io.Copy(s.w, bytes.NewReader(s.proto[:]))
			if err != nil {
				return err
			}
			if n != int64(len(shard)) {
				return errors.New("failed to write the entire shard tag")
			}
			// TODO: write checksum
			// TODO: write telomere
		}
		i++
	}

	// TODO: write telomere
	return nil
}

// type commitWriter func(io.Writer, *commit) error
//
// type commit struct {
// 	shard         []byte
// 	padding       uint32
// 	shardSequence uint8
// 	batchSequence uint32
// }
//
// func newWriter(w io.Writer, proto shard.TagPrototype) commitWriter {
// 	return func(w io.Writer, c *commit) (err error) {
// 		if _, err = io.Copy(w, bytes.NewReader(c.shard)); err != nil {
// 			return
// 		}
// 		proto.SetPadding(c.padding)
// 		proto.SetShardSequence(c.shardSequence)
// 		proto.SetBatchSequence(c.batchSequence)
//
// 		n, err := io.Copy(w, bytes.NewReader(proto[:]))
// 		if err != nil {
// 			return
// 		}
// 		if n != shard.TagSize {
// 			return errors.New("could not fit the tag in")
// 		}
// 		return nil
// 	}
// }
