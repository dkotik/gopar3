package gopar3

import "bytes"

type swapReference uint16

// another option is mmap backing https://github.com/edsrzf/mmap-go

type Swap struct {
	limit   swapReference
	buffers map[swapReference]*bytes.Buffer
	cursor  swapReference
	next    chan (swapReference)
}

func (s *Swap) NextAvailable() <-chan (swapReference) {
	if s.next == nil {
		s.next = make(chan (swapReference), s.limit)
	}

	var i swapReference
	for ; i < s.limit; i++ {
		s.buffers[i] = &bytes.Buffer{}
		s.next <- i
	}

	return s.next
}

// Release resets and puts the buffer back into s.next?
func (s *Swap) Release(ref swapReference) {
	// delete(s.buffers, ref)
	s.buffers[ref].Reset() // does this eliminate the capacity as well?
	s.next <- ref
}
