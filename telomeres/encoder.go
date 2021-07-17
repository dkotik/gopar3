package telomeres

import (
	"bytes"
	"io"
)

// TelomereStreamEncoder appends telomereEscapeByte to each telomereMarkByte. Do not forget to call WriteTelomere at the end or flush manually.
type TelomereStreamEncoder struct {
	t      []byte
	b      []byte
	cursor int
	w      io.Writer
}

// NewTelomereStreamEncoder sets up the encoder.
func NewTelomereStreamEncoder(w io.Writer, telomereLength, bufferSize int) *TelomereStreamEncoder {
	telomeres := make([]byte, telomereLength)
	for i := 0; i < telomereLength; i++ {
		telomeres[i] = telomereMarkByte
	}

	return &TelomereStreamEncoder{
		t: telomeres,
		b: make([]byte, bufferSize),
		w: w,
	}
}

// Flush commits the contents of the buffer to the underlying Writer.
func (t *TelomereStreamEncoder) Flush() (err error) {
	_, err = io.Copy(t.w, bytes.NewReader(t.b[:t.cursor]))
	t.cursor = 0
	return err
}

func (t *TelomereStreamEncoder) Write(b []byte) (n int, err error) {
	// one less for cursor, because may write two bytes per loop iteration
	max := len(t.b) - 1
	for ; n < len(b) && t.cursor < max; n++ {
		current := b[n]
		if current == telomereMarkByte || current == telomereEscapeByte {
			t.b[t.cursor] = telomereEscapeByte
			t.b[t.cursor+1] = current
			t.cursor += 2
		} else {
			t.b[t.cursor] = b[n]
			t.cursor++
		}
	}
	if t.cursor == max {
		if err = t.Flush(); err != nil {
			return 0, err
		}
	}
	return n, nil
}

// WriteTelomere flushes the buffer and writes repeated telomereEscapeBytes to the underlying Writer.
func (t *TelomereStreamEncoder) WriteTelomere() (n int, err error) {
	if err = t.Flush(); err != nil {
		return 0, err
	}
	j, err := io.Copy(t.w, bytes.NewReader(t.t))
	return int(j), err
}
