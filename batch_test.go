package gopar3

import (
	"bytes"
	"io"
	"testing"
)

func TestBatchLoader(t *testing.T) {
	b := bytes.NewBufferString("12345678901234567890")
	total := b.Len()
	loader := &BatchLoader{
		Quorum:    3,
		Shards:    5,
		ShardSize: 19,
	}

	batch, loaded, err := loader.Load(b)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if len(batch) != loader.Shards {
		t.Fatal("wrong number of shards:", len(batch), "vs", loader.Shards)
	}
	if loaded != total {
		t.Fatal("did not load the correct number of bytes:", loaded, "vs", total)
	}

	for _, shard := range batch[:loader.Quorum] {
		t.Logf("shard: %q", shard)
		if len(shard) != loader.ShardSize {
			t.Fatal("size does not match expectation:", len(shard), "vs", loader.ShardSize)
		}
	}

	for _, shard := range batch[loader.Quorum:loader.Shards] {
		if shard != nil {
			t.Logf("shard: %q", shard)
			t.Fatal("parity shard is not empty")
		}
	}
}

func TestBatchPadding(t *testing.T) {
	b := bytes.NewBufferString("a")
	total := b.Len()
	loader := &BatchLoader{
		Quorum:    3,
		Shards:    5,
		ShardSize: 8,
	}

	batch, loaded, err := loader.Load(b)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if len(batch) != loader.Shards {
		t.Fatal("wrong number of shards:", len(batch), "vs", loader.Shards)
	}
	if loaded != total {
		t.Fatal("did not load the correct number of bytes:", loaded, "vs", total)
	}

	for _, shard := range batch[:loader.Quorum] {
		t.Logf("shard: %q", shard)
		if len(shard) != loader.ShardSize {
			t.Fatal("size does not match expectation:", len(shard), "vs", loader.ShardSize)
		}
	}

	for _, shard := range batch[loader.Quorum:loader.Shards] {
		if shard != nil {
			t.Logf("shard: %q", shard)
			t.Fatal("parity shard is not empty")
		}
	}

	if !bytes.Equal([]byte("a???????"), batch[0]) {
		t.Fatal("first shard was not filled with padding bytes")
	}

	for _, shard := range batch[1:loader.Quorum] {
		if !bytes.Equal([]byte("????????"), shard) {
			t.Fatal("first shard was not filled with padding bytes")
		}
	}
}
