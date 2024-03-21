package main

import (
	"github.com/urfave/cli/v2"
)

var (
	// TODO: SDTIN SDTOUT flags
	// growth      = flag.Float64P("growth", "g", 1.3, "all fragments together will take up this much\nmore space than the input")
	// telomeres   = flag.UintP("telomeres", "t", 8, "length of telomere padding protecting\nshard boundaries more padding increases\noutput resilience")
	// version = "Alpha"

	flagOutput = &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Value:   ".",
		Usage:   "`destination` for created files",
	}

	flagQuorum = &cli.UintFlag{
		Name:    "quorum",
		Aliases: []string{"q"},
		Value:   5,
		Usage:   "`number` of intact shards required for restoration",
	}

	flagParity = &cli.UintFlag{
		Name:    "parity",
		Aliases: []string{"p"},
		Value:   3,
		Usage:   "`number` of parity shards",
	}

	flagSize = &cli.UintFlag{
		Name:    "size",
		Aliases: []string{"s"},
		Value:   64, // TODO: fix
		Usage:   "size of each shard in `bytes` without the metadata",
	}
)
