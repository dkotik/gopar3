package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/dkotik/gopar3"
)

func createFilterByBlock(block uint32, expectedLength int) func([]byte) bool {
	matcher := make([]byte, 4)
	binary.BigEndian.PutUint32(matcher, block)
	return func(b []byte) bool {
		return bytes.Compare(b[gopar3.TagBatchSequencePosition:gopar3.TagShardSequencePosition], matcher) == 0
	}
}

func (d *Decoder) selectRelatedShards(block uint32, q [][]byte) ([][]byte, error) {
	shards := make([][]byte, d.batchSize)
	var shardSequence int
	nfiltered := 0
	found := 0
	var duplicationTest, tag []byte
	filter := createFilterByBlock(block, d.shardSize)

	// TODO: make sure q is not sorted backwards, so getting the right shards first
	for _, s := range q {
		tag = s[len(s)-gopar3.TagSize-4:] // gopar3.TagSize-4
		if filter(tag) {
			shardSequence = int(tag[gopar3.TagShardSequencePosition])
			if duplicationTest = shards[shardSequence]; duplicationTest != nil {
				return nil, fmt.Errorf("shard #%d-%d is duplicated", block, shardSequence)
			}
			shards[shardSequence] = s
			found++
			if found == d.batchSize {
				break
			}
			continue
		}
		q[nfiltered] = s
		nfiltered++
	}

	// if found < d.requiredShards {
	// 	return nil, fmt.Errorf("block #%d contains only %d/%d required shards", block, found, d.requiredShards)
	// }

	q = q[:nfiltered]
	return shards, nil
}

func (d *Decoder) orderAndGroup(in <-chan ([]byte)) <-chan ([][]byte) {
	out := make(chan ([][]byte))
	go func() {
		defer close(out)
		var (
			depth   = int(d.sniffDepth)
			queue   = make([][]byte, depth)
			current []byte
			batch   [][]byte
			i       int
			block   uint32
			err     error
		)

		for {
			for i = len(queue); i < depth; i++ {
				current = <-in
				if current == nil {
					break
				}
				queue[i] = current
			}

			SortDescending(queue)
			batch, err = d.selectRelatedShards(block, queue)
			if err != nil {
				panic(err) // TODO: pass gracefully, close out channel?
			}
			out <- batch
			block++
		}
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
