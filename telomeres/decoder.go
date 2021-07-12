package telomeres

import "io"

type TelomereStreamDecoder struct {
	t []byte
	b []byte
	r io.Reader
}

func NewTelomereStreamDecoder(r io.Reader, telomereLength, bufferSize int) *TelomereStreamEncoder {
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

func (t *TelomereStreamDecoder) Read(b []byte) (n int, err error) {

}
