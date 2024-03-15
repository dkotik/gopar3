package scanner

import (
	"fmt"
)

// SnifferSample tracks the frequency of a shard and others similar to it.
type SnifferSample struct {
	Popular   []byte
	Frequency uint8
}

// Sniffer applies Differentiator to group collected samples, so that the most popular type can be selected later.
type Sniffer struct {
	Differentiator func(shard []byte) (group string)
	Samples        map[string]*SnifferSample
}

// Sample notes a shard for the coming selection by popularity. Shards are grouped using the Differentiator.
func (s *Sniffer) Sample(shard []byte) {
	group := s.Differentiator(shard)
	if saved, ok := s.Samples[group]; ok {
		saved.Frequency++
		return
	}
	s.Samples[group] = &SnifferSample{
		Popular:   shard,
		Frequency: 1,
	}
}

// GetPopular determines predominant shard qualities by taking the most popular sampled values grouped by Differentiator.
func (s *Sniffer) GetPopular() *SnifferSample {
	top := &SnifferSample{
		Popular:   nil,
		Frequency: 0,
	}
	for _, v := range s.Samples {
		if top.Frequency < v.Frequency {
			top = v
		}
	}
	return top
}

func TagDifferentiator(shard []byte) (group string) {
	length := len(shard)
	return fmt.Sprintf("%x^%d", shard[length-24:length-14], length) // TODO: update with tag positions after shard is refactored in
}
