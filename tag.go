package gopar3

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	TagBytesForCRC         = 4
	TagBytesForSourceSize  = 8
	TagBytesForShardQuorum = 1
	TagBytesForShardOrder  = 1
	TagBytesForShardGroup  = 2

	// Derivatives:

	TagBeginSourceCRC   = 0
	TagEndSourceCRC     = TagBeginSourceCRC + TagBytesForCRC
	TagBeginSourceSize  = TagEndSourceCRC
	TagEndSourceSize    = TagBeginSourceSize + TagBytesForSourceSize
	TagBeginShardQuorum = TagEndSourceSize
	TagEndShardQuorum   = TagBeginShardQuorum + TagBytesForShardQuorum
	TagBeginShardOrder  = TagEndShardQuorum
	TagEndShardOrder    = TagBeginShardOrder + TagBytesForShardOrder
	TagBeginShardGroup  = TagEndShardOrder
	TagEndShardGroup    = TagBeginShardGroup + TagBytesForShardGroup
	TagSize             = TagEndShardGroup - TagBeginSourceCRC
	DifferentiatorSize  = TagEndShardQuorum - TagBeginSourceCRC
)

// Tag holds the parameters to perform validated data reconstruction.
type Tag struct {
	SourceCRC   uint32
	SourceSize  uint64
	ShardQuorum uint8
	ShardOrder  uint8
	ShardGroup  uint16
}

func NewTag(b []byte) Tag {
	return Tag{
		SourceCRC: binary.BigEndian.Uint32(
			b[TagBeginSourceCRC:TagEndSourceCRC],
		),
		SourceSize: binary.BigEndian.Uint64(
			b[TagBeginSourceSize:TagEndSourceSize],
		),
		ShardQuorum: b[TagBeginShardQuorum],
		ShardOrder:  b[TagBeginShardOrder],
		ShardGroup: binary.BigEndian.Uint16(
			b[TagBeginShardGroup:TagEndShardGroup],
		),
	}
}

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
		b[TagBeginShardGroup:TagEndShardGroup],
		t.ShardGroup,
	)
	return b
}

type tagWriter struct {
	w          io.Writer
	encoded    []byte
	tag        Tag
	shardLimit uint8
}

func (tw *tagWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if tw.tag.ShardGroup == math.MaxUint16 && tw.tag.ShardOrder == math.MaxUint8 {
		return 0, errors.New("too many shards")
	}
	n, err = io.Copy(tw.w, r)
	if err != nil {
		return n, err
	}
	tn, err := tw.w.Write(tw.encoded)
	n += int64(tn)
	if tw.tag.ShardOrder < tw.shardLimit {
		tw.tag.ShardOrder++
		tw.encoded[TagBeginShardOrder] = tw.tag.ShardOrder
	} else {
		tw.tag.ShardOrder = 0
		tw.tag.ShardGroup++
		tw.encoded[TagBeginShardOrder] = 0
		binary.BigEndian.PutUint16(
			tw.encoded[TagBeginShardGroup:TagEndShardGroup],
			tw.tag.ShardGroup,
		)
	}
	// TODO: write checksum?
	return n, err
}

func NewTagWriter(to io.Writer, tag Tag, parityShards uint8) io.ReaderFrom {
	return &tagWriter{
		w:          to,
		encoded:    nil,
		tag:        tag,
		shardLimit: tag.ShardQuorum + parityShards,
	}
}
