package telomeres

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
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
		"\\\\:::::",
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
			t.Fatal(err)
		}
		b := &testBuffer{}
		encoder := telomeres.NewEncoder(b)
		n, err := encoder.Cut()
		if err != nil {
			t.Fatal(err)
		}
		if n != tl {
			t.Fatalf("wrote %d cut bytes, but expecting %d", n, tl)
		}
		n, err = encoder.Write([]byte(chunk))
		if err != nil {
			t.Fatal(err)
		}
		if n != len(chunk) {
			t.Fatalf("wrote %d chunk bytes, but expecting %d", n, len(chunk))
		}
		n, err = encoder.Cut()
		if err != nil {
			t.Fatal(err)
		}
		if n != tl {
			t.Fatalf("wrote %d cut bytes, but expecting %d", n, tl)
		}
		t.Log("  buffer:", b.String())
		cursor, err := b.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal("error seeking buffer:", err)
		}
		if cursor != 0 {
			t.Fatal("buffer cursor should be at 0, but it is at:", cursor)
		}

		d := &bytes.Buffer{}
		decoder := telomeres.NewDecoder(b)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		nn, err := decoder.StreamChunk(ctx, d)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		if nn != int64(len(chunk)) || d.String() != chunk {
			t.Log("   error:", err)
			t.Log("expected:", chunk)
			t.Log(" decoded:", d.String())
			t.Fatalf("wrote %d decoded chunk bytes, but expecting %d", nn, len(chunk))
		}
	})
}
