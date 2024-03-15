package telomeres

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func FuzzSingleChunk(f *testing.F) {
	for _, tc := range [...]string{
		"Hello, world",
		" ",
		"!12345",
		":",
		"::",
		"::::::::::",
		"\\",
		"\\\\",
	} {
		f.Add(tc) // seed corpus
	}
	f.Fuzz(func(t *testing.T, chunk string) {
		const tl = 4

		telomeres, err := New(
			WithMinimumCount(tl),
			WithBufferSize(79),
		)
		if err != nil {
			t.Error(err)
		}
		b := &bytes.Buffer{}
		encoder := telomeres.NewEncoder(b)
		n, err := encoder.Cut()
		if err != nil {
			t.Error(err)
		}
		if n != tl {
			t.Errorf("wrote %d cut bytes, but expecting %d", n, tl)
		}
		n, err = encoder.Write([]byte(chunk))
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
		if n != tl {
			t.Errorf("wrote %d cut bytes, but expecting %d", n, tl)
		}
		t.Log("  buffer:", b.String())

		d := &bytes.Buffer{}
		decoder := telomeres.NewDecoder(b)
		nn, err := decoder.WriteTo(d)
		if err != nil {
			if !(errors.Is(err, io.EOF) && len(chunk) == 0) {
				t.Log("expected:", chunk)
				t.Log("length:", len(chunk))
				t.Error(err)
			}
		}
		if nn != int64(len(chunk)) || d.String() != chunk {
			t.Log("expected:", chunk)
			t.Log(" decoded:", d.String())
			t.Errorf("wrote %d decoded chunk bytes, but expecting %d", nn, len(chunk))
		}
	})
}
