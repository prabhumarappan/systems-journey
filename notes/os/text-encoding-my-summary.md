# My summary: bytes, Unicode, UTF-8, and comparison

My own recap (S3, Oct 2026), in my words, touched up for accuracy after Claude's gap-check.
Full reference with worked examples: [[text-encoding-unicode-utf8]].

## The foundation: everything is bytes

Everything starts with bytes — that's the floor a systems programmer / software engineer works
at. Below that (voltages, silicon) is a different department — the hardware folks. Text is
**stored on disk as bytes, moved over the network as bytes, and most languages keep it in memory
as bytes too.** Some languages (Python) store **code points** instead, which helps with length
and indexing — but that costs extra encode/decode work at the I/O boundary.

## Three layers (each lets us extend independently)

1. **Unicode** — the universal catalog: maps every character / number / symbol / emoji to a
   unique **code point** (a number). One standard number per character, agreed worldwide.
2. **Encoding** (UTF-8, …) — converts **code points ↔ bytes**, both directions
   (encode = code point → bytes; decode = bytes → code points — they're *inverses*, not steps).
3. **Font / glyphs** — the drawing a font provides for a given code point; tells the screen what
   to actually render.

Because they're separate layers, each extends on its own — new characters (Unicode), new
encodings, or new fonts — without disturbing the others.

## History: ASCII → Unicode + UTF-8

Started with **ASCII** — 128 characters (letters, digits, symbols), fitting in 1 byte (really 7
bits). "Extended ASCII" pushed to 256, but that was a **mess of incompatible 8-bit code pages** —
partly *why* Unicode was created, to unify them. As we added the world's scripts and emoji, the
**Unicode catalog grew far past 256** (now 1.1M+ code points) — numbers that no longer fit in one
byte. That's where **UTF-8** came in: an encoding that represents big code points across multiple
bytes. UTF-8 is the default standard today (disk, web, network). *(Two separate things grew: the
catalog = Unicode, and the encoding = UTF-8 — not just "a new encoding.")*

## UTF-8 bit layout

The "how many bytes" info is packed into the high bits of the first byte:
- `0xxxxxxx` → **1 byte** (ASCII) — high bit 0 = standalone char, all 7 other bits are data
- `110xxxxx` → start of a **2-byte** char
- `1110xxxx` → start of a **3-byte** char
- `11110xxx` → start of a **4-byte** char
- `10xxxxxx` → a **continuation byte** (2nd/3rd/4th byte of a multibyte char) — NOT a start byte

So once the first byte starts with `1`, multiple bytes make up the character, and the count of
leading 1s says how many. The **leading byte carries data bits too**, not just the size tag.
Continuation bytes being marked `10` makes UTF-8 **self-synchronizing** — you can tell a start
byte from a continuation byte anywhere, so you can find character boundaries even mid-stream.

## Decoding: bytes → code point

Strip the tag bits off each byte, then **concatenate** the leftover data bits (leading byte's
bits first = the high-order part), and read the combined binary as one number. E.g. é = `c3 a9`:
- `c3` = `110 00011` → data `00011`
- `a9` = `10 101001` → data `101001`
- concat: `00011 101001` = **233** = U+00E9 = é

**Concatenate, don't add** — the lead bits are higher-order (it's a shift-and-OR, same trick as
building an int from bytes in [[bytes-and-numbers]]).

## Same bits, different meaning

The same byte means different things depending on which layer reads it. `0xC3` is `195` as a raw
byte but "2-byte lead" as UTF-8; `0xFF` is `255` raw but *invalid* as UTF-8. A "character" is
never a stored thing — it's an **interpretation computed on demand** from bytes. Underneath it's
always binary; hex/decimal are just notations.

## Comparison

**Python** holds code points (an array), so it compares code points. The real win here isn't
comparison — it's **length and indexing** (`len` counts characters, `s[n]` = the Nth character).
The cost is decode on input / encode on output.

**Go** holds UTF-8 bytes and compares the **bytes directly, no decoding.** It walks byte by byte
(lexicographically); at the first differing byte, the smaller value is "less" — a magnitude
comparison the CPU does bit-by-bit from the high bit (`0` vs `1` at the top settles it). For `==`,
Go checks **length first**, then the bytes.

This works without decoding because **UTF-8 is order-preserving**: sorting by raw bytes gives the
same order as sorting by code points.
- ASCII case (trivial — byte == code point): `a` = `01100001` = 97 < `b` = `01100010` = 98.
- Multibyte case (the clever part): `a` (`61`) sorts before `é` (`c3 a9`) because `61 < c3`,
  matching `97 < 233` — even though `c3 a9` as raw numbers (195, 169) look nothing like 233.
- A 1-byte char (lead `0…`, ≤127) always sorts before a multibyte char (lead `11…`, ≥192),
  decided at the first byte. Two same-length chars are decided at the first differing byte.

## Caveat: 1 character ≠ always 1 code point

`é` can be **one** code point (U+00E9) or **two** (`e` + a combining accent U+0301). They look
identical but compare unequal, so robust text comparison **normalizes** first (NFC / NFD).

## For gokv

A KV store treats keys/values as **opaque bytes** — never decodes. Key comparison = byte
comparison (fast, order-preserving, no decode). Encoding only matters at the edges where text is
actually displayed or processed.
