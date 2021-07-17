# GoPar3 Candidate Alpha

Go implementation of self-healing archive manager with Reed-Solomon error correction inspired by old Par1 and Par2 tools.

## Telomeres

GoPar3 uses a telomere encoder to guard block boundaries. Telomeres are repetitions of ":" padding characters. Occurrences of ":" and "\\" within the block data are escaped using "\\". The telomere encoder helps preserve block boundaries in severely damaged files. Even if some blocks are thrown out of alignment by shortening, they can be isolated from healthy blocks and partially recovered.

## Cross-check hashing

GoPar3 injects additional check sums between blocks at a regular interval.

## Forensics mode
