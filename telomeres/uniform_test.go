package telomeres

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUniformRandomEncodingDecoding(t *testing.T) {
	t.Skip("impl")
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
	filePath := filepath.Join(t.TempDir(), "uniformTest.txt")
	// filePath := filepath.Join(os.TempDir(), "uniformTest.txt")
	b, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// t.Log("temp file:", filePath)
		if err = b.Close(); err != nil {
			t.Fatal(err)
		}
	})
	encoder := telomeres.NewEncoder(b)

	for _, tc := range testCases {
		chunk := randomData(tc)
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
		if _, err = b.Seek(0, io.SeekStart); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		d := &bytes.Buffer{}
		decoder := telomeres.NewDecoder(b)
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
		if nn != int64(len(chunk)) || !bytes.Equal(d.Bytes(), chunk) {
			t.Logf("expected: %q", string(chunk))
			t.Logf(" decoded: %q", d.String())
			t.Fatalf("wrote %d decoded chunk bytes, but expecting %d", nn, len(chunk))
		}
	}
}
