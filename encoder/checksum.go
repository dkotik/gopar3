package encoder

import (
	"hash"
	"io"
)

type checkSumWriter struct {
	io.Writer
	checksum hash.Hash32
}

func (c *checkSumWriter) Write(b []byte) (int, error) {
	c.checksum.Write(b)
	return c.Writer.Write(b)
}

func (c *checkSumWriter) Cut() error {
	// var b [4]byte
	// binary.BigEndian.PutUint32(b[:], c.checksum.Sum32())
	// _, err := c.Writer.Write(b[:])
	_, err := c.Writer.Write(c.checksum.Sum(nil))
	if err != nil {
		return err
	}
	c.checksum.Reset()
	return nil
}
