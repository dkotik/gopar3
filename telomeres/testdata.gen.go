// +build ignore

package main

import (
	"bytes"
	"math/rand"
	"os"
	"time"

	"github.com/dkotik/gopar3/telomeres"
)

// to make new test data: `go generate testdata.gen.go`
const target = "testdata.txt"

//go:generate go run testdata.gen.go

func main() {
	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	encoder := telomeres.NewEncoder(f, 0, 1024)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		solution := &bytes.Buffer{}
		for j := 3; j < 4+rand.Intn(5); j++ {
			word := rData(100)
			solution.Write(word)
			solution.WriteString(`||`)
			_, err = encoder.Write(word)
			if err != nil {
				panic(err)
			}
			err = encoder.Flush()
			if err != nil {
				panic(err)
			}
			_, err = f.Write(rBoundary(20))
			if err != nil {
				panic(err)
			}
		}
		_, err = f.Write([]byte("\n"))
		if err != nil {
			panic(err)
		}
		solution.Write([]byte("\n"))
		_, err = f.Write(solution.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

func rBoundary(n int) []byte {
	b := make([]byte, 4+rand.Intn(n))
	for i := 0; i < len(b); i++ {
		b[i] = ':'
	}
	return b
}

var runes = []byte(`::::::::::::\\\\\\\\\\\\\abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`)

func rData(n int) []byte {
	b := make([]byte, rand.Intn(n))
	limit := len(runes)
	for i := 0; i < len(b)-1; i++ {
		b[i] = runes[rand.Intn(limit)]
	}
	return b
}
