package gopar3

import (
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/reedsolomon"
)

func TestParityGeneration(t *testing.T) {
	b := bytes.NewBufferString("12345678901234567890")
	loader := &BatchLoader{
		Quorum:    2,
		Shards:    3,
		ShardSize: 2,
	}
	rs, err := reedsolomon.New(loader.Quorum, loader.Shards-loader.Quorum)
	if err != nil {
		t.Fatal(err)
	}

	batches := make(chan [][]byte)
	go func() {
		for {
			batch, _, err := loader.Load(b)
			switch err {
			case nil:
				batches <- batch
			case io.EOF:
				batches <- batch
				close(batches)
				return
			default:
				close(batches)
				t.Fatal(err)
			}
		}
	}()

	out := make(chan [][]byte)
	go func() {
		err = AddParity(out, batches, rs)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// t.Skip("borked")

	got := 0
	for p := range out {
		t.Logf("Got batch: %q %q", p[0], p[1])
		t.Logf("Got parity: %q", p[2])
		got++
	}

	if got != 6 {
		t.Fatal("not the right number of batches:", got, "vs", 6)
	}
}
