package main

import (
	"flag"
	"fmt"
)

var fragments uint8 = 8

var (
	// fragments   = flag.Uint("fragments", 13, "break input into this many parts")
	growth      = flag.Float64("growth", 1.3, "all fragments together will take up this much  more space than the input")
	telomeres   = flag.Uint("telomeres", 8, "length of telomere padding protecting shard boundaries  more padding increases output resilience")
	memoryLimit = flag.Int("memoryMB", 128, "memory usage limit for operations  in Megabytes")

	requiredShards, redundantShards uint8
	command, target                 string
)

func main() {
	must(parseFlags())
	switch command {
	case "encode":
		// flagset.Parse(os.Args[2:])
		// flagset.PrintDefaults()
		fmt.Printf("%d %t", fragments, *growth)
		// fmt.Printf("%d %t", uint16(*fragments), uint8(*fragments))
		return
	case "decode":
		// *memoryLimit * 1024 * 1024
		return
	}
	flag.Usage()
}
