package gopar3

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShardLoading(t *testing.T) {
	testCases := [...]struct {
		TelomereEncoded string
		TelomereDecoded string
		Shard           Shard
	}{
		// TODO: fill out
	}

	testFile := filepath.Join(t.TempDir(), "testdata.txt")
	for _, tc := range testCases {
		f, err := os.Create(testFile)
		if err != nil {
			t.Fatal("could not create temporary file:", err)
		}
		if _, err = io.Copy(f, strings.NewReader(tc.TelomereEncoded)); err != nil {
			t.Fatal("failed to write to temporary file:", err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
