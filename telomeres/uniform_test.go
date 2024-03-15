package telomeres

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestUniformRandomEncodingDecoding(t *testing.T) {
	testCases := [...]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73,
		125, 126, 127, 128, 129, 130, 131, 136,
	}

	telomeres, err := New(
		WithMinimumCount(5),
		WithBufferSize(67),
	)
	if err != nil {
		t.Error(err)
	}
	b := &bytes.Buffer{}
	for _, tc := range testCases {
		chunk := randomData(tc)
		encoder := telomeres.NewEncoder(b)
		n, err := encoder.Cut()
		if err != nil {
			t.Error(err)
		}
		if n != 5 {
			t.Errorf("wrote %d cut bytes, but expecting %d", n, 5)
		}
		n, err = encoder.Write(chunk)
		if err != nil {
			t.Error(err)
		}
		if n != len(chunk) {
			t.Errorf("wrote %d chunk bytes, but expecting %d", n, len(chunk))
		}
		n, err = encoder.Cut()
		if err != nil {
			t.Error(err)
		}
		if n != 5 {
			t.Errorf("wrote %d cut bytes, but expecting %d", n, 5)
		}

		d := &bytes.Buffer{}
		decoder := telomeres.NewDecoder(b)
		nn, err := decoder.WriteTo(d)
		if err != nil {
			if !(errors.Is(err, io.EOF) && len(chunk) == 0) {
				t.Log("expected:", string(chunk))
				t.Log("length:", len(chunk))
				t.Fatal(err)
			}
		}
		if nn != int64(len(chunk)) || !bytes.Equal(d.Bytes(), chunk) {
			t.Logf("expected: %q", string(chunk))
			t.Logf(" decoded: %q", d.String())
			t.Fatalf("wrote %d decoded chunk bytes, but expecting %d", nn, len(chunk))
		}
		b.Reset()
	}
}
