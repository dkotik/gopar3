package decoder

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/dkotik/gopar3"
)

func (d *Decoder) SniffAndSetupFilter(in chan<- ([]byte), streams []io.Reader) error {
	shards := make([][]byte, 0, d.sniffDepth)
	streamCount := uint16(len(streams))
	for i := uint16(0); i < d.sniffDepth; i++ {
		shard, err := d.ReadShard(streams[i%streamCount])
		if err != nil {
			return err
		}
		shards = append(shards, shard)
	}

	// which shard tag signature is most common amoung the sniff set?
	popular, err := SniffDemocratically(shards)
	if err != nil {
		return err
	}

	lengthRequirement := len(popular)
	tagMatcher := make([]byte, gopar3.TagPaddingPosition, gopar3.TagPaddingPosition)
	n := copy(tagMatcher, popular[lengthRequirement-gopar3.TagSize+4:])
	if n != gopar3.TagPaddingPosition {
		return errors.New("could not set up a tag matcher for shard filter")
	}

	// shard filter rejects shards that do not match the most popular sniff set tag signature
	// TODO: should move it into its own function, just syncronize the parameters
	d.shardFilter = func(a []byte) bool {
		if len(a) != lengthRequirement {
			// TODO: broadcast event or use errc like in StartReading
			return false
		}
		pos := lengthRequirement - gopar3.TagSize + 4
		if bytes.Compare(tagMatcher, a[pos:pos+gopar3.TagPaddingPosition]) != 0 {
			// TODO: broadcast event or use errc like in StartReading
			return false
		}
		return true
	}

	for _, shard := range shards {
		if d.shardFilter(shard) {
			in <- shard
		}
	}
	return nil
}

func SniffDemocratically(q [][]byte) ([]byte, error) {
	length := len(q)
	if length <= 5 {
		return nil, errors.New("cannot use less than 5 shards to discover common values")
	}
	type rec struct {
		Index int
		Count int
	}

	// count similar values grouped by a slice and total length
	cc := make(map[string]*rec)
	for i := 0; i < len(q); i++ {
		length := len(q[i])
		if length < gopar3.TagSize+1 {
			continue // skip over short byte slices
		}
		mark := fmt.Sprintf("%x::%d", q[i][length-gopar3.TagSize:length-4], len(q[i]))
		if saved, ok := cc[mark]; ok {
			saved.Count++
			continue
		}
		cc[mark] = &rec{
			Index: i,
			Count: 1,
		}
	}

	// spew.Dump(cc)

	// select most common signature
	var top *rec
	for _, v := range cc {
		if top == nil || top.Count < v.Count {
			top = v
		}
	}

	if top == nil || top.Count <= length/3 {
		return nil, errors.New("there was not even a third of common values")
	}
	return q[top.Index], nil
}
