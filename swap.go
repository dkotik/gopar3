package gopar3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"runtime"
	"sync"
)

// SwapReference identifies a buffer provided by Swap.
type SwapReference uint32

// another option is mmap backing https://github.com/edsrzf/mmap-go

// Swap manages buffers needed for processing.
type Swap struct {
	mutex                   *sync.Mutex
	buffers                 map[SwapReference]*bytes.Buffer
	cursor                  SwapReference
	defaultBufferCapacity   int
	memoryUsageLimitInBytes uint64
}

// NewSwap sets up a valid swap. Memory usage refers to the total allocated memory used by the program. It is checked every time a new buffer is reserved.
func NewSwap(defaultBufferCapacity int, memoryUsageLimitInBytes uint64) *Swap {
	return &Swap{
		mutex:                   &sync.Mutex{},
		buffers:                 make(map[SwapReference]*bytes.Buffer),
		defaultBufferCapacity:   defaultBufferCapacity,
		memoryUsageLimitInBytes: memoryUsageLimitInBytes,
	}
}

func (s *Swap) nextAvailableReference() (SwapReference, error) {
	// do not call this function without locking the mutex
	// repeat at least once
	for i := s.cursor; i < math.MaxUint32; i++ {
		if _, ok := s.buffers[i]; !ok {
			s.cursor = i
			return i, nil
		}
	}
	s.cursor = 0 // wrap around
	return 0, errors.New("there are no available buffers")
}

// Reserve provides a reference to the next buffer that can be worked on.
func (s *Swap) Reserve() (SwapReference, io.Writer, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Alloc > s.memoryUsageLimitInBytes {
		return 0, nil, errors.New("memory limit exceeded")
	}

	b := &bytes.Buffer{}
	b.Grow(s.defaultBufferCapacity)
	s.mutex.Lock()
	ref, err := s.nextAvailableReference()
	if err != nil {
		ref, err = s.nextAvailableReference() // retry
		if err != nil {
			s.mutex.Unlock()
			return 0, nil, err
		}
	}
	s.buffers[ref] = b
	s.mutex.Unlock()
	return ref, b, nil
}

// Retrieve locates the correct reserved buffer and returns it.
func (s *Swap) Retrieve(ref SwapReference) (io.Reader, error) {
	s.mutex.Lock()
	b, ok := s.buffers[ref]
	s.mutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("swap buffer %q does not exist", ref)
	}
	return b, nil
}

// Release destroys the buffer.
func (s *Swap) Release(ref SwapReference) {
	s.mutex.Lock()
	delete(s.buffers, ref)
	s.mutex.Unlock()
}
