package gopar3

import (
	"encoding/binary"
	"hash"
	"hash/crc32"
	"io"
)

// Official documentation says that Koopman is superior for error detection.
var crc32PolynomialTable = crc32.MakeTable(crc32.Koopman)

// func ChecksumValidate(b []byte) bool {
// 	// TODO: make sure this function does not cause race conditions due to using the table?
// 	length := len(b)
// 	if length < blockHashSize {
// 		return false
// 	}
// 	length -= blockHashSize
// 	return 0 == bytes.Compare(
// 		b[length:], Hash(b[:length]))
// }

type checkSumWriter struct {
	io.Writer
	checksum hash.Hash32
}

func (c *checkSumWriter) Write(b []byte) (int, error) {
	c.checksum.Write(b)
	return c.Writer.Write(b)
}

func (c *checkSumWriter) Cut() error {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], c.checksum.Sum32())
	_, err := c.Writer.Write(b[:])
	if err != nil {
		return err
	}
	c.checksum.Reset()
	return nil
}
