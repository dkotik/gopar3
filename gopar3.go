package gopar3

const (
	// Version is the current version.
	Version = "0.0.1"
	// VersionByte is written to the shard tag when encoding.
	VersionByte = 'a'

	// MaximumPossibleSourceFileBytes represents how big of a file gopar3 can encode. It is calculated by multiplying ...
	// add telomereLength
	MaximumPossibleSourceFileBytes = 512 * (2 ^ 16) // finish this
)
