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
		Usage: "protect data from partial loss or corruption",
		Commands: []*cli.Command{
			{
				Name:    "inflate",
				Aliases: []string{"i"},
				Usage:   "one output file for each input file",
				Flags: []cli.Flag{
					flagOutput,
					flagQuorum,
					flagParity,
					flagSize,
				},
				Action: commandInflate,
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
