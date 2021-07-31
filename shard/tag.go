package shard

import (
	"encoding/binary"
	"io"
	"math/rand"
	"time"
)

// Tag pieces mark byte boundaries of each element:
const (
	// VersionBytePosition is first at 0.
	TagBlockDifferentiatorPosition = 0 + 1
	TagRequiredShardsPosition      = TagBlockDifferentiatorPosition + 4
	TagRedundantShardsPosition     = TagRequiredShardsPosition + 1
	TagPaddingPosition             = TagRedundantShardsPosition + 2
	TagBlockSequencePosition       = TagPaddingPosition + 2
	TagShardSequencePosition       = TagBlockSequencePosition + 8
	TagChecksumPosition            = TagShardSequencePosition + 2

	// Derivatives
	TagSize            = TagChecksumPosition + 4 // should be 24b
	DifferentiatorSize = TagRequiredShardsPosition - TagBlockDifferentiatorPosition
	// MaxPadding determines how large of a pad the tag supports.
	MaxPadding = (1 << (8 * (TagBlockSequencePosition - TagPaddingPosition))) - 1
)

// Tag holds the all the neccessary hints to perform full data reconstruction.
type Tag struct {
	Version             uint8                    // Encoder version used to create the tag.
	BlockDifferentiator [DifferentiatorSize]byte // Random block group UUID for identifying blocks that belong together.
	RequiredShards      uint8                    // Number of valid shards required for restoration.
	RedundantShards     uint8                    // Number of additional shards that can be used in place of invalid shards.
	Padding             uint16                   // Number of bytes to discard after restoration. Typically zero, except for the very last block.
	BlockSequence       uint64
	ShardSequence       uint16
}

// Differentiate fills the block differentiator with random bytes.
func (t *Tag) Differentiate() (err error) {
	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(t.BlockDifferentiator[:])
	return
}

func (t *Tag) Write(b []byte) (n int, err error) {
	if len(b) < TagSize {
		return 0, io.ErrShortBuffer
	}
	b[0] = byte(t.Version)
	copy(b[TagBlockDifferentiatorPosition:TagRequiredShardsPosition], t.BlockDifferentiator[:])
	b[TagRequiredShardsPosition] = byte(t.RequiredShards)
	b[TagRedundantShardsPosition] = byte(t.RedundantShards)
	b[TagPaddingPosition] = byte(t.Padding)
	binary.BigEndian.PutUint64(b[TagBlockSequencePosition:TagShardSequencePosition], t.BlockSequence)
	binary.BigEndian.PutUint16(b[TagShardSequencePosition:TagChecksumPosition], t.ShardSequence)
	return 0, nil
}

// Prototype creates a clonable prototype.
func (t *Tag) Prototype() TagPrototype {
	var result TagPrototype
	t.Write(result[:])
	return result
}
