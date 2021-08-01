package decoder

import (
	"fmt"
)

// SniffDemocratically determines predominant shard tag qualities by taking the most popular tag values from a given set of slices. Start and end mark the common part of the slice.
func SniffDemocratically(q [][]byte, start, end int) []byte {
	type rec struct {
		Index int
		Count uint16
	}

	// count the values
	cc := make(map[string]*rec)
	for i := 0; i < len(q); i++ {
		mark := fmt.Sprintf("%x", q[i][start:end]) // group by hex value
		if saved, ok := cc[mark]; ok {
			saved.Count++
			continue
		}
		cc[mark] = &rec{i, 1}
	}

	// spew.Dump(cc)

	// select most common signature
	var top *rec
	for _, v := range cc {
		if top == nil || top.Count < v.Count {
			top = v
		}
	}

	if top == nil {
		return nil
	}
	return q[top.Index]
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
