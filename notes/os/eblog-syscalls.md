# Starting Systems Programming, pt 2 — Syscalls & Files

Source: https://eblog.fly.dev/startingsystems2.html (S4)
Date: Oct 2026
Related: [[binary-formats-my-summary]] (execve loads a Mach-O/ELF; libSystem = syscall gateway)

## ⭐ To capture when I write the S4 summary (reminders)

- **BEFORE the S4 summary: quiz Prabhu specifically on the BIT-MANIPULATION intuition** with
  concrete worked examples — extract a field `(x >> p) & ((1<<n)-1)`, pick the right mask
  (`2ⁿ-1`: 1/3/7/15/0xFF), build via shift-OR `(d0<<6)|(d1<<3)|d2`, AND=read/extract vs
  OR=set/build, `s[i]-'0'` digit trick, octal↔bits. Use made-up mode values and ask him to
  decode/encode by hand.

- **chmod exercise (§5) — several sharp gotchas to quiz:**
  - **`Errno` vs `error` nil trap** — `syscall.Syscall`'s 3rd return is `syscall.Errno`; assigning
    `Errno(0)` into an `error` variable makes it **non-nil** (boxed zero-value interface) → always
    "errors". Check `errno != 0`, don't box it.
  - **C-string null-termination** — `open`/`chmod`/`stat` take a path **pointer only, no length**,
    so the kernel reads to a `\0`. Go `[]byte`/strings are **(ptr,len), no terminator** → must
    `append(..., 0)` or `BytePtrFromString`. Length-prefix (Go) vs sentinel (C string) — same
    choice as the `\n` line sentinel + length-prefix tradeoff in [[bytes-and-numbers]].
  - **`string(int)` rune trap** — `string(644)` = the char at code point 644, NOT `"644"`
    (use `strconv`). And `s[i]` is a **byte**, `'0'` is a number (code point), so `s[i]-'0'` is
    arithmetic on code values (digits are consecutive). **No implicit numeric conversions in Go** —
    hence the `int(...)` to line up `byte` with `int`.
  - symbolic↔octal: build with shift-OR, extract with shift+mask (two directions of the same bits).

- **After the summary: Prabhu wants a QUIZ across ALL core concepts so far (S2–S4), not just
  S4.** Spans: what systems programming is (own-the-byte-layout) · bytes/numbers (shifts,
  serialization, endianness) · UTF-8/Unicode (3 layers, decode, bytes-don't-carry-encoding) ·
  binary formats (ELF/Mach-O, static-vs-dynamic linking, instruction pointer, endianness) ·
  file I/O & memory (ReadAll vs streaming vs mmap, why DBs avoid mmap, buffer sizing) · syscalls
  (uintptr vs unsafe.Pointer, fd table, read/write loop + partial writes, blocking read,
  fork+execve, env vars as a list) · file permissions/types (9-bit rwx, octal-vs-decimal,
  bit-masking, S_IFMT file types, "everything is a file"). Draw from the notes/*.md + my-summary
  files.

- **`uintptr` in the syscall interface** — it is NOT a pointer; it's a pointer-*sized* unsigned
  integer. Everything crossing into the kernel via `syscall.Syscall` is reduced to a single
  machine-word **number**: either the value itself (fd, count, flags, bool→0/1) or the **address**
  where bigger data lives (strings/buffers travel as an address + length, never by value). Key
  click: **an address is just a (big) integer**, so one type — `uintptr` — covers both "a plain
  integer" and "an address," because they're both numbers. The slot type is sized for the biggest
  occupant (an address), so even the tiny fd `3` rides in a `uintptr`. (§4.1.1)

- **File permissions & types (§5)** — the `mode` is one integer: low 9 bits = `rwx` for
  owner/group/other (bit-packed flags, tested with single-bit masks `mode & (1<<(8-i))`), above
  them the special bits (setuid/setgid/sticky), then the file-type field (`S_IFMT`). **Octal**
  because 3 bits = 1 octal digit = 1 `rwx` triple (perfect alignment — contrast `shexdump` where
  hex fits a byte's two nibbles). 7 file types ("everything is a file": regular/dir/symlink/char
  dev/block dev/FIFO/socket). Decode symbolic form from the raw `st.Mode` (via `syscall.Stat`)
  yourself — don't lean on `.Mode().String()`. (§5.0.1 printperms)

- **`unsafe.Pointer` vs `uintptr`** — same *size* and convertible, but opposite meaning to the
  **GC**: `unsafe.Pointer` = a **real, tracked pointer** (GC keeps the object alive, updates it if
  moved); `uintptr` = just a **number**, a dead snapshot of an address the GC ignores. Danger:
  once you convert a pointer to `uintptr`, the GC can free/move the object and your number goes
  stale. Golden rule in syscalls: write `uintptr(unsafe.Pointer(p))` **in the call's argument
  list** — never store the `uintptr` in a variable and use it later (storing the `unsafe.Pointer`
  is fine). (§4.2)

## Questions (jot the moment they occur — don't stop to research)

- ~~Does `uintptr(fd)` turn the fd into a pointer to its data?~~ **No.** fd is an opaque
  **index** into the kernel's per-process file-descriptor table; `uintptr(3)` is just the number
  `3`. At syscall time the **kernel** uses it to look up its own `struct file` (offset, mode,
  inode link) — kernel-private, never a user-space address. (§4.1.1)
