package encoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/dkotik/gopar3"
)

// https://www.scadacore.com/tools/programming-calculators/online-checksum-calculator/

func TestCheckSumWriter(t *testing.T) {
	cases := []struct {
		Checksum string
		Data     string
	}{
		{"021d82ca", "really big mess"},
		{"b56b19f7", "cxnjkd ou04390 sdflksdjf 84u3 skdjflsjdflks"},
		{"aa90ee95", "9032089sdf n8sdf0 8f43u0-8340 jfsp98 f34f lsfsd. sdfkuhusdf"},
	}

	b := &bytes.Buffer{}
	w := &checkSumWriter{b, gopar3.NewChecksum()}
	var c [4]byte
	for _, cs := range cases {
		w.checksum.Reset()
		w.Write([]byte(cs.Data))
		binary.BigEndian.PutUint32(c[:], w.checksum.Sum32())
		if s := fmt.Sprintf("%x", c[:]); s != cs.Checksum {
			t.Fatalf("checksum %q does not match %q", s, cs.Checksum)
		}
	}
}
