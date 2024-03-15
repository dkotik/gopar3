package telomeres

import (
	"bytes"
	"fmt"
)

func ExampleEncoder_Encode() {
	b := &bytes.Buffer{}
	t, err := New(WithMinimumCount(4))
	if err != nil {
		panic(err)
	}
	e := t.NewEncoder(b)

	_, _ = e.Cut()
	_, _ = e.Write([]byte("hello"))
	_, _ = e.Cut()
	_, _ = e.Write([]byte("world"))
	_, _ = e.Cut()

	fmt.Print(b.String())
	// Output: ::::hello::::world::::
}

func ExampleDecoder_Decode() {
	b := bytes.NewBuffer([]byte("::::hello::::world::::"))
	t, err := New(WithMinimumCount(4))
	if err != nil {
		panic(err)
	}
	d := t.NewDecoder(b)

	s1 := &bytes.Buffer{}
	_, err = d.WriteTo(s1)
	if err != nil {
		panic(err)
	}

	s2 := &bytes.Buffer{}
	_, err = d.WriteTo(s2)
	if err != nil {
		panic(err)
	}

	fmt.Println(s1.String(), "||", s2.String())
	// Output: hello || world
}
