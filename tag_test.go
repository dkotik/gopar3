package gopar3

import "testing"

func TestTag(t *testing.T) {
	if TagSize != 24 {
		t.Fatal("tag size got bent out of standard", TagSize)
	}
	if 2<<40/MaxBlocks < 512 { // 512b will be the minimum block size for a TB of data
		t.Fatal("tag size does not support a TB of input data", 2<<40/MaxBlocks)
	}
}

func TestPadding(t *testing.T) {
	if MaxPadding != 4294967295 {
		t.Fatal("max padding is out of standard", MaxPadding)
	}
}
