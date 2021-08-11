package decoder

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/dkotik/gopar3"
)

const (
	shards = 60
	cutoff = gopar3.TagChecksumPosition
)

func makeShards(tag *gopar3.Tag) [][]byte {
	proto := tag.Prototype()
	sl := make([][]byte, shards)
	var shardn uint8
	var block uint32
	for i := 0; i < shards; i++ {
		b := make([]byte, cutoff)
		if n := copy(b, proto[:]); n != cutoff {
			panic("could not fit all the tag")
		}
		sl[i] = b
		shardn = uint8(i) % (tag.RequiredShards + tag.RedundantShards)
		// fmt.Println(shardn)
		proto.SetShardSequence(shardn)
		if shardn == 0 {
			proto.SetBatchSequence(block)
			block++
		}
	}
	return sl
}

func TestSort(t *testing.T) {
	// generate some tags using gopar3.Protype
	tag := &gopar3.Tag{
		Version:         gopar3.VersionByte,
		RequiredShards:  9,
		RedundantShards: 3,
	}
	if err := tag.Differentiate(); err != nil {
		t.Fatal(err)
	}

	original, shuffled := makeShards(tag), makeShards(tag)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(shards, func(i int, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// try to sort the slice
	SortAscending(shuffled)

	// encode and compare output
	if shardsToString(original) != shardsToString(shuffled) {
		t.Fatal("shuffled slice is still different from the original")
		// t.Fatal(shardsToString(original), shardsToString(shuffled))
	}
}

func shardsToString(s [][]byte) string {
	result := &strings.Builder{}
	result.WriteRune('[')

	for i := 0; i < len(s); i++ {
		result.WriteRune('\n')
		fmt.Fprintf(result, "%x", s[i][gopar3.TagBatchSequencePosition:])
	}

	result.WriteRune(']')
	return result.String()
}
