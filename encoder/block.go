package encoder

import (
	"bytes"
	"context"
	"io"
)

// var m runtime.MemStats
// runtime.ReadMemStats(&m)
// if m.Alloc > s.memoryUsageLimitInBytes {
//     // wait instead of issuing error!
//     return 0, errors.New("memory limit exceeded")
// }

// TODO: add memory usage checks
func (e *Encoder) issueBlocks(
	ctx context.Context, r io.Reader, stream chan<- (*bytes.Buffer),
) (err error) {
	// read blocksize*RequiredShards of bytes
	// pad to the required length
	limit := int64(e.shardSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		b := &bytes.Buffer{}
		// b.Grow(blockSize * e.RequiredShards) ???
		// ref, err := s.Reserve(b)
		// if err != nil {
		// 	return err
		// }
		_, err = io.CopyN(b, r, limit)
		if err != nil {
			return err
		}
		stream <- b
	}
	return
}
