package telomeres

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
)

func ExampleEncoder_Encode() {
	b := &bytes.Buffer{}
	t, err := New(WithMinimumCount(4))
	if err != nil {
		panic(err)
	}
	e := t.NewEncoder(b)

	_, _ = e.Cut()
	_, _ = e.Write([]byte("hello"))
	_, _ = e.Cut()
	_, _ = e.Write([]byte("world"))
	_, _ = e.Cut()

	fmt.Print(b.String())
	// Output: ::::hello::::world::::
}

func ExampleDecoder_Decode() {
	b := newTestBuffer([]byte("::::hello::::world::::"))
	t, err := New(WithMinimumCount(4))
	if err != nil {
		panic(err)
	}
	d := t.NewDecoder(b)

	s1 := &bytes.Buffer{}
	_, err = d.WriteTo(s1)
	if err != nil {
		panic(err)
	}

	s2 := &bytes.Buffer{}
	_, err = d.WriteTo(s2)
	if err != nil {
		panic(err)
	}

	fmt.Println(s1.String(), "||", s2.String())
	// Output: hello || world
}

func randomBoundary(n int) []byte {
	b := make([]byte, 4+rand.Intn(n))
	for i := 0; i < len(b); i++ {
		b[i] = ':'
	}
	return b
}

func randomData(n int) []byte {
	const runes = `::::::::::::\\\\\\\\\\\\\abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	b := make([]byte, n)
	limit := len(runes)
	for i := 0; i < len(b); i++ {
		b[i] = runes[rand.Intn(limit)]
	}
	return b
}

// testBuffer implements io.ReadWriteSeeker for testing purposes.
type testBuffer struct {
	buffer []byte
	offset int64
}

// Creates new buffer that implements io.ReadWriteSeeker for testing purposes.
func newTestBuffer(initial []byte) io.ReadWriteSeeker {
	return &testBuffer{buffer: initial}
}

func (tb *testBuffer) Bytes() []byte {
	return tb.buffer
}

func (tb *testBuffer) String() string {
	return string(tb.buffer)
}

func (tb *testBuffer) Len() int {
	return len(tb.buffer)
}

func (tb *testBuffer) Read(b []byte) (int, error) {
	available := len(tb.buffer) - int(tb.offset)
	if available == 0 {
		return 0, io.EOF
	}
	size := len(b)
	if size > available {
		size = available
	}
	copy(b, tb.buffer[tb.offset:tb.offset+int64(size)])
	tb.offset += int64(size)
	return size, nil
}

func (tb *testBuffer) Write(b []byte) (int, error) {
	copied := copy(tb.buffer[tb.offset:], b)
	if copied < len(b) {
		tb.buffer = append(tb.buffer, b[copied:]...)
	}
	tb.offset += int64(len(b))
	return len(b), nil
}

func (tb *testBuffer) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = tb.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(tb.buffer)) + offset
	default:
		return 0, errors.New("Unknown Seek Method")
	}
	if newOffset > int64(len(tb.buffer)) || newOffset < 0 {
		return 0, fmt.Errorf("Invalid Offset %d", offset)
	}
	tb.offset = newOffset
	return newOffset, nil
}

func (tb *testBuffer) Reset() {
	tb.buffer = nil
	tb.offset = 0
}
