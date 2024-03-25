package main

import (
	"errors"
	"os"

	"github.com/dkotik/gopar3"
	"github.com/urfave/cli/v2"
)

func commandRestore(cliCtx *cli.Context) (err error) {
	sources := cliCtx.Args().Slice()
	if len(sources) == 0 {
		return cli.ShowSubcommandHelp(cliCtx)
	}

	index, err := gopar3.NewIndex(cliCtx.Context, sources...)
	if err != nil {
		return err
	}
	if len(index) == 0 {
		return errors.New("no files to restore")
	}

	var w *os.File
	for differentiator, file := range index {
		w, err = os.Create(differentiator + ".tmp") // TODO: check if exists
		if err != nil {
			return err
		}
		err = errors.Join(gopar3.Restore(cliCtx.Context, w, file), w.Close())
		break
	}

	return err
}
