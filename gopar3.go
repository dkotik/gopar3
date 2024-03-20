package gopar3

import (
	"errors"
	"hash"
	"hash/crc32"
	"io"
)

const (
	PaddingByte = '?'

	ShardLimit      = 1<<(TagBytesForShardOrder*8) - 1
	ShardGroupLimit = 1<<(TagBytesForShardGroup*8) - 1
	SourceSizeLimit = 1<<(TagBytesForSourceSize*8) - 1
)

// castagnoliTable sources [crc.New] with 0x82f63b78
// polynomial. It is known for superior error detection
// and use for BitTorrent and iSCSI protocols.
var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

type crcWriteCloser struct {
	io.Writer
	hash hash.Hash32
}

// NewCRC wraps a writer with crc32.Castagnoli that
// is added when closed.
func NewCRC(w io.Writer) io.WriteCloser {
	return &crcWriteCloser{
		Writer: w,
		hash:   crc32.New(castagnoliTable),
	}
}

func (w *crcWriteCloser) Write(b []byte) (n int, err error) {
	n, err = w.Writer.Write(b)
	_, _ = w.hash.Write(b[:n])
	return
}

func (w *crcWriteCloser) Close() (err error) {
	n, err := w.Writer.Write(w.hash.Sum(nil))
	if err != nil {
		return err
	}
	if n != TagBytesForCRC {
		return errors.New("failed to write checksum")
	}
	return
}
