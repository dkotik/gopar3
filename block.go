package gopar3

import (
	"hash/crc32"
)

const (
	blockSize        = 512
	blockBufferSize  = blockSize * 2
	blockHashSize    = crc32.Size
	blockDecodeLimit = blockSize + blockHashSize
)

// Block is a minimal unit of the armored file.
type Block struct {
	Body   [blockSize]byte
	Length int
	Index  int
}

// Write ignores any additional data written past the available byte space.
func (b *Block) Write(data []byte) (n int, err error) {
	n = copy(b.Body[b.Length:], data)
	b.Length += n
	return
}

// // IsValid returns true if the hash value matches the body.
// func (b *Block) IsValid() bool {
// 	return 0 == bytes.Compare(
// 		b.Body[blockSize:], ChecksumCompute(b.Body[:blockSize]))
// }

// // Bytes returns the proper slice of bytes representing its contents.
// func (b *Block) Bytes() []byte {
// 	return b.Body[:b.Length]
// }

// // IsCrossCheck returns true if the block is small and contains many encodingBoundaryRunes. Boundary blocks are always invalid, because they do not contain a hash
// func (b *Block) IsCrossCheck() bool {
// 	if b.Length < blockHashSize * 3 {
// 		return false
// 	}
// 	return false
// }

// func (b *Block) String() string {
// 	return fmt.Sprintf("Block#%x@%d-%d",
// 		b.Body[blockSize+metaTagPositionChecksum:blockSize+metaTagPositionChecksum+blockHashSize],
// 		b.Index, b.Index+b.Length) // byte range within the file
// }
