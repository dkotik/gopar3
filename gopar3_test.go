package gopar3

import (
	"bytes"
	"context"
	"hash/crc32"
	"testing"
	"time"

	"github.com/klauspost/reedsolomon"
)

func TestInflate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := Inflate(
		ctx,
		t.TempDir(),
		"README.md",
		5,
		3,
		64,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestIdenticalQuorumShardsWithDifferentParity(t *testing.T) {
	quorum := [][]byte{
		[]byte("aaa"),
		[]byte("bbb"),
		[]byte("ccc"),
		[]byte("ddd"),
		[]byte("eee"),
	}

	enc, err := reedsolomon.New(len(quorum), 3)
	if err != nil {
		t.Fatal(err)
	}
	threeMore := append(quorum, nil, nil, nil)
	err = enc.Reconstruct(threeMore)
	if err != nil {
		t.Fatal(err)
	}

	enc, err = reedsolomon.New(len(quorum), 5)
	if err != nil {
		t.Fatal(err)
	}
	fiveMore := append(quorum, nil, nil, nil, nil, nil)
	err = enc.Reconstruct(fiveMore)
	if err != nil {
		t.Fatal(err)
	}

	for i, shard := range threeMore {
		if !bytes.Equal(shard, fiveMore[i]) {
			t.Logf("threeMore: %q", shard)
			t.Logf("fiveMore: %q", fiveMore[i])
			t.Fatal("shards do not match")
		}
	}
}

func TestIfCastagnoliSumAlwaysSame(t *testing.T) {
	crc := crc32.New(castagnoliTable)
	_, err := crc.Write([]byte("a"))
	if err != nil {
		t.Fatal(err)
	}

	sum := crc.Sum32()
	crc.Sum(nil)
	sum2 := crc.Sum32()

	if sum != sum2 {
		t.Fatal("unequal sums")
	}
}
