package gopar3

import (
	"crypto/rand"
	"fmt"
	orand "math/rand"
	"testing"

	"github.com/dkotik/gopar3/shard"
)

func TestSniffDemocratically(t *testing.T) {
	var b [shard.TagSize + 1]byte
	shuffle := func() {
		_, err := rand.Read(b[:])
		if err != nil {
			t.Fatal(err)
		}
	}
	shuffle()
	original := fmt.Sprintf("%x", b[:])

	q := make([][]byte, 16)
	for i := 0; i < len(q); i++ {
		a := b
		q[i] = a[:]
		if i == 10 {
			shuffle()
		}
	}
	orand.Shuffle(len(q), func(i int, j int) {
		q[i], q[j] = q[j], q[i]
	})

	sniffer := &Sniffer{
		Differentiator: TagDifferentiator,
		Samples:        make(map[string]*SnifferSample),
	}
	for _, v := range q {
		sniffer.Sample(v)
	}

	popular, frequency := sniffer.GetPopular()
	captured := fmt.Sprintf("%x", popular)
	if popular == nil || frequency != 11 { // 10 is the expected number of matches
		t.Fatal("could not determine popular shard", captured, frequency)
	}
	if captured != original {
		for i := 0; i < len(q); i++ {
			t.Logf("%x", q[i])
		}

		t.Fatal("values do not match", captured, original)
	}
}
