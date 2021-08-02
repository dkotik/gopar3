package encoder

import (
	"github.com/klauspost/reedsolomon"
)

// CompleteWithReedSolomon generates missing redundant shards.
func (e *Encoder) CompleteWithReedSolomon(bb <-chan (*Batch)) <-chan (*Batch) {
	req, red := int(e.RequiredShards), int(e.RedundantShards)
	out := make(chan (*Batch), 0)

	go func() {
		var err error
		defer func() {
			close(out)
			if err != nil {
				e.errc <- err
			}
		}()
		enc, err := reedsolomon.New(req, red)
		if err != nil {
			return
		}
		// total := req + red
		// if l := len(base); l != req {
		// 	return fmt.Errorf("need %d pieces, but only got %d", req, l)
		// }

		for b := range bb {
			if err = enc.Encode(b.shards); err != nil {
				return
			}
			out <- b
		}
	}()

	return out
}
