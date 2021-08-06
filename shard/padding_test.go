package shard

import "testing"

func TestPadding(t *testing.T) {
	if MaxPadding != 4294967295 {
		t.Fatal("max padding is out of standard", MaxPadding)
	}
}
