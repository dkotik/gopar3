package telomeres

import (
	"errors"
	"io"
)

// ErrTelomereBoundaryReached represents a state of encountering a repeated telomereMarkByte.
var ErrTelomereBoundaryReached = errors.New("reach a telomere boundary")

// TelomereStreamDecoder reads the stream, strips telomereEscapeBytes, and detects telomere boundaries.
type TelomereStreamDecoder struct {
	l         int
	b, window []byte
	r         io.Reader
}

// NewTelomereStreamDecoder sets up the decoder. Telomere length determines how many telomereMarkByte are found in a sequence before ErrTelomereBoundaryReached state is triggered.
func NewTelomereStreamDecoder(r io.Reader, telomereLength, bufferSize int) *TelomereStreamDecoder {
	return &TelomereStreamDecoder{
		l: telomereLength,
		b: make([]byte, bufferSize),
		r: r,
	}
}

func (t *TelomereStreamDecoder) fillBuffer() (err error) {
	n := copy(t.b[:], t.window[:])

	var j int
	var l = len(t.b)
	for {
		j, err = t.r.Read(t.b[n:])
		n += j
		if err != nil || n == l {
			break
		}
	}
	t.window = t.b[:n]
	if n > 0 && err == io.EOF {
		err = nil
	}
	return
}

func (t *TelomereStreamDecoder) Read(b []byte) (n int, err error) {
	if len(t.window) < t.l { // try to refill when close
		if err = t.fillBuffer(); err != nil {
			return 0, err
		}
	}

	// try to detect a boundary
	boundary := 0
	for i := 0; i < len(t.window); i++ {
		if t.window[i] == telomereMarkByte {
			boundary++
		} else {
			break
		}
	}
	if boundary >= t.l {
		t.window = t.window[boundary:]
		return 0, ErrTelomereBoundaryReached
	}

	escaped := false
	var i int
loop:
	for i = 0; i < len(t.window); i++ {
		if escaped {
			escaped = false
		} else {
			switch t.window[i] {
			case telomereEscapeByte:
				boundary = 0
				escaped = true
				continue
			case telomereMarkByte:
				b[n] = t.window[i]
				n++
				if boundary++; boundary == t.l {
					break loop
				}
				continue
			default:
				boundary = 0
			}
		}
		b[n] = t.window[i]
		n++
	}
	t.window = t.window[i-boundary:]
	return n - boundary, nil
}
