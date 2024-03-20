package telomeres

import (
	"context"
	"errors"
	"io"
)

// Cursor returns the underlying reader position.
// Data chunk boundary can be determined by checking the
// cursor position before and after calling [Decoder.Stream].
func (d *Decoder) Cursor() (n int64, err error) {
	n, err = d.r.Seek(0, io.SeekCurrent)
	n -= d.telomereTail
	return
}

func (d *Decoder) SeekChunk(ctx context.Context) error {
	return d.SeekChunkBuffer(ctx, d.makeDefaultBuffer())
}

func (d *Decoder) SeekChunkBuffer(
	ctx context.Context,
	buffer []byte,
) (err error) {
	var (
		n int
		i int
		c byte
	)

	for {
		n, err = d.r.Read(buffer)
		for i, c = range buffer[:n] {
			if c != Mark {
				_, err = d.r.Seek(-int64(n-i), io.SeekCurrent)
				return err
			}
		}

		select {
		case <-ctx.Done():
			return errors.Join(ctx.Err(), err)
		default:
		}

		if err != nil {
			return err
		}
	}
}
