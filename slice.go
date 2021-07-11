package gopar3

// SliceType provides a hint as to how the Slice should be processed.
type SliceType uint8

const (
	// SliceSizeLimit determines when the slice constructor stops retaining data. It is calculated by adding Shard size, one byte for the number of required shards, one byte for the number of redundancy shards, two bytes for the padding length (by which the recovered data is truncated), and the Hash size.
	SliceSizeLimit = blockSize + ReedSolomonMetaBinaryLength + blockHashSize
)

// Slice can contain either a shard, a shard family boundary, or a tail marker.
type Slice struct {
	Body   [SliceSizeLimit]byte // contents buffer
	Index  int                  // the cursor position where the slice began
	Length int                  // contents length
}

// WriteChecksum adds a checksum to the body. If there is not enough space in the buffer, the Slice will overrun in length and will be marked as invalid.
func (s *Slice) WriteChecksum() {
	s.Write(Hash(s.Body[:s.Length]))
}
