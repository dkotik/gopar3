package shard

import "testing"

func TestTag(t *testing.T) {
	if TagSize != 24 {
		t.Fatal("tag size got bent out of standard", TagSize)
	}
}
