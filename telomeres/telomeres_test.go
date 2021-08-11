package telomeres

import (
	"bytes"
	"fmt"
)

func ExampleEncoder_Encode() {
	b := &bytes.Buffer{}
	e := NewEncoder(b, 4, 1024)

	e.Write([]byte("hello"))
	e.Cut()
	e.Write([]byte("world"))
	e.Cut()
	// e.Flush() // why is this NOT needed?

	fmt.Print(b.String())
	// Output: hello::::world::::
}
