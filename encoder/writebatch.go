package encoder

import "os"

// BatchWriter consumes batches and commits them to disk.
type BatchWriter func(<-chan (*Batch), chan<- (error))

// NewSingleDestinationWriter puts all the batches sequentially into the same Writer.
func (e *Encoder) NewSingleDestinationWriter(path string) BatchWriter {
	return func(bb <-chan (*Batch), errc chan<- (error)) {
		ow, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0700)
		if err != nil {
			errc <- err
			return
		}
		// defer ow.Close() // closed within the go routine

		go func() {
			var err error
			defer func() {
				ow.Close()
				errc <- err // capture error
			}()

			w := e.NewWriter(ow, e.prototype)
			if _, err = w.t.Cut(); err != nil {
				return
			}

			var i uint32
			var j uint8

			for b := range bb {
				j = 0
				w.TagPrototype.SetPadding(b.padding)
				for _, shard := range b.shards {
					if shard == nil {
						break // channel was closed
					}
					w.TagPrototype.SetBatchSequence(i)
					w.TagPrototype.SetShardSequence(j)
					j++
					if _, err = w.Write(shard); err != nil {
						return
					}
				}
				i++
			}
			_, err = w.t.Cut()
		}()
	}
}
