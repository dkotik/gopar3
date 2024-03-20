# GoPar3: File Integrity for the Paranoid

Go implementation of self-healing archive manager with Reed-Solomon error correction inspired by old Par1 and Par2 tools.

## Roadmap to Beta

- [ ] Reap ideas from <https://jacobfilipp.com/arvid-vhs/>
    - > As best as I could figure out, data was stored to tape using Non-Return-to-Zero encoding (according to some Fido7 comments). On the tape itself, files were recorded in sections separated by 5-second blank intervals.
    - > Home-use VHS tapes had worse quality than commercial-grade magnetic tape, so the makers took extra measures to detect and fix errors on tape. (already doing some of this)
    - > ArVid read and wrote data using an error correction algorithm called “Reed-Solomon with Interleaving” (I also came across mentions of a Galois algorithm). They claimed that this let the ArVid software correct up to 3 defective bytes in a code group, and a loss of up to 450 consecutive bytes could be corrected. After reading data from tape, the software performed a CRC32 check for errors, operating on every 512-byte block.
- [ ] Make sure small shard sizes can accommodate [gopar3.BlockLimit] for a given source file size

## Telomeres

GoPar3 uses a telomere encoder to guard block boundaries. Telomeres are repetitions of ":" padding characters. Occurrences of ":" and "\\" within the block data are escaped using "\\". The telomere encoder helps preserve block boundaries in severely damaged files. Even if some blocks are thrown out of alignment by shortening, they can be isolated from healthy blocks and partially recovered.

## Index Inspection

GoPar3 can produce the list of all shards in given sources. The index file may be used to attempt manual restoration of damaged shards.
