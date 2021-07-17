package gopar3

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func swapWork(s *Swap) error {
	data := []byte(fmt.Sprintf("%d", 9999999999+rand.Intn(999999999999)))
	ref, b, err := s.Reserve()
	if err != nil {
		return err
	}
	_, err = io.Copy(b, bytes.NewReader(data))
	if err != nil {
		return err
	}
	r, err := s.Retrieve(ref)
	if err != nil {
		return err
	}
	check := &bytes.Buffer{}
	_, err = io.Copy(check, r)
	if err != nil {
		return err
	}
	s.Release(ref)
	if bytes.Compare(data, check.Bytes()) != 0 {
		return fmt.Errorf("data %q is not equal to %q", string(data), check)
	}
	return nil
}

func TestSwap(t *testing.T) {
	s := NewSwap(0, 99999999)
	g := &errgroup.Group{}

	for i := 0; i < 10; i++ {
		g.Go(func() error {
			// spew.Dump(s)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(7)))
			return swapWork(s)
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	// spew.Dump(s)
}
