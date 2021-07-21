package shard

import (
	"encoding/binary"
)

// TagPrototype represents a binary encoded tag that can be partially updated and cloned.
type TagPrototype [TagSize]byte

// SetPadding updates tag with a padding value.
func (t TagPrototype) SetPadding(n uint16) {
	binary.BigEndian.PutUint16(
		t[TagPaddingPosition:TagBlockSequencePosition], n)
}

// SetBlockSequence updates tag with a new block sequence.
func (t TagPrototype) SetBlockSequence(n uint64) {
	binary.BigEndian.PutUint64(
		t[TagBlockSequencePosition:TagShardSequencePosition], n)
}

// SetShardSequence updates tag with a new shard sequence.
func (t TagPrototype) SetShardSequence(n uint8) {
	t[TagShardSequencePosition] = byte(n)
}

// func (t TagPrototype) Write(b []byte) (n int, err error) {
// 	if len(b) < TagSize {
// 		return 0, io.ErrShortBuffer
// 	}
// 	return copy(b, t[:]), nil
// }
