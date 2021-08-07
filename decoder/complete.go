package decoder

import (
	"github.com/klauspost/reedsolomon"
)

// CompleteWithReedSolomon generates missing redundant shards.
func (d *Decoder) CompleteWithReedSolomon(in <-chan ([][]byte)) <-chan ([][]byte) {
	req, red := int(d.requiredShards), int(d.redundantShards)
	out := make(chan ([][]byte), 0)

	go func() {
		var err error
		defer func() {
			close(out)
			if err != nil {
				d.errc <- err
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

		for b := range in {
			if err = enc.ReconstructData(b); err != nil {
				return
			}
			out <- b[:req]
		}
	}()

	return out
}
