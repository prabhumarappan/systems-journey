# My summary: executable formats, linking, the instruction pointer, endianness

My own recap (S3, Oct 2026), in my words, touched up for accuracy after Claude's gap-check.
Source thread: eblog §5.4 + the instruction-pointer / endianness follow-ons.
Related: [[eblog-programs]], [[text-encoding-unicode-utf8]] (the UTF-8-vs-instruction-length
untangle, and "bytes don't carry their encoding" ↔ "bytes don't carry their endianness").

## Executable file formats (ELF / Mach-O / PE)

ELF and Mach-O are two **executable file formats** (PE is the Windows one). ELF is used by
Linux; **Mach-O by macOS — on both Intel and Apple Silicon.** Each has its own header layout and
structure.

**Format is tied to the OS, not the architecture** — the architecture is just a *field inside*
the header. That's why my binary investigation differed from the blog: the blog used **ELF**
(Linux); I'm on **macOS**, so my executable was **Mach-O**. Different magic bytes, different field
offsets — same hello world.

## After the header: load commands, then the code

The header is followed by **load commands** (I'd first called this "loader data"). Some of them
name the **dynamic/shared libraries** the program needs at runtime — on macOS the big one is
**`libSystem`**, the single gateway for syscalls (writing to stdout, opening a socket, DNS, etc.
all go through it).

Then come the actual segments. **`__TEXT`** holds the program's real machine instructions.
(`__rodata` holds read-only data like string constants — see the "why is hello world 2.3MB"
note in [[eblog-programs]].)

## Static vs dynamic linking (the key distinction)

- **Static — `fmt` (and the whole Go runtime):** *compiled into `__TEXT` at build time.* The
  linker copies the machine code bodily into my binary. **No runtime import step.** Always
  present. This is *why* a Go hello world is ~2.3MB — it carries its own `fmt`, GC, scheduler.
- **Dynamic — `libSystem`:** *not* copied in. It's **named in a load command** and loaded by
  **`dyld`** at launch. The binary just holds a reference.

So "imported when needed" describes the **dynamic** path only. `fmt` is never "imported" — it's
baked in. Trade-off: static = bigger file, zero dependencies, self-contained; dynamic = smaller
file, but the library must exist on the machine at run time.

## The instruction pointer (IP)

- **One IP per thread**, pointing at the instruction it's currently executing.
- The instructions are a **contiguous list**; the CPU runs through them one by one:
  **fetch → execute → advance**.
- **Advance = by the size of the instruction just executed:**
  - **ARM64: fixed 4 bytes** → always +4.
  - **x86-64: variable, 1–15 bytes** → advance by however long that instruction was.
- The UTF-8 untangle: I'd confused this with UTF-8's variable width. The irony is that **x86-64
  is the UTF-8-like one** (variable length, you don't know the size until you read it); ARM64 is
  the fixed-width one.

**Functions (CALL / RET):** a CALL **jumps** to the function *and* **pushes the return address**
(the instruction right after the CALL) **onto the stack**. When the function finishes, **RET pops
that address** and jumps back — so execution resumes at the next instruction after the call.
Code is contiguous, but execution *jumps*.

## Endianness

- **Endianness is a property of the CPU architecture, not the OS or the file format.**
  Apple Silicon (and x86-64) are **little-endian**.
- For a multi-byte integer (2, 4, or 8 bytes):
  - **Little-endian:** the first (lowest-address) byte holds the **least-significant** byte.
  - **Big-endian:** the first byte holds the **most-significant** byte ("network byte order";
    reads naturally in a hex dump).
- **Which order a CPU uses is mostly historical convention** — neither is clearly better. The
  reason a program *stores values in its CPU's native order* is to avoid byte-swapping on every
  read/write (zero cost). The choice of *which* order the chip itself uses is the arbitrary part.

## Why this matters for gokv / godb (the real takeaway)

**In memory**, data is always the CPU's native order — no choice. **On disk / on the wire**, it's
whatever **the format says** — *I* decide when I write.

- If I naively dump native-order bytes, a file written on a little-endian machine is **garbage**
  to a big-endian reader. The bytes don't carry their own endianness (same lesson as "bytes don't
  carry their encoding" in [[text-encoding-unicode-utf8]]).
- So the on-disk format is a **contract, independent of the machine running the code.** Pick one
  byte order, write it in the spec, and have every reader/writer honor it (in Go:
  `binary.LittleEndian` / `binary.BigEndian` on every encode/decode — never dump a raw struct).
- Payoff: **backup/restore and replication across different architectures just work.** Write on
  little-endian node A, restore/replicate to a different machine B → same bytes mean the same
  thing.

**Endianness isn't the only machine-specific trap** in an on-disk format — the same discipline
covers:
- **Fixed field widths** — write `uint32`/`uint64`, never Go's `int` (32- or 64-bit by platform).
- **No struct padding** — never `memcpy` a struct to disk; the compiler inserts invisible padding
  that differs by platform. Serialize field by field.

All one principle: **the on-disk format must be self-describing and architecture-neutral, because
the machine that reads it may not be the one that wrote it.** This is what makes a SQLite file or
a Postgres dump portable — and it's a design decision I'll make for gokv's log-record format.
