/*
Package main provides a command line interface to:

- [gopar3.Inflate]
- [gopar3.Split]
- [gopar3.Scatter]
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dkotik/gopar3"
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
				Action: func(ctx *cli.Context) (err error) {
					sources := ctx.Args().Slice()
					if len(sources) == 0 {
						return errors.New("provide at least one source")
					}
					for _, source := range sources {
						fmt.Println("inflating: ", source)
						if err = gopar3.Inflate(
							ctx.Context,
							".",
							source,
							5,
							3,
							64,
						); err != nil {
							return err
						}
					}
					return nil
				},
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
