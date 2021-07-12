package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	flagset := flag.NewFlagSet(command, flag.ContinueOnError)

	switch command {
	case "encode":
		fragments := flagset.Uint("fragments", 13, "break input into this many parts")
		growth := flagset.Float64("growth", 1.3,
			"all fragments together will take up\nthis much  more space than the input")
		flagset.Parse(os.Args[2:])
		flagset.PrintDefaults()
		fmt.Printf("%d %t", *fragments, *growth)
		// fmt.Printf("%d %t", uint16(*fragments), uint8(*fragments))
	case "decode":
		// memoryLimit := flagset.Uint("memoryMB", 200, "the amount of memory to use")
	default:
		fmt.Printf(`gopar3 %s data resilience utility

`, "gopar3.Version")

		fmt.Printf("%v", os.Args)
		flag.PrintDefaults()
	}
}
