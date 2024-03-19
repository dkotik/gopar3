package telomeres

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func TestDecoding(t *testing.T) {
	telomeres, err := New(WithMinimumCount(4))
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	b := &bytes.Buffer{}
	var n int64
	for _, tc := range encodingTestCases {
		d := telomeres.NewDecoder(newTestBuffer([]byte(tc.out)))

		for _, chunk := range tc.in {
			n, err = d.StreamChunk(ctx, b)
			if err != nil {
				t.Error(err)
			}
			if chunk != b.String() {
				t.Log("expecting:", chunk)
				t.Log("    given:", b.String())
				t.Error("decoder output does not match expectation")
			}
			if n != int64(b.Len()) {
				t.Errorf("decoded %d bytes, but expected %d instead", b.Len(), n)
			}
			b.Reset()
		}

		n, err = d.StreamChunk(ctx, b)
		if !errors.Is(err, io.EOF) {
			t.Errorf("expecing io.EOF got this instead: %+v", err)
		}
		if n != 0 {
			t.Errorf("got %d written bytes, but should have been 0", n)
		}
	}
}

func TestEmptyDecoding(t *testing.T) {
	telomeres, err := New()
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	decoder := telomeres.NewDecoder(&testBuffer{})
	b := &bytes.Buffer{}
	n, err := decoder.StreamChunk(ctx, b)
	if !errors.Is(err, io.EOF) {
		t.Error("expected io.EOF but instead got:", err)
	}
	if n > 0 {
		t.Errorf("decoded %d bytes, but should have been zero", n)
	}
}
