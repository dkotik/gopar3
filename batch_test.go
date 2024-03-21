package gopar3

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
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

func TestBatchStreaming(t *testing.T) {
	b := bytes.NewBufferString("12345678901234567890")
	loader := &BatchLoader{
		Quorum:    2,
		Shards:    3,
		ShardSize: 1,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan [][]byte)
	go func() {
		if err := loader.Stream(ctx, c, b); err != nil {
			t.Fatal(err)
		}
	}()

	gotBack := 0
	for batch := range c {
		t.Logf("Got batch: %q %q %q", batch[0], batch[1], batch[2])
		gotBack++
	}
	if gotBack != 11 {
		t.Fatal("got the wrong number of batches:", gotBack, "vs", 11)
	}
}
