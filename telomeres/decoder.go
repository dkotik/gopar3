package telomeres

import (
	"errors"
	"io"
)

// ErrEndReached represents a state of encountering a repeated telomereMarkByte.
var ErrEndReached = errors.New("reader has no more data")

// Decoder reads the stream, strips telomereEscapeBytes, and detects telomere boundaries.
type Decoder struct {
	minimum, maximum, boundary int
	b, window                  []byte
	r                          io.Reader
}

// NewDecoder sets up the decoder. Telomere length determines how many telomereMarkByte are found in a sequence before ErrTelomereBoundaryReached state is triggered.
func NewDecoder(
	r io.Reader,
	minimumTelomereLength int,
	maximumTelomereLength int,
	bufferSize int,
) *Decoder {
	if minimumTelomereLength == 0 {
		panic("minimum telomere length cannot be 0")
	}
	if maximumTelomereLength == 0 {
		panic("maximum telomere length cannot be 0")
	}
	if minimumTelomereLength > maximumTelomereLength {
		panic("minimum telomere length cannot be greater than the maximum")
	}
	if minimumTelomereLength >= bufferSize {
		panic("telomere length must be less than the allocated buffer size")
	}
	return &Decoder{
		minimum: minimumTelomereLength,
		maximum: maximumTelomereLength,
		b:       make([]byte, bufferSize),
		r:       r,
	}
}

func (t *Decoder) fillBuffer() (err error) {
	n := copy(t.b[:], t.window[:]) // move remainder to the front

	var j int
	var max = len(t.b)
	for {
		j, err = t.r.Read(t.b[n:])
		n += j
		if err != nil || n == max {
			break
		}
	}
	t.window = t.b[:n] // rebuild window
	return
}

// Skip keeps reading telomereMarkBytes until something else is found. Returns the number of telomereMarkBytes read and an error.
func (t *Decoder) Skip() (boundary int, err error) {
	// var i int
	for {
		if len(t.window) < t.minimum {
			// try to refill the buffer when window gets small
			if err = t.fillBuffer(); err != nil {
				if err == io.EOF {
					return 0, nil // swallow io.EOF
				}
				return 0, err
			}
		}

		// try to detect a boundary and keep going if you do
		i := 0
		for ; i < len(t.window); i++ {
			if t.window[i] == telomereMarkByte {
				boundary++
				if boundary > t.maximum {
					t.window = t.window[i:]
					return 0, errors.New("maximum telomere length exceeded")
				}
			} else {
				t.window = t.window[i:]
				return boundary, nil
			}
		}
		t.window = t.window[i:]
	}
}

func (t *Decoder) Read(b []byte) (n int, err error) {
	boundary, err := t.Skip()
	if err != nil {
		return 0, err
	}
	boundary += t.boundary // add the number from previous iteration

	if boundary >= t.minimum {
		t.boundary = 0
		return 0, io.EOF
	}
	if len(t.window) == 0 {
		// no more data is coming
		return 0, ErrEndReached
	}

	// write the boundary back into b, because it was discovered to not be long enough
	for ; boundary > 0 && n < len(b); boundary-- {
		b[n] = telomereMarkByte
		n++
	}

	var i int
	escaped := false
loop:
	for ; i < len(t.window) && n < len(b); i++ {
		if escaped {
			escaped = false
		} else {
			switch t.window[i] {
			case telomereEscapeByte:
				boundary = 0
				escaped = true
				continue
			case telomereMarkByte:
				if boundary++; boundary >= t.minimum {
					n++ // count this byte
					break loop
				}
			default:
				boundary = 0
			}
		}
		b[n] = t.window[i]
		n++
	}
	t.window = t.window[i:]
	t.boundary = boundary // remember the boundary
	return n - boundary, nil
}
