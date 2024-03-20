package telomeres

import (
	"bytes"
	"errors"
	"io"
)

// Encoder appends an [Escape] byte to each [Mark] byte.
type Encoder struct {
	w io.Writer
	b *bytes.Buffer
	t []byte
}

// NewEncoder creates a telomere encoder.
func NewEncoder(w io.Writer, telomereCount int) (*Encoder, error) {
	if telomereCount < 1 {
		return nil, errors.New("telomere count must be greater than one")
	}

	telomeres := make([]byte, telomereCount)
	for i := range telomeres {
		telomeres[i] = Mark
	}

	return &Encoder{
		w: w,
		b: &bytes.Buffer{},
		t: telomeres,
	}, nil
}

func (t *Encoder) Write(b []byte) (n int, err error) {
	var (
		c      byte
		window = b
	)
	t.b.Reset()

encode:
	for n, c = range window {
		switch c {
		case Mark, Escape:
			_, _ = t.b.Write(window[:n])
			_ = t.b.WriteByte(Escape)
			_ = t.b.WriteByte(c)
			window = window[n+1:]
			goto encode
		}
	}
	_, _ = t.b.Write(window)
	if n, err = t.w.Write(t.b.Bytes()); err != nil {
		escaped := false
		for _, c = range t.b.Bytes()[n:] {
			if escaped {
				escaped = false
				continue
			}
			if c == Escape {
				escaped = true
				n--
			}
		}
		return n, err
	}
	return len(b), nil
}

// Cut writes [Mark]s to the underlying Writer to indicate the end
// of a data chunk.
func (t *Encoder) Cut() (n int, err error) {
	return t.w.Write(t.t)
}
