package gopar3

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"

	"github.com/dkotik/gopar3/telomeres"
)

var (
	ErrShardTooSmall = errors.New("there are not enough bytes to decode the shard checksum and tag")
)

type CheckSumError struct {
	Shard         *Shard
	CastagnoliSum uint32
}

func (e *CheckSumError) Error() string {
	return fmt.Sprintf("corrupted shard: Castagnoli CRC32 sum %d does not match %d", e.Shard.CastagnoliSum, e.CastagnoliSum)
}

type Reader struct {
	*telomeres.Decoder
	Source   string
	buffer   []byte
	shardCRC hash.Hash32
}

func NewReader(source string, r io.ReadSeeker) *Reader {
	return &Reader{
		Source:   source,
		Decoder:  telomeres.NewDecoder(r),
		buffer:   make([]byte, 32*1024),
		shardCRC: crc32.New(castagnoliTable),
	}
}

func (r *Reader) NextShard(ctx context.Context, w io.Writer) (s *Shard, err error) {
	s = &Shard{
		Source: r.Source,
	}
	defer func() {
		var cerr error
		if s.LastByte, cerr = r.Decoder.Cursor(); err != nil {
			err = errors.Join(err, cerr)
		}
		if err != nil {
			s.Error = err.Error()
		} else if realSum := r.shardCRC.Sum32(); realSum != s.CastagnoliSum {
			s.Error = (&CheckSumError{Shard: s, CastagnoliSum: realSum}).Error()
		}
	}()
	if err = r.Decoder.SeekChunk(ctx); err != nil {
		return s, err
	}
	s.FirstByte, err = r.Decoder.Cursor()
	if err != nil {
		return s, err
	}

	var werr error
	n, err := io.ReadFull(r.Decoder, r.buffer)
	if n < TagSize+TagBytesForCRC {
		return s, ErrShardTooSmall
	}
	s.Size += int64(n)
	r.shardCRC.Reset()
	s.CastagnoliSum = binary.BigEndian.Uint32(r.buffer[:TagBytesForCRC])
	if _, werr = r.shardCRC.Write(r.buffer[TagBytesForCRC:n]); werr != nil {
		return s, werr
	}
	s.Tag = NewTagFromBytes(r.buffer[TagBytesForCRC : TagBytesForCRC+TagSize])
	switch err {
	case nil:
	case io.EOF, io.ErrUnexpectedEOF, telomeres.ErrBoundary:
		if n, err = w.Write(r.buffer[TagBytesForCRC+TagSize : n]); err != nil {
			return s, err
		}
		fallthrough
	default:
		return s, err
	}

	n, err = w.Write(r.buffer[TagBytesForCRC+TagSize : n])
	// log.Fatalf("%s", r.buffer[TagBytesForCRC+TagSize:n])
	if err != nil {
		return s, err
	}

	for err == nil {
		select {
		case <-ctx.Done():
			return s, ctx.Err()
		default:
		}
		n, err = io.ReadFull(r.Decoder, r.buffer)
		if n > 0 {
			s.Size += int64(n)
			if _, werr = w.Write(r.buffer[:n]); werr != nil {
				return s, werr
			}
			if _, werr = r.shardCRC.Write(r.buffer[:n]); werr != nil {
				return s, werr
			}
		}
	}

	switch err {
	case telomeres.ErrBoundary:
		return s, nil
	default:
		return s, err
	}
}
