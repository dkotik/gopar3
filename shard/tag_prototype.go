package shard

import "encoding/binary"

// TagPrototype represents a binary encoded tag that can be partially updated and cloned.
type TagPrototype [shardTagSize]byte

// SetBlockSequence updates tag with a new block sequence.
func (t TagPrototype) SetBlockSequence(n uint64) {
	binary.BigEndian.PutUint64(
		t[shardTagBlockSequencePosition:shardTagShardSequencePosition], n)
}

// SetShardSequence updates tag with a new shard sequence.
func (t TagPrototype) SetShardSequence(n uint8) {
	t[shardTagShardSequencePosition] = byte(n)
}
