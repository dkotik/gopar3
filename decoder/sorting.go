package decoder

import (
	"bytes"
	"sort"
)

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
