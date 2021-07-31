package shard

import "testing"

func TestPadding(t *testing.T) {
	if MaxPadding != 65535 {
		t.Fatal("max padding is out of standard", MaxPadding)
	}
}
