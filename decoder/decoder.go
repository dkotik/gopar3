package decoder

import (
	"context"
	"errors"
	"hash"
	"io"
)

// NewDecoder constructs decoder.
func NewDecoder(withOptions ...Option) (d *Decoder, err error) {
	d = &Decoder{}

	if err = WithOptions(withOptions...)(d); err != nil {
		return nil, err
	}

	if d.checksumFactory == nil {
		return nil, errors.New("decode must have a checksum provider")
	}

	return d, nil
}

// Decoder restores the original data from a set of streams.
type Decoder struct {
	batchSize       int
	checksumFactory func() hash.Hash32
}

// Decode reads the streams, orders shards, recovers data, and writes it out to destination writer.
func (d *Decoder) Decode(ctx context.Context, w io.Writer, streams []io.Reader) (err error) {
	// pipe
	// read streams concurrently while rejecting bad chunks
	// discard chunks with less than tag's length +1

	// order chunks and try to restore

	// collect restored data and write it out

	return
}

func (d *Decoder) orderAndGroup(
	in []<-chan ([]byte), out chan<- ([][]byte),
) (err error) {
	queue := make([][]byte, d.batchSize*2)

	var i, c, count int
	for {
		for i = 0; i < d.batchSize; i++ {
			for c = 0; c < len(in); c++ {
				select { // non-blocking take
				case queue[count] = <-in[c]:
					if queue[count] != nil {
						count++
					}
				default:
				}
			}
		}

		if count == 0 {
			break // no more data is coming
		}
		SortDescending(queue)

		// determine shard count once // TODO: auto
		shards := 9 // TODO: auto
		limit := count - shards
		if limit < 0 {
			limit = 0
		}
		out <- queue[limit:count]                       // next stage
		queue = append(queue[:limit], queue[count:]...) // drop the used part
	}

	return nil
}
