package gopar3

import "bytes"

type swapReference uint16

// another option is mmap backing https://github.com/edsrzf/mmap-go

type Swap struct {
	buffers map[swapReference]*bytes.Buffer
	cursor  swapReference
	next    chan (swapReference)
}

func NewSwap(shardSizeInBytes uint16, memoryUsageLimitInBytes uint64) (*Swap, error) {
	s := &Swap{}

	// can calculate precisely since telomeres
	maxShards := memoryUsageLimitInBytes / uint64(shardSizeInBytes)

	s.next = make(chan (swapReference), s.limit)
	var i swapReference
	for ; i < swapReference(maxShards); i++ {
		s.buffers[i] = &bytes.Buffer{}
		s.buffers[i].Grow(shardSize)
		s.next <- i
	}

	return s, nil
}

func (s *Swap) NextBufferAvailable() <-chan (swapReference) {
	return s.next
}

// Release resets and puts the buffer back into s.next?
func (s *Swap) Release(ref swapReference) {
	s.buffers[ref].Reset()
	s.next <- ref
}
