package gopar3

import (
	"bytes"
	"testing"

	"github.com/klauspost/reedsolomon"
)

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
