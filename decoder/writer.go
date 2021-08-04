package decoder

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func (d *Decoder) WriteAll(w io.Writer, in <-chan ([][]byte)) error {
	// f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 600)
	// if err != nil {
	// 	return err
	// }
	errc := make(chan (error))

	go func() {
		var (
			err error
			n   int64
			i   uint64
		)
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not write piece №%d: %w", i, err)
			}
			errc <- err
		}()

		for batch := range in {
			for _, piece := range batch {
				n, err = io.Copy(w, bytes.NewReader(piece))
				if err != nil {
					return
				}
				if n != int64(len(piece)) { // TODO: is this needed?
					err = errors.New("the entire piece did not fit")
					return
				}
				i++
			}
		}
	}()

	return <-errc
}
