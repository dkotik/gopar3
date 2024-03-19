package telomeres

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

type chunkFollowedByTelomere struct {
	chunk     string
	telomeres int
}

var cursorTestCases = [...]*chunkFollowedByTelomere{
	&chunkFollowedByTelomere{"a", 4},
	&chunkFollowedByTelomere{"bb", 400},
	&chunkFollowedByTelomere{"ccc", 400},
	&chunkFollowedByTelomere{"dddd", 500},
	&chunkFollowedByTelomere{"eeeee", 6},
}

func TestCursorPositionReporting(t *testing.T) {
	telomeres, err := New(
		WithMinimumCount(4),
		WithBufferSize(71),
	)
	if err != nil {
		t.Error(err)
	}

	b := &bytes.Buffer{}
	encoder := telomeres.NewEncoder(b)
	_, _ = io.WriteString(b, `::::`)
	for _, tc := range cursorTestCases {
		_, err = io.WriteString(encoder, tc.chunk)
		// _, err = encoder.Write([]byte(tc.chunk))
		if err != nil {
			t.Fatal(err)
		}
		if err = encoder.Flush(); err != nil {
			t.Fatal(err)
		}
		_, _ = io.WriteString(b, strings.Repeat(`:`, tc.telomeres))
	}
	// t.Log("edge:", b.String())

	decoder := telomeres.NewDecoder(newTestBuffer(b.Bytes()))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err = decoder.SeekChunk(ctx); err != nil {
		t.Fatal("failed to seek:", err)
	}
	cursor, err := decoder.Cursor()
	if err != nil {
		t.Fatal("failed to get cursor position:", err)
	}
	if cursor != 4 {
		t.Fatalf("did not find expected data edge %d vs %d", cursor, 4)
	}
	// cursor += 4

	chunk := &bytes.Buffer{}
	// data := b.String()
	var n int64
	for _, tc := range cursorTestCases {
		// start, err := decoder.Cursor()
		if err != nil {
			t.Fatal("failed to get cursor position:", err)
		}

		if n, err = decoder.StreamChunk(ctx, chunk); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		if chunk.String() != tc.chunk {
			t.Log("     got:", chunk.String())
			t.Log("expected:", tc.chunk)
			t.Fatal("case mismatch")
		}
		cursor += n
		decoderCursor, err := decoder.Cursor()
		if err != nil {
			t.Fatal("failed to get cursor position:", err)
		}

		if cursor != decoderCursor {
			t.Logf("added: %d", n)
			t.Fatalf("did not find expected chunk edge cursor %d vs %d", cursor, decoderCursor)
		}
		// if data[start:cursor] != tc.chunk {
		// 	t.Logf("     got: %q", data[start:cursor])
		// 	t.Logf("expected: %q", tc.chunk)
		// 	t.Fatal("data mismatch when checking cursor bounds")
		// }
		chunk.Reset()
		cursor += int64(tc.telomeres)
	}

	// var n int64
	// for _, tc := range encodingTestCases {
	// 	d := telomeres.NewDecoder(strings.NewReader(tc.out))
	// }
}
