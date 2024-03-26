package gopar3

import (
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"math"
)

const (
	TagBytesForCRC         = 4
	TagBytesForSourceSize  = 8
	TagBytesForShardQuorum = 1
	TagBytesForShardOrder  = 1
	TagBytesForShardBatch  = 2

	TagBeginSourceCRC   = 0
	TagEndSourceCRC     = TagBeginSourceCRC + TagBytesForCRC
	TagBeginSourceSize  = TagEndSourceCRC
	TagEndSourceSize    = TagBeginSourceSize + TagBytesForSourceSize
	TagBeginShardQuorum = TagEndSourceSize
	TagEndShardQuorum   = TagBeginShardQuorum + TagBytesForShardQuorum
	TagBeginShardOrder  = TagEndShardQuorum
	TagEndShardOrder    = TagBeginShardOrder + TagBytesForShardOrder
	TagBeginShardBatch  = TagEndShardOrder
	TagEndShardBatch    = TagBeginShardBatch + TagBytesForShardBatch
	TagSize             = TagEndShardBatch - TagBeginSourceCRC
	DifferentiatorSize  = TagEndShardQuorum - TagBeginSourceCRC
)

// Tag holds the parameters to perform validated data reconstruction.
type Tag struct {
	SourceCRC   uint32
	SourceSize  uint64
	ShardQuorum uint8
	ShardOrder  uint8
	ShardBatch  uint16
}

func NewTag(
	ctx context.Context,
	r io.Reader,
	quorum uint8,
) (tag Tag, err error) {
	tag.ShardQuorum = quorum
	// 32 * 1024 is io package default
	b := make([]byte, 4*32*1024)
	crc := crc32.New(castagnoliTable)
	n := 0

	for err == nil {
		select {
		case <-ctx.Done():
			return tag, ctx.Err()
		default:
		}
		n, err = io.ReadFull(r, b)
		tag.SourceSize += uint64(n)
		_, _ = crc.Write(b[:n])
	}
	switch err {
	case io.EOF, io.ErrUnexpectedEOF:
		err = nil
	default:
		return
	}
	tag.SourceCRC = crc.Sum32()
	// panic(tag.SourceCRC)
	return
}

func NewTagFromBytes(b []byte) Tag {
	return Tag{
		SourceCRC: binary.BigEndian.Uint32(
			b[TagBeginSourceCRC:TagEndSourceCRC],
		),
		SourceSize: binary.BigEndian.Uint64(
			b[TagBeginSourceSize:TagEndSourceSize],
		),
		ShardQuorum: b[TagBeginShardQuorum],
		ShardOrder:  b[TagBeginShardOrder],
		ShardBatch: binary.BigEndian.Uint16(
			b[TagBeginShardBatch:TagEndShardBatch],
		),
	}
}

// Bytes encodes the tag into binary format.
func (t Tag) Bytes() (b []byte) {
	b = make([]byte, TagSize)
	binary.BigEndian.PutUint32(
		b[TagBeginSourceCRC:TagEndSourceCRC],
		t.SourceCRC,
	)
	binary.BigEndian.PutUint64(
		b[TagBeginSourceSize:TagEndSourceSize],
		t.SourceSize,
	)
	b[TagBeginShardQuorum] = t.ShardQuorum
	b[TagBeginShardOrder] = t.ShardOrder
	binary.BigEndian.PutUint16(
		b[TagBeginShardBatch:TagEndShardBatch],
		t.ShardBatch,
	)
	return b
}

type Tagger interface {
	Bytes() []byte
	Next() error
}

type sequentialTagger struct {
	encoded    []byte
	tag        Tag
	shardLimit uint8
}

// NewSequentialTagger prepares a tagger that increments
// shard counter to the limit. Then, it sets the shard counter
// to zero and increments shard batch counter.
func NewSequentialTagger(t Tag, shardLimit uint8) Tagger {
	return &sequentialTagger{
		encoded:    t.Bytes(),
		tag:        t,
		shardLimit: shardLimit,
	}
}

func (t *sequentialTagger) Bytes() []byte {
	// log.Printf("tag shard %d.%d", t.tag.ShardBatch, t.tag.ShardOrder)
	return t.encoded
}

func (t *sequentialTagger) Next() error {
	if t.tag.ShardBatch == math.MaxUint16 && t.tag.ShardOrder == math.MaxUint8 {
		return errors.New("too many shards")
	}
	t.tag.ShardOrder++
	if t.tag.ShardOrder < t.shardLimit {
		t.encoded[TagBeginShardOrder] = t.tag.ShardOrder
	} else {
		t.tag.ShardOrder = 0
		t.tag.ShardBatch++
		t.encoded[TagBeginShardOrder] = 0
		binary.BigEndian.PutUint16(
			t.encoded[TagBeginShardBatch:TagEndShardBatch],
			t.tag.ShardBatch,
		)
	}
	return nil
}

type latteralTagger struct {
	encoded []byte
	tag     Tag
}

// NewLateralTagger prepares a tagger that increments
// shard batch counter. Useful for writing data into
// separate shard files.
func NewLateralTagger(t Tag) Tagger {
	return &latteralTagger{tag: t}
}

func (t *latteralTagger) Bytes() []byte {
	return t.encoded
}

func (t *latteralTagger) Next() error {
	if t.tag.ShardBatch == math.MaxUint16 {
		return errors.New("too many shards")
	}
	t.tag.ShardBatch++
	binary.BigEndian.PutUint16(
		t.encoded[TagBeginShardBatch:TagEndShardBatch],
		t.tag.ShardBatch,
	)
	return nil
}
