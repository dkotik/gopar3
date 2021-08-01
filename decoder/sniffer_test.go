package decoder

import (
	"crypto/rand"
	"fmt"
	orand "math/rand"
	"testing"
)

func TestSniffDemocratically(t *testing.T) {
	var b [24]byte
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

	popular := fmt.Sprintf("%x", SniffDemocratically(q, 0, 24))
	if popular != original {
		for i := 0; i < len(q); i++ {
			t.Logf("%x", q[i])
		}

		t.Fatal("values do not match", popular, original)
	}
}
