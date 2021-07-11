# GoPar3 Candidate Alpha

Go implementation of self-healing archive manager with Reed-Solomon error correction inspired by old Par1 and Par2 tools.

## Telomeres

GoPar3 uses a telomere encoder to guard block boundaries.

## Cross-check hashing

GoPar3 injects additional check sums between blocks at a regular interval.
