package telomeres

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestStoredCases(t *testing.T) {
	f, err := os.Open("testdata/primary.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	telomeres, err := New(
		WithMinimumCount(4),
		WithBufferSize(65),
	)
	if err != nil {
		t.Error(err)
	}
	decoder := telomeres.NewDecoder(f)

	b := &bytes.Buffer{}
	for {
		_, err = decoder.WriteTo(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Error(err)
		}
	}
}
