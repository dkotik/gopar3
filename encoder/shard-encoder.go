package encoder

import (
	"encoding/binary"
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
func NewShardEncoder(w io.Writer, checkSum hash.Hash32, prefill *shard.ShardTag) (s *ShardEncoder) {
	if checkSum == nil {
		checkSum = shard.NewChecksum()
	}
	s = &ShardEncoder{
		w: w,
		c: checkSum,
	}

	if prefill == nil {
		prefill = &shard.ShardTag{}
	}
	// if prefill.Version == 0 {
	// 	prefill.Version = VersionByte
	// }
	prefill.Write(s.tag[:])
	return
}

func (s *ShardEncoder) Write(b []byte) (n int, err error) {
	n, err = s.w.Write(b)
	if n > 0 {
		s.c.Sum(b[:n])
	}
	return
}

// Seal writes tag and checksum.
func (s *ShardEncoder) Seal() (err error) {
	_, err = s.Write(s.tag[:])
	if err != nil {
		return
	}

	var checkSumBytes [4]byte
	binary.BigEndian.PutUint32(checkSumBytes[:], s.c.Sum32())
	s.c.Reset()
	_, err = s.w.Write(checkSumBytes[:])
	return
}
