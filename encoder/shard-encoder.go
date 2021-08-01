package encoder

import (
	"hash"
	"io"

	"github.com/dkotik/gopar3/shard"
)

const (
	// ShardLimit is constrained by klauspost/reedsolomon limit.
	ShardLimit = 256
)

// ShardEncoder commits and tags given shards to IO streams.
type ShardEncoder struct {
	w   io.Writer
	c   hash.Hash32
	tag shard.TagPrototype
}

// NewShardEncoder creates the recorded and populates it with reasonable defaults. If given check sum algorithm if nil, falls back on CRC32.
func NewShardEncoder(w io.Writer, checkSum hash.Hash32, prefill *shard.Tag) (s *ShardEncoder) {
	if checkSum == nil {
		checkSum = shard.NewChecksum()
	}
	s = &ShardEncoder{
		w: w,
		c: checkSum,
	}

	if prefill == nil {
		prefill = &shard.Tag{}
	}
	// if prefill.Version == 0 {
	// 	prefill.Version = VersionByte
	// }
	prefill.Write(s.tag[:])
	return
}
