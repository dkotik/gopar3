package main

import (
	"fmt"
	"os"

	"github.com/dkotik/gopar3"
	flag "github.com/spf13/pflag"
)

var (
	fragments   = flag.UintP("fragments", "f", 9, "break input into this many parts")
	growth      = flag.Float64P("growth", "g", 1.3, "all fragments together will take up this much\nmore space than the input")
	telomeres   = flag.UintP("telomeres", "t", 8, "length of telomere padding protecting\nshard boundaries more padding increases\noutput resilience")
	memoryLimit = flag.IntP("memoryMB", "m", 128, "memory usage limit for operations\nin Megabytes")
	help        = flag.BoolP("help", "h", false, "print help message")
)

func main() {
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.CommandLine.MarkHidden("help")
	var err error
	flag.Usage = func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
			return
		}
		fmt.Fprintf(os.Stderr, "gopar3 v%s data resilience utility\n\n  gopar3 [encode|decode] [FILE] [FLAGS]\n\n", gopar3.Version)
		flag.PrintDefaults()
	}

	if !*help {
		switch flag.Arg(0) {
		case "encode":
			// spew.Dump(fragments, growth)
			// flagset.Parse(os.Args[2:])
			// flagset.PrintDefaults()
			// fmt.Printf("%d %t", fragments, *growth)
			// fmt.Printf("%d %t", uint16(*fragments), uint8(*fragments))
			return
		case "decode":
			// *memoryLimit * 1024 * 1024
			return
		case "version":
			fmt.Fprintf(os.Stderr, "gopar3 v%s\n", gopar3.Version)
			return
		}
	}
	flag.Usage()
}
