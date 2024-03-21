package gopar3

import (
	"hash/crc32"
	"io"

	"github.com/dkotik/gopar3/telomeres"
)

func WriteShardsWithTagAndChecksum(
	w *telomeres.Encoder,
	shards <-chan []byte,
	tagger Tagger,
) (err error) {
	if _, err = w.Cut(); err != nil {
		return err
	}
	var (
		n     int
		shard []byte
		tag   []byte
		crc   = crc32.New(castagnoliTable)
	)

	for {
		select {
		case shard = <-shards:
			if shard == nil {
				return nil
			}
			// write body
			n, err = w.Write(shard)
			if err != nil {
				return err
			}
			if n != len(shard) {
				return io.ErrShortWrite
			}

			// write tag
			n, err = w.Write(tagger.Bytes())
			if err != nil {
				return err
			}
			if n != TagSize {
				return io.ErrShortWrite
			}
			if err = tagger.Next(); err != nil {
				return err
			}

			// write checksum
			crc.Reset()
			_, err = crc.Write(shard)
			if err != nil {
				return err
			}
			_, err = crc.Write(tag)
			if err != nil {
				return err
			}
			n, err = w.Write(crc.Sum(nil))
			if err != nil {
				return err
			}
			if n != TagBytesForCRC {
				return io.ErrShortWrite
			}

			if _, err = w.Cut(); err != nil {
				return err
			}
		}
	}
}
