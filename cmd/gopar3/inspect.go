package main

import (
	"encoding/json"
	"os"

	"github.com/dkotik/gopar3"
	"github.com/urfave/cli/v2"
)

func commandInspect(ctx *cli.Context) (err error) {
	sources := ctx.Args().Slice()
	if len(sources) == 0 {
		return cli.ShowSubcommandHelp(ctx)
	}
	index, err := gopar3.NewIndex(ctx.Context, sources...)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(index)
}
