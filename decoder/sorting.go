package decoder

import (
	"bytes"
	"sort"
)

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

func (d *Decoder) SortAndBatch(in <-chan ([]byte)) chan<- ([][]byte) {
	out := make(chan ([][]byte), d.batchSize)
	go func() {
		shards := make([][]byte, 0, d.batchSize*3)
		// TODO: instead of next loop,
		// 			// write a loop to fill up the shards, sort it, try to batch, repeat
		for shard := range in {
			// TODO: if does not match what we are looking for, append to shards

			// if matched, feed to out
			out <- [][]byte{shard}
			// sort the shards?
			// feed them one by one if match
			for _, shard = range shards {
				out <- [][]byte{shard}
			}
		}
		close(out)
	}()
	return out
}

// SortAscending re-arranges BigEndian encoded slice in ascending order.
func SortAscending(q [][]byte) {
	sort.Slice(q, func(i, j int) bool {
		return bytes.Compare(q[i], q[j]) < 0
	})
}

// SortDescending re-arranges BigEndian encoded slice in descending order.
func SortDescending(q [][]byte) {
	sort.Slice(q, func(i, j int) bool {
		return bytes.Compare(q[i], q[j]) > 0
	})
}
