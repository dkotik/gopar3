package gopar3

import (
	"encoding/binary"
)

// TagPrototype represents a binary encoded tag that can be partially updated and cloned.
type TagPrototype [TagSize]byte

// SetPadding updates tag with a padding value.
func (t *TagPrototype) SetPadding(n uint32) {
	binary.BigEndian.PutUint32(
		t[TagPaddingPosition:TagBatchSequencePosition], n)
}

// SetBatchSequence updates tag with a new block sequence.
func (t *TagPrototype) SetBatchSequence(n uint32) {
	binary.BigEndian.PutUint32(
		t[TagBatchSequencePosition:TagShardSequencePosition], n)
}

// SetShardSequence updates tag with a new shard sequence.
func (t *TagPrototype) SetShardSequence(n uint8) {
	t[TagShardSequencePosition] = byte(n)
}

// func (t TagPrototype) Write(b []byte) (n int, err error) {
// 	if len(b) < TagSize {
// 		return 0, io.ErrShortBuffer
// 	}
// 	return copy(b, t[:]), nil
// }
