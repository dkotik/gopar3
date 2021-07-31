package shard

import "fmt"

const ()

func IsPossible(shards uint8, shardSize uint64) error {
	if more := uint64(shards)*shardSize - MaxPadding; more > 0 {
		return fmt.Errorf("exceeding possible by %d", more)
	}
	return nil
}
