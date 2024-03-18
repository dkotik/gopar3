package telomeres

import (
	"errors"
	"fmt"
	"io"
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
	r                io.Reader
}

// NewDecoder sets up the decoder. Telomere length determines how many telomereMarkByte are found in a sequence before ErrTelomereBoundaryReached state is triggered.
func (t *Telomeres) NewDecoder(
	r io.Reader,
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

// Cursor returns the underlying reader position.
// Data chunk boundary can be determined by checking the
// cursor position before and after calling [Decoder.WriteTo].
func (d *Decoder) Cursor() (position int64) {
	return d.cursor - int64(d.telomereStreak)
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

func (d *Decoder) FindEdge(softReadLimit int64) (n int64, err error) {
	var (
		i int
		c byte
	)

	if len(d.window) < d.minimum {
		n, err = d.fillBufferWindow()
		softReadLimit -= n
		n = 0
	}

	for i, c = range d.window {
		if c == d.mark {
			n++
			continue
		}
		if d.telomereStreak+i >= d.minimum {
			// found edge that ended a streak
			d.window = nil
			d.telomereStreak = 0
			d.cursor += int64(i)
			return n, err
		}
		return n, err // do not move cursor - already at edge
	}

	more := 0
	for {
		if err != nil {
			break
		}
		more, err = d.r.Read(d.b)
		for i, c = range d.b[:more] {
			if c != d.mark {
				d.window = d.b[i:]
				d.telomereStreak += i
				d.cursor += int64(i)
				return n, err
			}
			n++
		}
		d.cursor += int64(more)
		softReadLimit -= int64(more)
		// TODO: add soft limit check
	}

	d.window = nil // exhausted read
	d.telomereStreak += int(n)
	d.cursor += n
	return n, err
}
