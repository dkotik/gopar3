package telomeres

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TODO: // decoder.Skip() does not seem to work by itself

func TestDecoding(t *testing.T) {
	telomeres, err := New(WithMinimumCount(4))
	if err != nil {
		t.Error(err)
	}

	b := &bytes.Buffer{}
	for _, tc := range encodingTestCases {
		b.Reset()
		e := telomeres.NewEncoder(b)
		e.Cut()
		for _, chunk := range tc.in {
			if _, err = io.Copy(e, strings.NewReader(chunk)); err != nil {
				t.Error(err)
			}
			if _, err = e.Cut(); err != nil {
				t.Error(err)
			}
		}

		if b.String() != tc.out {
			t.Log("expecting:", tc.out)
			t.Log("    given:", b.String())
			t.Error("encoder output does not match expectation")
		}
	}
}
