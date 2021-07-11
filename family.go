package gopar3

// ShardFamily is a ReedSolomon collection, from which data can be recovered. The median value of the accumulated meta properties is taken as truth.
type ShardFamily struct {
	shards                                [][]byte
	accumulatedRequiredShardsMetaRecords  []uint8
	accumulatedRedundantShardsMetaRecords []uint8
	accumulatedPaddingLengthMetaRecords   []uint16
}

// Load associates a Slice with a family.
func (s *ShardFamily) Load(slice *Slice) (err error) {

}

// OriginalBytes returns the data that was kept safe by the ShardFamily.
func (s *ShardFamily) OriginalBytes() ([]byte, error) {
	return nil, nil
}
