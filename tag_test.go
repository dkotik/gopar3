package gopar3

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestTagLimits(t *testing.T) {
	if TagSize != 16 {
		t.Fatal("tag size bent out of standard", TagSize)
	}
	if DifferentiatorSize != 13 {
		t.Fatal("differentiator size bent out of standard", DifferentiatorSize)
	}
	if ShardLimit != math.MaxUint8 {
		t.Fatal("unexpected shard limit", ShardLimit, math.MaxUint8)
	}
	if ShardBatchLimit != math.MaxUint16 {
		t.Fatal("unexpected shard limit", ShardBatchLimit, math.MaxUint16)
	}
	if SourceSizeLimit != math.MaxUint64 {
		t.Fatal("unexpected shard limit")
	}

	a := []byte("1234567890123456")
	b := NewTagFromBytes(a).Bytes()
	if !bytes.Equal(a, b) {
		t.Logf("a: %q", a)
		t.Logf("b: %q", b)
		t.Fatal("values do not match")
	}
	if !reflect.DeepEqual(NewTagFromBytes(a), NewTagFromBytes(b)) {
		t.Fatal("decoded tags do not match")
	}
}
