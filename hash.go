package gopar3

import (
	"bytes"
	"hash/crc32"
)

// Official documentation says that Koopman is superior for error detection.
var tablePolynomial = crc32.MakeTable(crc32.Koopman)

// ChecksumCompute returns a checksum for error correction.
func ChecksumCompute(in []byte) []byte {
	// TODO: make sure this function does not cause race conditions due to using the table?
	hash := crc32.New(tablePolynomial)
	return hash.Sum(in)[:]
}

func ChecksumValidate(b []byte) bool {
	// TODO: make sure this function does not cause race conditions due to using the table?
	length := len(b)
	if length < blockHashSize {
		return false
	}
	length -= blockHashSize
	return 0 == bytes.Compare(
		b[length:], Hash(b[:length]))
}
