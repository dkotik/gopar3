package gopar3

import (
	"hash"
	"hash/crc32"
	"io"

	"github.com/dkotik/gopar3/telomeres"
)

type writer struct {
	encoder *telomeres.Encoder
	tagger  Tagger
	crc     hash.Hash32
}

func NewWriter(w *telomeres.Encoder, t Tagger) (io.Writer, error) {
	if _, err := w.Cut(); err != nil {
		return nil, err
	}
	return &writer{
		encoder: w,
		tagger:  t,
		crc:     crc32.New(castagnoliTable),
	}, nil
}

// Write writes a Castagnoli sum, a shard tag, followed by given bytes
// to the [telomeres.Encoder]. Ends with a telomere sequence to designate
// the end of the shard.
func (w *writer) Write(b []byte) (n int, err error) {
	{ // write checksum
		w.crc.Reset()
		_, err = w.crc.Write(w.tagger.Bytes())
		if err != nil {
			return 0, err
		}
		_, err = w.crc.Write(b)
		if err != nil {
			return 0, err
		}
		n, err = w.encoder.Write(w.crc.Sum(nil))
		if err != nil {
			return 0, err
		}
		if n != TagBytesForCRC {
			return 0, io.ErrShortWrite
		}
	}

	{ // write tag
		n, err = w.encoder.Write(w.tagger.Bytes())
		if err != nil {
			return 0, err
		}
		if n != TagSize {
			return 0, io.ErrShortWrite
		}
	}

	n, err = w.encoder.Write(b)
	if err != nil {
		return n, err
	}
	if n != len(b) {
		return n, io.ErrShortWrite
	}
	if _, err = w.encoder.Cut(); err != nil {
		return n, err
	}
	return n, w.tagger.Next()
}
