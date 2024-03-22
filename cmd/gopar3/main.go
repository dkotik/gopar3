/*
Package main provides a command line interface to:

- [gopar3.Inflate]
- [gopar3.Split]
- [gopar3.Scatter]
*/
package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gopar3",
		Usage: "(Alpha) protect data from partial loss or corruption",
		Commands: []*cli.Command{
			{
				Name:      "inflate",
				Aliases:   []string{"i"},
				Usage:     "one output file for each input file",
				ArgsUsage: "[...FILES]",
				Flags: []cli.Flag{
					flagOutput,
					flagQuorum,
					flagParity,
					flagSize,
				},
				Action: commandInflate,
			},
			{
				Name:      "inspect",
				Aliases:   []string{"s"},
				Usage:     "scan each input file or directory for data shards",
				ArgsUsage: "[...FILES]",
				Flags:     []cli.Flag{
					// flagOutput,
					// flagQuorum,
					// flagParity,
					// flagSize,
				},
				Action: commandInspect,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// TODO: add cancellation

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
