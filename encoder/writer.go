package encoder

import (
	"bytes"
	"errors"
	"io"

	"github.com/dkotik/gopar3"
	"github.com/dkotik/gopar3/scanner"
	"github.com/dkotik/gopar3/telomeres"
)

var (
	// ErrShardDidNotFit triggers when there are not enough bytes to write the entire gopar3.
	ErrShardDidNotFit = errors.New("failed to write the entire shard")
	// ErrShardTagDidNotFit triggers when there are not enough bytes to write the entire shard tag.
	ErrShardTagDidNotFit = errors.New("failed to write the entire shard tag")
)

// Writer streams encoded shards into an underlying writer.
type Writer struct {
	TagPrototype gopar3.TagPrototype
	t            *telomeres.Encoder
	c            *checkSumWriter
}

// NewWriter sets up a nested encoder for commiting shards to IO.
func (e *Encoder) NewWriter(w io.Writer, proto gopar3.TagPrototype) *Writer {
	telw := telomeres.NewEncoder(w, e.telomeresLength, e.telomeresBufferSize)
	return &Writer{
		TagPrototype: proto,
		t:            telw,
		c:            &checkSumWriter{telw, scanner.NewChecksum()},
	}
}

func (w *Writer) Write(b []byte) (n int, err error) {
	j, err := io.Copy(w.c, bytes.NewReader(b))
	if err != nil {
		return
	}
	n = int(j)
	if n != len(b) {
		return n, ErrShardDidNotFit
	}
	return

	j, err = io.Copy(w.c, bytes.NewReader(w.TagPrototype[:]))
	if err != nil {
		return
	}
	if int(j) != gopar3.TagSize {
		return n, ErrShardTagDidNotFit
	}

	c := scanner.NewChecksum()
	if _, err = w.c.Write(c.Sum(b)); err != nil {
		return n, err
	}

	return w.t.Cut()
	// TODO: add periodic checksum
}
