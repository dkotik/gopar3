package gopar3

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dkotik/gopar3/telomeres"
)

func ExampleScanner_Pipe() {
	shards := bytes.NewBuffer([]byte(`:::`))
	for _, word := range []string{`one`, `two`, `three`} {
		shards.WriteString(word)
		shards.WriteString(strings.Repeat("-", 20)) // tag padding
		chk := NewChecksum()
		chk.Write([]byte(word))
		chk.Write([]byte(strings.Repeat("-", 20)))
		shards.Write(chk.Sum(nil))
		shards.WriteString(`:::`)
	}
	// panic(shards)

	scanner := &Scanner{
		MaxBytesBeforeGivingUp: 2 << 8,
		ChecksumFactory:        NewChecksum,
		ErrorHandler: func(err error) bool {
			// panic(err)
			fmt.Println("Error:", err)
			return true
		},
	}
	out := make(chan ([]byte))

	decoder := telomeres.NewDecoder(shards, 3, 8, 2<<8)
	scanner.Pipe(context.Background(), decoder, out)
	time.Sleep(time.Second)
	for i := 0; i < 4; i++ {
		fmt.Printf("%s\n", <-out)
	}

	// Output:
	// one--------------------
	// two--------------------
	// three--------------------
}
