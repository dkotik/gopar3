package telomeres

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// TODO: // decoder.Skip() does not seem to work by itself

func TestDecoding(t *testing.T) {
	b := bytes.NewBuffer([]byte("hello::::world::::"))
	d := NewDecoder(b, 4, 100, 1024)

	s1 := &bytes.Buffer{}
	_, err := io.Copy(s1, d)
	if err != nil {
		t.Fatal("first", err)
	}
	s2 := &bytes.Buffer{}
	_, err = io.Copy(s2, d)
	if err != nil {
		t.Fatal("second", err)
	}

	if h := s1.String(); h != "hello" {
		t.Fatal(fmt.Errorf("first part %q does not equal 'hello'", h))
	}
	if h := s2.String(); h != "world" {
		t.Fatal(fmt.Errorf("second part %q does not equal 'world'", h))
	}
}

func ExampleDecoder_Decode() {
	b := bytes.NewBuffer([]byte("hello::::world::::"))
	d := NewDecoder(b, 4, 100, 1024)

	s1 := &bytes.Buffer{}
	io.Copy(s1, d)
	s2 := &bytes.Buffer{}
	io.Copy(s2, d)

	fmt.Println(s1.String(), s2.String())
	// Output: hello world
}
