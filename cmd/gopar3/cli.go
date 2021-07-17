package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func printFlag(f *flag.Flag) {
	if len(f.Name) < 3 {
		return
	}
	defaultValue := ""
	if f.DefValue != "" && f.DefValue != "false" {
		defaultValue = fmt.Sprintf(" (default: %s)", f.DefValue)
	}

	fmt.Printf("  --%s%s%s%s\n\n",
		f.Name,
		strings.Repeat(" ", 12-len(f.Name)),
		strings.Replace(f.Usage, "  ", "\n                ", 20),
		defaultValue)
}

func parseFlags() (err error) {
	flag.Usage = func() {
		fmt.Printf("gopar3 v%s data resilience utility\n\n  gopar3 [encode|decode] [FILE] [FLAGS]\n\n", "gopar3.Version")
		flag.VisitAll(printFlag)
		printFlag(&flag.Flag{
			Name:  "help, -h",
			Usage: "display help message",
		})
	}

	flag.Var(&shardValue{&fragments}, "fragments", "fragments usage")
	// flag.UintVar(fragments, "f", 13, "shorthand for --fragments")
	flag.Float64Var(growth, "g", 1.3, "shorthand for --growth")
	flag.UintVar(telomeres, "t", 8, "shorthand for --telomeres")
	flag.IntVar(memoryLimit, "m", 200, "shorthand for --memoryMB")
	flag.Parse()
	command = flag.Arg(0)
	if command == "" {
		command = "help"
	} else {
		flag.CommandLine.Parse(flag.Args()[1:])
		if target = flag.Arg(0); target != "" {
			flag.CommandLine.Parse(flag.Args()[1:])
		}
	}

	if *growth < 1 {
		err = errors.New("growth cannot be less than 1")
		return
	}
	if *memoryLimit < 4 {
		err = errors.New("cannot use less than 4MB of memory")
	}

	// requiredShards = uint8(float64(*fragments) / *growth)
	// if requiredShards == 0 {
	// 	requiredShards = 1
	// }
	// redundantShards = uint8(*fragments) - requiredShards

	return
}

type shardValue struct {
	v *uint8
}

func (s shardValue) String() string {
	return fmt.Sprintf("%s", s.v)
}

func (s shardValue) Set(v string) error {
	if !regexp.MustCompile(`^\d+$`).MatchString(v) {
		return fmt.Errorf("invalid shard number value: %q", v)
	}
	fmt.Sscanf(v, "%d", s.v)
	if *s.v == 0 {
		return fmt.Errorf("shard number cannot be 0")
	}
	return nil
}

func must(err error) {
	if err == flag.ErrHelp {
		flag.Usage()
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("Error: %s.\n", err.Error())
		os.Exit(1)
	}
}
