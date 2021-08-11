package inspector

import (
	"strings"

	"github.com/dkotik/gopar3"
)

type Description struct {
	Corrupted bool
	Tag       *gopar3.Tag
}

func (d *Description) String() string {
	w := &strings.Builder{}

	return w.String()
}

func Describe(shard []byte) *Description {

	return nil
}
