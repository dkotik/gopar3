package telomeres

import (
	"bytes"
	"io"
)

type TelomereStreamEncoder struct {
	t []byte
	b []byte
	w io.Writer
}

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

func (t *TelomereStreamEncoder) Write(b []byte) (n int64, err error) {

}

func (t *TelomereStreamEncoder) WriteTelomere() (n int, err error) {
	j, err := io.Copy(t.w, bytes.NewReader(t.t))
	return int(j), err
}
