# Bytes ↔ numbers: shifts, serialization, endianness

Source: S3 side-thread (from eblog §4.7's `2ⁿ` shift trick), Sep 9 2026
Related: [[eblog-programs]] (string blob, hex), [[go-file-io-memory]] (buffer pool)

Reusable systems note — this is the machinery gokv (S32) and godb (M6) are built from.
Feeds directly into "the byte layout IS your interface" (S2 litmus test).

## Left shift `<<`

`x << n == x * 2ⁿ`. Bits slide left n positions, zeros fill in from the right; each shift-by-1
doubles the value (each bit's place value doubles). `>>` is the inverse: `x >> n == x / 2ⁿ`.

```
5      = 0000 0101
5 << 1 = 0000 1010 = 10
5 << 3 = 0010 1000 = 40
```

"Fast multiply" is the least interesting use. In systems code shifts are for **placing bits
into a slot** — assembling/disassembling byte layouts and bit flags:

```go
// pack a flag set
const CanRead  = 1 << 0   // 0001
const CanWrite = 1 << 1   // 0010
perms := CanRead | CanWrite          // 0011
if perms & CanWrite != 0 { ... }     // test one bit  (1<<n = "only bit n set")
```

## Building a multi-byte number from bytes (deserialization)

The outside world (disk, network, other programs) only speaks **bytes**. An `int` in RAM is a
CPU convenience; the moment a number crosses to disk/socket/another-program's-file it must
become a byte sequence, and coming back it's reassembled. That round-trip is **serialization**.

```go
// 4 bytes → uint32, big-endian: each byte shifted into its slot, OR'd together
n := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])

// uint32 → bytes (the reverse: shift down, mask off the rest)
b0 := byte(n >> 24)
b1 := byte((n >> 16) & 0xFF)
```

This is exactly what `encoding/binary`'s `BigEndian.Uint32` does internally.

**Where it's needed:**
- **Length-prefixed records (gokv!):** to store `key→val` in a log you write each length as
  fixed bytes then the data: `[keylen:4][key][vallen:4][val]`. On replay you read 4 bytes →
  reassemble to the number → read exactly that many bytes → repeat. Without turning the bytes
  back into the length integer you can't tell where one record ends.
- **Reading binary formats:** Mach-O/ELF headers, PNG dimensions, zip offsets — every numeric
  field is fixed-width bytes you reassemble. (When `size` reported `__gopclntab` = 498840, it
  reassembled that from header bytes.)
- **Network protocols:** ports (16-bit), packet lengths, sequence numbers — all multi-byte ints
  on the wire.

## Why you can't just cast bytes → int: endianness

CPUs store an integer's bytes in different orders (big- vs little-endian). So a file format or
protocol picks **one** canonical order and everyone assembles/disassembles explicitly with
shifts to agree, regardless of their CPU. That's *why* the shift-and-OR — it pins down "byte 0
is most significant" so a file written on one machine reads right on another.

## How big can the int be? — a design decision

N bytes → `2^(8N)` distinct values → max unsigned = `2^(8N) − 1`.

| prefix width | max value        | note |
|---|---|---|
| 2 bytes `uint16` | 65,535 (~64 KB)  | cheap; caps record size |
| 4 bytes `uint32` | 4,294,967,295 (~4 GB) | 4 bytes overhead **per record** |
| 8 bytes `uint64` | ~1.8×10¹⁹ (~18 EB) | future-proof but wasteful |

Tension: a 4-byte length prefix on a 2-byte value = 200% overhead. Fine for big records,
wasteful for tiny ones. Real systems (Protobuf, SQLite, LevelDB) use a **varint**: small
numbers take 1 byte, big numbers grow — best of both, more complex encode/decode. Choosing the
prefix width (or varint) is a real gokv decision.

## Hex is the notation for all this

Bytes are viewed/written in hex because 16 = 2⁴: 1 byte = exactly 2 hex digits (fixed-width
columns) AND 1 hex digit = 4 bits (read the bit pattern directly). Decimal aligns to neither.
Purely a display choice — the bytes are identical. (Full reasoning in [[eblog-programs]] S3.)

```mermaid
mindmap
  root((bytes ↔ numbers))
    shift <<
      x<<n = x*2^n
      real use: place bits in a slot
      flags 1<<n, mask &0xFF
    serialize
      RAM int is a CPU convenience
      disk/net/file = flat byte stream
      assemble: b0<<24 | b1<<16 | ...
      disassemble: n>>24, &0xFF
    where
      length-prefixed records (gokv)
      binary headers (Mach-O/ELF/PNG)
      network protocols
    endianness
      CPUs disagree on byte order
      format picks one canonical order
      shifts pin down which byte is which
    width = design choice
      2^(8N)-1 max
      overhead vs max size
      varint (protobuf/sqlite)
```
