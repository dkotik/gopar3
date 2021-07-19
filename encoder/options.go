package encoder

import "errors"

// Option modifies the encoder.
type Option func(e *Encoder) error

// WithSwap sets up a buffer swap provider.
func WithSwap() {}

// WithShardSize sets the size of created chunks. Smaller chunks make the output more resilient at the cost of disk space and recovery speed.
func WithShardSize(inbytes int) Option {
	return func(e *Encoder) error {
		if inbytes < 1 {
			return errors.New("cannot encode using 0 or negative shard size")
		}
		e.shardSize = inbytes
		return nil
	}
}
