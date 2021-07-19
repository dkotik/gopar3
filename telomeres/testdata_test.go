package telomeres

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestStoredCases(t *testing.T) {
	f, err := os.Open("testdata.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		encoded := sc.Bytes()
		if !sc.Scan() {
			t.Fatalf("data pair %q is not matched", string(encoded))
		}
		pieces := bytes.Split(sc.Bytes(), []byte(`||`))
		if len(pieces) == 0 {
			t.Fatalf("data %q has not pieces", sc.Text())
		}
		// spew.Dump(pieces)

		decoder := NewDecoder(
			bytes.NewReader(encoded), 4, 1024, 4096)

		comp := &bytes.Buffer{}
		for i := 0; i < len(pieces)-1; i++ { // last piece is empty, so stop sooner
			comp.Reset()
			_, err = io.Copy(comp, decoder)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(comp.Bytes(), pieces[i]) != 0 {
				t.Fatalf("%d: %q is not equal to %q",
					i*2, comp.String(), string(pieces[i]))
			}
		}
		if _, err = io.Copy(comp, decoder); err != ErrEndReached {
			if err == nil {
				t.Fatalf("error is nil, when should be ErrEndReached")
			} else {
				t.Fatalf("error %t %q does not match ErrEndReached", err, err.Error())
			}
		}
	}
}
