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
	{
		in:  []string{"000000000000000000000000000000000000000000000000000000000000000000000000"},
		out: "::::000000000000000000000000000000000000000000000000000000000000000000000000::::",
	},
}

func TestEncoding(t *testing.T) {
	b := &bytes.Buffer{}
	for _, tc := range encodingTestCases {
		b.Reset()
		e, err := NewEncoder(b, 4)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = e.Cut(); err != nil {
			t.Fatal(err)
		}
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

func TestEmptyEncoding(t *testing.T) {
	b := &bytes.Buffer{}
	encoder, err := NewEncoder(b, 4)
	if err != nil {
		t.Fatal(err)
	}
	n, err := encoder.Write(nil)
	if err != nil {
		t.Fatal(err)
	}
	if n > 0 {
		t.Fatalf("encoded %d bytes, but should have been zero", n)
	}
}
