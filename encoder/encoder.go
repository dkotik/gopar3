package encoder

// Encoder adds data resiliency to its input.
type Encoder struct {
	requiredShards      uint8
	redundantShards     uint8
	shardSize           int // TODO: replace with shard size // int64?
	telomeresLength     int
	telomeresBufferSize int
	crossCheckFrequency uint
	errc                chan (error)
}

// NewEncoder initializes the encoder with options. Default options are used, if no options were specified.
func NewEncoder(withOptions ...Option) (e *Encoder, err error) {
	e = &Encoder{
		errc: make(chan (error)),
	}

	withOptions = append(withOptions, WithDefaultOptions())
	if err = WithOptions(withOptions...)(e); err != nil {
		return nil, err
	}

	return nil, nil
}
