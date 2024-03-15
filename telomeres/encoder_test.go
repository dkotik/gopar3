package telomeres

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

var encodingTestCases = [...]struct {
	in  []string
	out string
}{
	{
		in:  []string{"1", "2", "3", ":"},
		out: "::::1::::2::::3::::\\:::::",
	},
	{
		in:  []string{"111", "2222", "33333", "44:abc"},
		out: "::::111::::2222::::33333::::44\\:abc::::",
	},
	{
		in:  []string{"111", "2222", "33333", "44\\abc"},
		out: "::::111::::2222::::33333::::44\\\\abc::::",
	},
}

func TestEncoding(t *testing.T) {
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
