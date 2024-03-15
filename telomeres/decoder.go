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
				d.window = d.window[n+1:]
			}
		}

		// add carry-over boundary bytes, if any
		for i := 0; i < boundary; i++ {
			decoded = append(decoded, d.mark)
		}

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
			if err == io.EOF && total > 0 {
				return total, nil
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

// Cursor returns the underlying reader position with the
// number of detected telomere streak on the end of the buffer.
// Data chunk boundary can be determined by checking the
// cursor position before and after calling [Decoder.WriteTo].
func (d *Decoder) Cursor() (position int64, tailTelomeres int64) {
	return d.cursor, int64(d.telomereStreak)
}

func (d *Decoder) Next(w io.Writer) (n int64, err error) {
	// TODO: implement
	// TODO: update cursor
	// TODO: test that it works!
	return 0, nil
}
