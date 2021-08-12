package inspector

import "io"

type Inspector struct {
}

func (i *Inspector) Analyze(r io.Reader) <-chan (*Description) {

	out := make(chan (*Description))
	// go func() {
	// }()
	return out
}
