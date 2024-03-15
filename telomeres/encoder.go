package telomeres

import (
	"bytes"
	"io"
)

// Encoder appends an escape byte to each mark byte.
// Call [Encoder.Cut] manually at the begining of writing,
// after copying each data chunk, and at the end.
type Encoder struct {
	mark   byte
	escape byte
	t      []byte
	b      []byte
	cursor int
	w      io.Writer
}

// NewEncoder creates a telomere encoder.
func (t *Telomeres) NewEncoder(w io.Writer) *Encoder {
	telomeres := make([]byte, t.minimum)
	for i := range telomeres {
		telomeres[i] = t.mark
	}

	return &Encoder{
		mark:   t.mark,
		escape: t.escape,
		t:      telomeres,
		b:      make([]byte, t.bufferSize),
		w:      w,
	}
}

// Flush commits the contents of the buffer to the underlying Writer.
func (t *Encoder) Flush() (err error) {
	if _, err = io.Copy(t.w, bytes.NewReader(t.b[:t.cursor])); err != nil {
		return err
	}
	t.cursor = 0
	return nil
}

func (t *Encoder) Write(b []byte) (n int, err error) {
	available := len(t.b)
	for _, c := range b {
		if available-t.cursor < 2 {
			// always have at least two bytes available
			// in case need to escape
			if err = t.Flush(); err != nil {
				return 0, err
			}
		}

		switch c {
		case t.mark, t.escape:
			t.b[t.cursor] = t.escape
			t.cursor++
			fallthrough
		default:
			t.b[t.cursor] = c
			t.cursor++
		}
	}
	return len(b), nil
}

// Cut flushes the buffer and writes repeated telomereEscapeBytes to the underlying Writer.
func (t *Encoder) Cut() (n int, err error) {
	if err = t.Flush(); err != nil {
		return 0, err
	}
	j, err := io.Copy(t.w, bytes.NewReader(t.t))
	return int(j), err
}

func (t *Encoder) Close() (err error) {
	_, err = t.Cut()
	return err
}
