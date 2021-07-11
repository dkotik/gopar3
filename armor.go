package gopar3

import (
	"io"
)

const (
	encodingEscapeRune     = '\\'
	encodingTelomereRune   = ':'
	encodingTelomereLength = 8
)

// StreamBlocks decodes the first block from the reader and pipes it into the buffer.
func StreamBlocks(out chan<- Block, r io.Reader) (err error) {
	var (
		result        [blockDecodeLimit]byte
		resultIndex   int
		b             [blockBufferSize]byte
		n, j          int
		c             byte
		telomereCount uint8
		sawEscapeRune bool
	)
	for n, err = r.Read(b[:]); err != nil; n, err = r.Read(b[:]) {
		for j = 0; j < n; j++ {
			switch c = b[j]; c {
			case encodingEscapeRune:
				if sawEscapeRune {
					sawEscapeRune = false
					break
				}
				sawEscapeRune = true
				continue
			case encodingTelomereRune:
				if sawEscapeRune {
					sawEscapeRune = false
					continue
				}
				telomereCount++
				continue
			}
			if telomereCount > 0 {
				out <- Block{result, resultIndex}
				resultIndex = 0
				continue
			}
			result[resultIndex] = c
			resultIndex++
		}
	}
	return
}
