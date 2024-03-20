/*
Package telomeres provides an encoder and a decoder that write
chunks of data surrounded by repeated uniform byte sequences
that designate chunk boundaries.

Telomeres amplify data resiliency by preventing two data chunks
from being corrupted by having an error on their boundary.

In addition, telomeres provide a more reliable mechanism
of chunk detection that does not depend on counting bytes
associated with chunk length. The file can suffer more
damage while its chunks remain recognizable.
*/
package telomeres

const (
	// Mark is repeated to form a telomere sequence
	// that indicates data chunk boundary.
	Mark = ':'

	// Escape indicates that the next byte should
	// be treated as raw data.
	Escape = '\\'
)
