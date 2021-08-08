package inspector

import (
	"strings"

	"github.com/dkotik/gopar3/shard"
)

type Description struct {
	Corrupted bool
	Tag       *shard.Tag
}

func (d *Description) String() string {
	w := &strings.Builder{}

	return w.String()
}

func Describe(shard []byte) *Description {

	return nil
}
