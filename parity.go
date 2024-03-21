package gopar3

import (
	"github.com/klauspost/reedsolomon"
)

// AddParity fills out the missing chards with parity data
// that can be used to restore the original data shards.
// Closes the out channel when finished.
func AddParity(
	out chan<- [][]byte,
	in <-chan [][]byte,
	rs reedsolomon.Encoder,
) (err error) {
	var batch [][]byte
	for {
		select {
		case batch = <-in:
			if batch == nil {
				close(out)
				return nil // finished
			}
			if err = rs.Reconstruct(batch); err != nil {
				close(out)
				return err
			}
			out <- batch
		}
	}
}
