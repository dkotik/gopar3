package main

import (
	"github.com/dkotik/gopar3"
	"github.com/urfave/cli/v2"
)

func commandInflate(ctx *cli.Context) (err error) {
	sources := ctx.Args().Slice()
	if len(sources) == 0 {
		return cli.ShowSubcommandHelp(ctx)
	}
	for _, source := range sources {
		// fmt.Println("inflating: ", source)
		if err = gopar3.Inflate(
			ctx.Context,
			ctx.String("output"),
			source,
			uint8(ctx.Uint("quorum")),
			uint8(ctx.Uint("parity")),
			ctx.Int("size"),
		); err != nil {
			return err
		}
	}
	return nil
}
