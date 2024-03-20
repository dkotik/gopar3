package telomeres

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"
)

func TestStoredCases(t *testing.T) {
	data, err := os.ReadFile("testdata/primary.txt")
	if err != nil {
		t.Fatal(err)
	}
	decoder := NewDecoder(newTestBuffer(data))

	b := &bytes.Buffer{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	for {
		_, err = decoder.StreamChunk(ctx, b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Error(err)
		}
		b.Reset()
	}
}
