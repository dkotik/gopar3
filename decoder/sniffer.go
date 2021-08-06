package decoder

import (
	"errors"
	"fmt"

	"github.com/dkotik/gopar3/shard"
)

func (d *Decoder) Sniff(in chan ([]byte)) (stack [][]byte, popular []byte, err error) {
	shards := make([][]byte, 0, d.sniffDepth)
	// var i int
	for i := uint16(0); i < d.sniffDepth; i++ {
		shard := <-in
		if shard == nil {
			break
		}
		shards = append(shards, shard)
	}
	popular, err = SniffDemocratically(shards)
	if err != nil {
		return nil, nil, err
	}
	return stack, popular, nil
}

// SniffDemocratically determines predominant shard tag qualities by taking the most popular tag values from a given set of slices.
func SniffDemocratically(q [][]byte) ([]byte, error) {
	length := len(q)
	if length <= 5 {
		return nil, errors.New("cannot use less than 5 shards to discover common values")
	}
	type rec struct {
		Index int
		Count int
		// Length int
	}

	// count similar values grouped by a slice and total length
	cc := make(map[string]*rec)
	for i := 0; i < len(q); i++ {
		length := len(q[i])
		if length < shard.TagSize+1 {
			continue // skip over short byte slices
		}
		mark := fmt.Sprintf("%x::%d", q[i][length-shard.TagSize:length-4], len(q[i]))
		if saved, ok := cc[mark]; ok {
			saved.Count++
			continue
		}
		cc[mark] = &rec{
			Index: i,
			// Length: length,
			Count: 1,
		}
	}

	// spew.Dump(cc)

	// select most common signature
	var top *rec
	for _, v := range cc {
		if top == nil || top.Count < v.Count {
			top = v
		}
	}

	if top == nil || top.Count <= length/3 {
		return nil, errors.New("there was not even a third of common values")
	}
	return q[top.Index], nil
}

// func medianUint8(bunch []uint8) uint8 {
// 	sort.Slice(bunch, func(i int, j int) bool {
// 		return bunch[i] > bunch[j]
// 	})
// 	return bunch[len(bunch)/2]
// }
//
// func medianUint16(bunch []uint16) uint16 {
// 	sort.Slice(bunch, func(i int, j int) bool {
// 		return bunch[i] > bunch[j]
// 	})
// 	return bunch[len(bunch)/2]
// }
