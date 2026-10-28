# Go raw-syscall gotchas (reference)

Lookup notes for when I do raw `syscall.Syscall` work again (surfaced in the §5 `chmod` exercise,
S4). These are **implementation footguns, not concepts** — concepts live in [[syscalls-my-summary]].

## 1. `Errno` vs `error` — the non-nil-interface trap
`syscall.Syscall` returns `(r1, r2 uintptr, errno syscall.Errno)` — the 3rd value is a
**`syscall.Errno`**, not an `error`. If you assign it into an `error` variable, `Errno(0)`
(success) boxes into a **non-nil** interface, so `err != nil` is **always true** → you "error" even
on success.
**Fix:** keep it typed as `Errno` and compare to `0`:
```go
_, _, errno := syscall.Syscall(...)
if errno != 0 { ... }
```

## 2. C-string null-termination (pointer-only args)
Path args to `open`/`chmod`/`stat` are passed as a **pointer with no length**, so the kernel reads
bytes until a `\0`. Go strings/`[]byte` are **(pointer, length)** with **no `\0`** → the kernel
reads past your data into garbage (undefined behavior; may "work" by luck).
**Fix:** null-terminate — `append([]byte(path), 0)`, or `syscall.BytePtrFromString(path)`, or use
the `syscall.Chmod` / `os.*` wrapper (does it for you).
(Length-prefix [Go] vs sentinel [C string] — same choice as the `\n` line sentinel. Contrast
`read`/`write`, which pass a pointer **and** a length, so they need no terminator.)

## 3. `string(int)` rune trap
`string(644)` is **not** `"644"` — it's the character at **code point 644** (a rune → its UTF-8
bytes). `string(someInt)` always means "the character at this code point," never the number's text
form.
**Fix:** for number→text use `strconv.Itoa` / `strconv.FormatInt`. (This bit the first
`parseOctetToDecimal`, which did `string(octet)` on an `int` expecting `"644"`.)

## 4. No implicit numeric conversions
Go never auto-converts between numeric types — `byte`, `int`, `uint32`, `uintptr` must be converted
**explicitly**. That's why bit code is littered with `int(...)`, `uint32(...)` (e.g.
`int(s[i]-'0')` to line a `byte` up with an `int` for the `|`). A type mismatch is a **compile
error**, not a silent promotion (unlike C).
