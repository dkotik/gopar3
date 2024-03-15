package telomeres

import (
	"bytes"
	"errors"
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
	var n int64
	for _, tc := range encodingTestCases {
		d := telomeres.NewDecoder(strings.NewReader(tc.out))

		for _, chunk := range tc.in {
			n, err = d.WriteTo(b)
			if err != nil {
				t.Error(err)
			}
			if chunk != b.String() {
				t.Log("expecting:", chunk)
				t.Log("    given:", b.String())
				t.Error("decoder output does not match expectation")
			}
			if n != int64(b.Len()) {
				t.Errorf("decoded %d bytes, but expected %d instead", b.Len(), n)
			}
			b.Reset()
		}

		n, err = d.WriteTo(b)
		if !errors.Is(err, io.EOF) {
			t.Errorf("expecing io.EOF got this instead: %+v", err)
		}
		if n != 0 {
			t.Errorf("got %d written bytes, but should have been 0", n)
		}
	}
}

func TestEmptyDecoding(t *testing.T) {
	telomeres, err := New()
	if err != nil {
		t.Error(err)
	}

	decoder := telomeres.NewDecoder(bytes.NewReader(nil))
	b := &bytes.Buffer{}
	n, err := decoder.WriteTo(b)
	if !errors.Is(err, io.EOF) {
		t.Error("expected io.EOF but instead got:", err)
	}
	if n > 0 {
		t.Errorf("decoded %d bytes, but should have been zero", n)
	}
}
