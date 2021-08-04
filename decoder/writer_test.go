package decoder

import (
	"fmt"
	"os"
)

func ExampleDecoder_WriteAll() {
	d := &Decoder{}
	in := make(chan ([][]byte))
	go func() {
		defer close(in)
		in <- [][]byte{
			[]byte("| hey"),
			[]byte(" ... "),
			[]byte("now | "),
		}
	}()

	err := d.WriteAll(os.Stdout, in)
	fmt.Println("and the error is", err)
	// Output: | hey ... now | and the error is <nil>
}
