package gopar3

const (
	ShardLimit      = 1<<(TagBytesForShardOrder*8) - 1
	ShardBatchLimit = 1<<(TagBytesForShardBatch*8) - 1
	SourceSizeLimit = 1<<(TagBytesForSourceSize*8) - 1
)
