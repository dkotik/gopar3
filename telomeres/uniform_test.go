package telomeres

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
)

func TestUniformRandomEncodingDecoding(t *testing.T) {
	testCases := [...]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73,
		125, 126, 127, 128, 129, 130, 131, 136,
	}
	b := &bytes.Buffer{}
	encoder, err := NewEncoder(b, 5)
	if err != nil {
		t.Fatal(err)
	}

	var dataByte byte = '!'
	// var dataByte byte = Mark
	// var dataByte byte = Escape
	for _, tc := range testCases {
		chunk := bytes.Repeat([]byte{dataByte}, tc)
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
	}

	decoder := NewDecoder(newTestBuffer(b.Bytes()))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	d := &bytes.Buffer{}
	for _, tc := range testCases {
		nn, err := decoder.StreamChunk(ctx, d)
		if err != nil && err != io.EOF {
			// if !(errors.Is(err, io.EOF) && len(chunk) == 0) {
			// 	// t.Log("  buffer:", b.String())
			// 	t.Log("error:", err)
			// 	t.Logf("expected: %q", string(chunk))
			// 	t.Log("length:", len(chunk))
			// }
			t.Fatal(err)
		}
		if nn != int64(tc) {
			t.Logf(" decoded: %q", d.String())
			t.Fatalf("wrote %d decoded chunk bytes, but expecting %d", nn, tc)
		}
		for _, c := range d.Bytes() {
			if c != dataByte {
				t.Logf(" decoded: %q", d.String())
				t.Fatalf("expected all %q bytes, but got one %q", dataByte, c)
			}
		}
	}
}
