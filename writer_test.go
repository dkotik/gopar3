package gopar3

import (
	"bytes"
	"encoding/hex"
	"hash/crc32"
	"testing"

	"github.com/dkotik/gopar3/telomeres"
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
		w := crc32.New(castagnoliTable)
		_, err := w.Write([]byte(tc.RawData))
		if err != nil {
			t.Fatal(err)
		}

		dehex, err := hex.DecodeString(tc.HexCRC32)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(w.Sum(nil), dehex) {
			t.Fatal("did not match crc32.Castagnoli check")
		}
	}
}

func TestWriteShardsWithTagAndChecksum(t *testing.T) {
	testCases := [...]string{
		"1",
		"22",
		"333",
		"4444",
		"55555",
		"6666:66666",
	}

	b := &bytes.Buffer{}
	tlm, err := telomeres.NewEncoder(b, 4)
	if err != nil {
		t.Fatal(err)
	}
	w, err := NewWriter(tlm, NewSequentialTagger(Tag{}, 5))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		if _, err = w.Write([]byte(tc)); err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("%q", b.String())
	// t.Fatal("check result")
}
