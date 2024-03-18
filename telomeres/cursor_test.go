package telomeres

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type chunkFollowedByTelomere struct {
	chunk     string
	telomeres int
}

var cursorTestCases = [...]*chunkFollowedByTelomere{
	&chunkFollowedByTelomere{"a", 4},
	&chunkFollowedByTelomere{"bb", 4},
	&chunkFollowedByTelomere{"ccc", 4},
	&chunkFollowedByTelomere{"dddd", 5},
	&chunkFollowedByTelomere{"eeeee", 6},
}

func TestCursorPositionReporting(t *testing.T) {

	telomeres, err := New(WithMinimumCount(4))
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
	t.Log("edge:", b.String())

	decoder := telomeres.NewDecoder(bytes.NewReader(b.Bytes()))
	n, err := decoder.FindEdge(9999999999)
	if n != 4 {
		t.Fatalf("did not find expected data edge %d vs %d", n, 4)
	}
	cursor := decoder.Cursor()
	if cursor != 4 {
		t.Fatalf("did not find expected data edge cursor %d vs %d", cursor, 4)
	}

	chunk := &bytes.Buffer{}
	data := b.String()
	for _, tc := range cursorTestCases {
		if n, err = decoder.WriteTo(chunk); err != nil {
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
		if cursor != decoder.Cursor()+n {
			t.Fatalf("did not find expected chunk edge cursor %d vs %d", cursor, decoder.Cursor()+n)
		}
		if data[cursor-n:cursor] != tc.chunk {
			t.Log("     got:", chunk.String())
			t.Log("expected:", tc.chunk)
			t.Fatal("data mismatch when checking cursor bounds")
		}
	}

	// var n int64
	// for _, tc := range encodingTestCases {
	// 	d := telomeres.NewDecoder(strings.NewReader(tc.out))
	// }
}
