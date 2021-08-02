package encoder

import (
	"bytes"
	"errors"
	"io"

	"github.com/dkotik/gopar3/shard"
	"github.com/dkotik/gopar3/telomeres"
)

var (
	// ErrShardDidNotFit triggers when there are not enough bytes to write the entire shard.
	ErrShardDidNotFit = errors.New("failed to write the entire shard")
	// ErrShardTagDidNotFit triggers when there are not enough bytes to write the entire shard tag.
	ErrShardTagDidNotFit = errors.New("failed to write the entire shard tag")
)

// Writer streams encoded shards into an underlying writer.
type Writer struct {
	TagPrototype shard.TagPrototype
	t            *telomeres.Encoder
	c            *checkSumWriter
}

// NewWriter sets up a nested encoder for commiting shards to IO.
func (e *Encoder) NewWriter(w io.Writer, proto shard.TagPrototype) *Writer {
	telw := telomeres.NewEncoder(w, e.telomeresLength, e.telomeresBufferSize)
	return &Writer{
		TagPrototype: proto,
		t:            telw,
		c:            &checkSumWriter{telw, shard.NewChecksum()},
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
	if int(j) != shard.TagSize {
		return n, ErrShardTagDidNotFit
	}

	c := shard.NewChecksum()
	if _, err = w.c.Write(c.Sum(b)); err != nil {
		return n, err
	}

	return w.t.Cut()
	// TODO: add periodic checksum
}
