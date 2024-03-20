package telomeres

import (
	"context"
	"errors"
	"io"
)

var (
	ErrBoundary       = errors.New("encountered telomere boundary")
	ErrUnpairedEscape = errors.New("unpaired escaped byte")
)

// Decoder reads until a [Mark] byte and strips [Escape] bytes.
type Decoder struct {
	r            io.ReadSeeker
	telomereTail int64
}

// NewDecoder sets up the decoder.
func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{r: r}
}

// StreamChunk is a calls [Decoder.StreamChunkBuffer] with default buffer.
func (d *Decoder) StreamChunk(ctx context.Context, to io.Writer) (written int64, err error) {
	return d.StreamChunkBuffer(ctx, to, d.makeDefaultBuffer())
}

// StreamChunkBuffer decodes the underlying [io.ReadSeeker] into
// a writer using a specified buffer. Context is checked for expiration
// between writes. [ErrBoundary] ends current chunk. Call again
// to get the next chunk. Returns io.EOF if there are no more chunks.
// Buffer must have capacity for **at least two bytes**.
func (d *Decoder) StreamChunkBuffer(ctx context.Context, to io.Writer, buffer []byte) (written int64, err error) {
	var n int
	var werr error
	for {
		n, err = d.Read(buffer)
		if n > 0 {
			if n, werr = to.Write(buffer[:n]); werr != nil {
				return written, errors.Join(err, werr)
			}
			written += int64(n)
		}

		select {
		case <-ctx.Done():
			return written, errors.Join(ctx.Err(), err)
		default:
		}

		if err != nil {
			if err == ErrBoundary {
				if written == 0 {
					continue
				}
				return written, nil
			}
			if err == io.EOF && written > 0 {
				return written, nil
			}
			return written, err
		}
	}
}

// Read satisfies [io.Reader] interface. Returns [ErrBoundary]
// if a [Mark] byte is detected. Return [ErrUnpairedEscape]
// if an escaped character is not paired with another.
func (d *Decoder) Read(b []byte) (n int, err error) {
	var (
		c         byte
		index     int
		lastIndex int
	)
	n, err = d.r.Read(b)
	if n < 1 {
		return n, err
	}
	window := b[:n]
	if d.telomereTail > 0 && window[0] != Mark {
		d.telomereTail = 0
	}
	// log.Printf("in: %q", b[:n])
	// defer func() {
	// 	log.Printf("out: %q %d", b[:n], n)
	// }()

decode:
	for index, c = range window {
		switch c {
		case Mark:
			n = n - len(window) + index
			window = window[index+1:]
			// log.Printf("window: %q buffer: %q", string(window), b[:n])
			// log.Printf("int: %q %d %d", b[:n], n, n+len(window)+index)
			goto drain
		case Escape:
			n--
			lastIndex = len(window) - 1
			if index == lastIndex {
				if err == nil {
					_, err = d.r.Seek(-1, io.SeekCurrent)
					return n, err
				}
				err = ErrUnpairedEscape
				break decode
			}
			copy(window[index:lastIndex], window[index+1:]) // cut current byte
			window = window[index+1 : lastIndex]
			// log.Println("escaped window:", string(window), string(b[:n]))
			goto decode
		}
	}
	return n, err

drain: // discard any remaining [Mark] bytes
	d.telomereTail++ // for the previous byte that got us to drain
	for index, c = range window {
		if c != Mark {
			d.telomereTail += int64(index)
			// log.Printf("seeking back: %d %q %q", -int64(len(window)-index), string(window), c)
			// log.Printf("buffer: %q", string(b[n:]))
			_, err = d.r.Seek(-int64(len(window)-index), io.SeekCurrent)
			if err != nil {
				return n, err
			}
			return n, ErrBoundary
		}
	}
	d.telomereTail += int64(len(window))
	return n, ErrBoundary
}

func (d *Decoder) makeDefaultBuffer() []byte {
	size := 32 * 1024
	if l, ok := d.r.(io.Reader).(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	return make([]byte, size)
}
