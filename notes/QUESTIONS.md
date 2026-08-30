# Question parking lot

Questions too deep (or too early) to answer in the session they came up. Swept during weekly reviews and phase-exit sessions. Move answered ones to the bottom with a one-line answer — watching this section grow is a progress signal.

Format: `- [ ] (S# / source) question` → when answered: `- [x] (S# → S#) question — short answer or pointer to notes/diagram`

## Open

## Answered

- [x] (side thread, S2→S3) Does `io.ReadAll` read the whole file into memory, or fetch more as you loop over it? — whole file, eagerly; and `os.ReadFile` doesn't fix files > RAM. Full thread incl. mmap-vs-buffer-pool: `notes/os/go-file-io-memory.md`
