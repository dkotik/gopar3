package telomeres

import "log"

// Cursor returns the underlying reader position.
// Data chunk boundary can be determined by checking the
// cursor position before and after calling [Decoder.Stream].
func (d *Decoder) Cursor() (position int64) {
	return d.cursor
}

func (d *Decoder) FindEdge(softReadLimit int64) (n int64, err error) {
	var (
		i int
		c byte
	)

	if len(d.window) < d.minimum {
		n, err = d.fillBufferWindow()
		softReadLimit -= n
		n = 0
	}

	for i, c = range d.window {
		if c == d.mark {
			n++
			continue
		}
		if d.telomereStreak+i >= d.minimum {
			// found edge that ended a streak
			d.window = d.window[i:]
			d.cursor += int64(i + d.telomereStreak)
			d.telomereStreak = 0
			log.Println("buffer:", string(d.window))
			log.Println("cursor:", d.cursor, n)
			// log.Println("tstreak:", d.telomereStreak)
			return n, err
		}
		return n, err // do not move cursor - already at edge
	}

	more := 0
	for {
		if err != nil {
			break
		}
		more, err = d.r.Read(d.b)
		for i, c = range d.b[:more] {
			if c != d.mark {
				d.window = d.b[i:]
				d.telomereStreak += i
				d.cursor += int64(i)
				return n, err
			}
			n++
		}
		d.cursor += int64(more)
		softReadLimit -= int64(more)
		// TODO: add soft limit check
	}

	d.window = nil // exhausted read
	d.telomereStreak += int(n)
	d.cursor += n
	return n, err
}
