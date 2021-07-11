package gopar3

import (
	"encoding/binary"
	"sort"
)

const (
	shardLimit = 256 // tied to Block.order and klauspost/reedsolomon limit

	// version|required|redundant|padding|sequence|checksum
	metaTagPositionVersionNumber   = 0
	metaTagPositionSequenceNumber  = metaTagPositionVersionNumber + blockHashSize
	metaTagPositionRequiredShards  = metaTagPositionSequenceNumber + 1
	metaTagPositionRedundantShards = metaTagPositionRequiredShards + 1
	metaTagPositionPaddingLength   = metaTagPositionRedundantShards + 1
	metaTagPositionChecksum        = metaTagPositionPaddingLength + 2
	// MetaTagTotalLength corresponds to the number of tag bytes.
	MetaTagTotalLength = metaTagPositionChecksum + blockHashSize
)

// MetaTag holds the all the neccessary hints to perform data reconstruction.
type MetaTag struct {
	b []byte
}

func vrrpStaticMeta(required, redundant uint8, padding uint16) (vrrp [5]byte) {
	vrrp[0] = byte(11)
	vrrp[1] = byte(required)
	vrrp[2] = byte(redundant)
	binary.BigEndian.PutUint16(vrrp[3:5], padding)
	return
}

func (m *MetaTag) getSequenceNumber() uint8 {
	return uint8(m.b[metaTagPositionSequenceNumber])
}

func (m *MetaTag) setSequenceNumber(n uint8) {
	m.b[metaTagPositionSequenceNumber] = byte(n)
}

func (m *MetaTag) getRequiredShards() uint8 {
	return uint8(m.b[metaTagPositionRequiredShards])
}

func (m *MetaTag) getRedundantShards() uint8 {
	return uint8(m.b[metaTagPositionRedundantShards])
}

func (m *MetaTag) getPaddingLength() uint16 {
	return binary.BigEndian.Uint16(
		m.b[metaTagPositionPaddingLength:metaTagPositionChecksum])
}

// type ReedSolomonMeta struct {
// 	SequenceNumber  uint8
// 	RequiredShards  uint8
// 	RedundantShards uint8
// 	PaddingLength   uint16 // truncate recovered data by this much
// }
//
// // Encode ReedSolomonMeta to a binary array.
// func (r *ReedSolomonMeta) Encode() (b [ReedSolomonMetaBinaryLength]byte) {
// 	b[0] = byte(r.SequenceNumber)
// 	b[1] = byte(r.RequiredShards)
// 	b[2] = byte(r.RedundantShards)
// 	binary.BigEndian.PutUint16(b[3:5], r.PaddingLength)
// 	return
// }
//
// // ReadMeta recovers the meta from the slice.
// func (s *Slice) ReadMeta() ReedSolomonMeta {
// 	return ReedSolomonMeta{
// 		SequenceNumber:  uint8(s.Body[blockSize]),
// 		RequiredShards:  uint8(s.Body[blockSize+1]),
// 		RedundantShards: uint8(s.Body[blockSize+2]),
// 		PaddingLength:   binary.BigEndian.Uint16(s.Body[blockSize+3 : blockSize+5]),
// 	}
// }

func medianUint8(bunch []uint8) uint8 {
	sort.Slice(bunch, func(i int, j int) bool {
		return bunch[i] > bunch[j]
	})
	return bunch[len(bunch)/2]
}

func medianUint16(bunch []uint16) uint16 {
	sort.Slice(bunch, func(i int, j int) bool {
		return bunch[i] > bunch[j]
	})
	return bunch[len(bunch)/2]
}
