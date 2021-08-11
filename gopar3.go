package gopar3

import (
	"hash"
	"hash/crc32"
)

const (
	// Version is the current version.
	Version = "0.0.1"
	// VersionByte is written to the shard tag when encoding.
	VersionByte = 'a'
	// PaddingByte is used to fill up incomplete data slots of Reed Solomon processing.
	PaddingByte = '?'

	// MaximumPossibleSourceFileBytes represents how big of a file gopar3 can encode. It is calculated by multiplying ...
	// add telomereLength
	MaximumPossibleSourceFileBytes = 512 * (2 ^ 16) // finish this
)

// Official documentation says that Koopman is superior for error detection.
var crc32PolynomialTable = crc32.MakeTable(crc32.Koopman)

// NewChecksum creates a shard checksum that is used for encoding and decoding.
func NewChecksum() hash.Hash32 {
	return crc32.New(crc32PolynomialTable)
}
