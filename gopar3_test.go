package gopar3

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/klauspost/reedsolomon"
)

func TestCRC32(t *testing.T) {
	testCases := [...]struct {
		RawData  string
		HexCRC32 string
	}{ // https://crccalc.com/
		{RawData: "123456789", HexCRC32: "E3069283"},
		{RawData: "This site uses cookies for analytics and ads.", HexCRC32: "7EC18835"},
		{RawData: "9f834hkjfnfo783ofhlkdfh7834", HexCRC32: "FBA9ED08"},
	}

	for _, tc := range testCases {
		b := &bytes.Buffer{}
		w := NewCRC(b)
		_, err := w.Write([]byte(tc.RawData))
		if err != nil {
			t.Fatal(err)
		}
		if err = w.Close(); err != nil {
			t.Fatal(err)
		}

		c := bytes.NewBuffer([]byte(tc.RawData))
		dehex, err := hex.DecodeString(tc.HexCRC32)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = c.Write(dehex)

		if !bytes.Equal(b.Bytes(), c.Bytes()) {
			t.Fatal("did not match crc32.Castagnoli check")
		}
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
