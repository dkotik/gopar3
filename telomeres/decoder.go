package telomeres

import (
	"context"
	"errors"
	"fmt"
	"io"
)

var (
	ErrBoundary       = errors.New("encountered telomere boundary")
	ErrUnpairedEscape = errors.New("unpaired escaped byte")
)

// Decoder reads the stream, strips telomereEscapeBytes, and detects telomere boundaries.
type Decoder struct {
	mark   byte
	escape byte

	minimum          int
	telomereStreak   int
	b                []byte
	window           []byte
	maximumChunkSize int64
	cursor           int64
	r                io.ReadSeeker
}

// NewDecoder sets up the decoder. Telomere length determines how many telomereMarkByte are found in a sequence before ErrTelomereBoundaryReached state is triggered.
func (t *Telomeres) NewDecoder(
	r io.ReadSeeker,
) *Decoder {
	return &Decoder{
		mark:             t.mark,
		escape:           t.escape,
		minimum:          t.minimum,
		maximumChunkSize: 0,
		b:                make([]byte, t.bufferSize),
		r:                r,
	}
}

// StreamChunk is a calls [Decoder.StreamChunkBuffer] with default buffer.
func (d *Decoder) StreamChunk(ctx context.Context, to io.Writer) (written int64, err error) {
	b := make([]byte, 13) // TODO: update
	return d.StreamChunkBuffer(ctx, to, b)
}

// StreamChunkBuffer decodes the underlying [io.ReadSeeker] into
// a writer using a specified buffer. Context is checked for expiration
// between writes. [ErrBoundary] ends current chunk. Call again
// to get the next chunk. Returns io.EOF if there are no more chunks.
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
// if a telomere mark is detected. Return [ErrUnpairedEscape]
// if an escaped character is not paired with another.
func (d *Decoder) Read(b []byte) (n int, err error) {
	var (
		c         byte
		index     int
		lastIndex int
	)
	n, err = d.r.Read(b)
	window := b[:n]
	// log.Printf("in: %q", b[:n])
	// defer func() {
	// 	log.Printf("out: %q %d", b[:n], n)
	// }()

decode:
	for index, c = range window {
		switch c {
		case d.mark:
			n = n - len(window) + index
			window = window[index+1:]
			// log.Printf("window: %q buffer: %q", string(window), b[:n])
			// log.Printf("int: %q %d %d", b[:n], n, n+len(window)+index)
			goto drain
		case d.escape:
			n--
			lastIndex = len(window) - 1
			if index == lastIndex {
				if err == nil {
					d.cursor, err = d.r.Seek(-1, io.SeekCurrent)
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
	d.cursor += int64(n)
	return n, err

drain: // discard any remaining mark bytes
	for index, c = range window {
		if c != d.mark {
			// log.Printf("seeking back: %d %q %q", -int64(len(window)-index), string(window), c)
			// log.Printf("buffer: %q", string(b[n:]))
			d.cursor, err = d.r.Seek(-int64(len(window)-index), io.SeekCurrent)
			// panic("boo")
			if err != nil {
				return n, err
			}
			return n, ErrBoundary
		}
	}
	d.cursor += int64(n)
	return n, ErrBoundary
}

// WriteTo streams the next data chunk into a writer according
// to [io.WriterTo] expectations. Returns [io.EOF] if
// no more data chunks are coming.
func (d *Decoder) WriteTo(w io.Writer) (total int64, err error) {
	var (
		n          int
		windowSize int
		boundary   int
		escaped    bool
		decoded    []byte
	)

	for {
		windowSize = len(d.window)
		if windowSize < d.minimum {
			// TODO: use d.fillBufferWindow here?
			// total, err = d.fillBufferWindow()
			// not enough bytes to detect boundary, get more
			if windowSize > 0 {
				// move remaining bytes to the front of the buffer
				// source and destination may overlap
				copy(d.b[:windowSize], d.window)
				// fill the rest of the buffer
				n, err = d.r.Read(d.b[windowSize:])
				windowSize += n
			} else {
				n, err = d.r.Read(d.b)
				windowSize = n
			}
			d.window = d.b[:windowSize]
			d.cursor += int64(n)
			// d.cursor += total
			// total = 0 // reset, because variable was borrowed
		}

		if windowSize > 0 {
			var c byte
			if d.telomereStreak > 0 { // drain telomeres
				for n, c = range d.window {
					if c != d.mark { // end of telomeres
						// d.cursor += int64(d.telomereStreak)
						d.telomereStreak = 0
						d.window = d.window[n:]
						windowSize = len(d.window)
						break
					}
				}
				if d.telomereStreak > 0 {
					d.telomereStreak += n
					d.window = nil
					continue // possibly more coming
				}
			}

			if windowSize > 0 {
				decoded = make([]byte, 0, windowSize)
			decoding:
				for n, c = range d.window {
					if escaped {
						decoded = append(decoded, c)
						escaped = false
						continue
					}

					switch c {
					case d.escape:
						boundary = 0
						escaped = true
					case d.mark:
						// do not append boundary bytes
						// because they will either
						// trigger a telomere streak
						// or carry over until the next read
						boundary++
						if boundary >= d.minimum {
							d.telomereStreak = boundary
							boundary = 0
							break decoding
						}
					default:
						boundary = 0
						decoded = append(decoded, c)
					}
				}
				// shorten the window, except
				// for any carry-over boundary
				// bytes, because they could
				// complete a minimum length
				// requirement on the next read
				d.window = d.window[n+1-boundary:]
				// d.window = d.window[n+1:]
			}
		}

		// add carry-over boundary bytes, if any
		// for i := 0; i < boundary; i++ {
		// 	decoded = append(decoded, d.mark)
		// }
		// fmt.Println("carry over", string(decoded))
		// fmt.Println("window", string(d.window))

		if len(decoded) > 0 {
			// fmt.Println("decoded:", string(decoded), "?", string(d.window))
			n, werr := w.Write(decoded)
			total += int64(n)
			d.cursor += int64(n)
			if werr != nil {
				if err == io.EOF {
					return total, werr
				}
				return total, errors.Join(err, werr)
			}
		}

		if err != nil {
			if err == io.EOF {
				if len(d.window) > 0 {
					// some telomere bytes remaining
					// that were not yet written
					// because the window was shortened
					// write them now
					n, werr := w.Write(d.window)
					total += int64(n)
					d.cursor += int64(n)
					d.window = d.window[n:]
					if werr != nil {
						return total, werr
					}
				}

				if total > 0 {
					return total, nil
				}
			}
			return total, err
		}

		if total > 0 {
			if d.telomereStreak > 0 {
				// finished chunk
				return total, nil
			}
			if d.maximumChunkSize > 0 && total > d.maximumChunkSize {
				// TODO: add give-up limit for blocks too long here
				return total, fmt.Errorf("reached the maximum chunk size %d, giving up decoding, use Next(io.Discard) to dump bytes until the next telomere boundary", d.maximumChunkSize)
			}
		}
	}
}

func (d *Decoder) fillBufferWindow() (n int64, err error) {
	windowSize := len(d.window)
	if windowSize > 0 {
		// move remaining bytes to the front of the buffer
		// source and destination may overlap
		copy(d.b[:windowSize], d.window)
	}

	more := 0
	max := len(d.b)
	for {
		more, err = d.r.Read(d.b[windowSize:])
		if more == 0 {
			d.window = d.b[:windowSize]
			return
		}
		windowSize += more
		n += int64(more)
		if err != nil || windowSize >= max {
			d.window = d.b[:windowSize]
			return
		}
	}
}
