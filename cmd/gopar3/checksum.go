package main

import (
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"sync"

	"github.com/dkotik/gopar3"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func commandChecksum(cliCtx *cli.Context) (err error) {
	sources := cliCtx.Args().Slice()
	if len(sources) == 0 {
		return cli.ShowSubcommandHelp(cliCtx)
	}
	wg, ctx := errgroup.WithContext(cliCtx.Context)
	wg.SetLimit(runtime.NumCPU())

	type sumResult struct {
		Source        string
		CastagnoliSum uint32
	}

	results := make([]sumResult, 0, len(sources))
	mu := &sync.Mutex{}
	defer func() {
		if len(results) == 0 {
			return
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err = errors.Join(err, encoder.Encode(results))
	}()

	for _, source := range sources {
		wg.Go(func() (err error) {
			f, err := os.Open(source)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, f.Close())
			}()
			sum, err := gopar3.CastagnoliSum(ctx, f)
			if err != nil {
				return err
			}
			mu.Lock()
			results = append(results, sumResult{
				Source:        source,
				CastagnoliSum: sum,
			})
			mu.Unlock()
			return
		})
	}
	return wg.Wait()
}
