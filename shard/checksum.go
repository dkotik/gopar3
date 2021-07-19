package shard

import (
	"hash"
	"hash/crc32"
)

// Official documentation says that Koopman is superior for error detection.
var crc32PolynomialTable = crc32.MakeTable(crc32.Koopman)

// NewChecksum creates a shard checksum that is used for encoding and decoding.
func NewChecksum() hash.Hash32 {
	return crc32.New(crc32PolynomialTable)
}
