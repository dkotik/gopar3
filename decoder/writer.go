package decoder

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func (d *Decoder) WriteAll(w io.Writer, in <-chan ([][]byte)) error {
	errc := make(chan (error))

	go func() {
		var (
			err error
			n   int64
			i   uint64
		)
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not write shard №%d: %w", i, err)
			}
			errc <- err
		}()

		for batch := range in {
			for i, shard := range batch {
				if shard == nil || i > d.batchSize { // TODO: reduce limit to req size later
					continue
				}
				n, err = io.Copy(w, bytes.NewReader(shard))
				if err != nil {
					return
				}
				if n != int64(len(shard)) { // TODO: is this needed?
					err = errors.New("the entire shard did not fit")
					return
				}
				i++
			}
		}
	}()

	return <-errc
}
