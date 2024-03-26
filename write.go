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

func (w *writer) Write(b []byte) (n int, err error) {
	var pn int
	// write body
	n, err = w.encoder.Write(b)
	if err != nil {
		return n, err
	}
	if n != len(b) {
		return n, io.ErrShortWrite
	}

	// write tag
	tag := w.tagger.Bytes()
	pn, err = w.encoder.Write(tag)
	if err != nil {
		return n, err
	}
	if pn != TagSize {
		return n, io.ErrShortWrite
	}

	// write checksum
	w.crc.Reset()
	_, err = w.crc.Write(b)
	if err != nil {
		return n, err
	}
	_, err = w.crc.Write(tag)
	if err != nil {
		return n, err
	}
	pn, err = w.encoder.Write(w.crc.Sum(nil))
	if err != nil {
		return n, err
	}
	if pn != TagBytesForCRC {
		return n, io.ErrShortWrite
	}

	if _, err = w.encoder.Cut(); err != nil {
		return n, err
	}
	return n, w.tagger.Next()
}
