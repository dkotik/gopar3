package gopar3

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"sort"
	"time"
)

const (
	shardTagBlockDifferentiatorPosition = 1
	shardTagRequiredShardsPosition      = shardTagBlockDifferentiatorPosition + 4
	shardTagRedundantShardsPosition     = shardTagRequiredShardsPosition + 1
	shardTagPaddingPosition             = shardTagRedundantShardsPosition + 1
	shardTagBlockSequencePosition       = shardTagPaddingPosition + 2
	shardTagShardSequencePosition       = shardTagBlockSequencePosition + 8
	shardTagChecksumPosition            = shardTagShardSequencePosition + 2
	shardTagSize                        = shardTagChecksumPosition + 4 // should be 24b
)

// ShardTag holds the all the neccessary hints to perform full data reconstruction.
type ShardTag struct {
	Version             uint8   // Encoder version used to create the tag.
	BlockDifferentiator [4]byte // Random block group UUID for identifying blocks that belong together.
	RequiredShards      uint8   // Number of valid shards required for restoration.
	RedundantShards     uint8   // Number of additional shards that can be used in place of invalid shards.
	Padding             uint16  // Number of bytes to discard after restoration. Typically zero, except for the very last block.
	BlockSequence       uint64
	ShardSequence       uint16
}

// Differentiate fills the block differentiator with random bytes.
func (t *ShardTag) Differentiate() (err error) {
	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(t.BlockDifferentiator[:])
	return
}

func (t *ShardTag) Write(b []byte) (n int, err error) {
	if len(b) < shardTagSize {
		return 0, errors.New("not enough bytes given for a shard tag")
	}
	b[0] = byte(t.Version)
	copy(b[shardTagBlockDifferentiatorPosition:shardTagRequiredShardsPosition], t.BlockDifferentiator[:])
	b[shardTagRequiredShardsPosition] = byte(t.RequiredShards)
	b[shardTagRedundantShardsPosition] = byte(t.RedundantShards)
	b[shardTagPaddingPosition] = byte(t.Padding)
	binary.BigEndian.PutUint64(b[shardTagBlockSequencePosition:shardTagShardSequencePosition], t.BlockSequence)
	binary.BigEndian.PutUint16(b[shardTagShardSequencePosition:shardTagChecksumPosition], t.ShardSequence)
	return 0, nil
}

func medianUint8(bunch []uint8) uint8 {
	sort.Slice(bunch, func(i int, j int) bool {
		return bunch[i] > bunch[j]
	})
	return bunch[len(bunch)/2]
}

func medianUint16(bunch []uint16) uint16 {
	sort.Slice(bunch, func(i int, j int) bool {
		return bunch[i] > bunch[j]
	})
	return bunch[len(bunch)/2]
}
