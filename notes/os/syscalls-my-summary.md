# My summary: syscalls, processes & permissions (eblog pt 2, §3–5)

My own recap (S4, Oct 2026), in my words, touched up for accuracy after Claude's gap-check.
Full source thread + reminders: [[eblog-syscalls]]. Go implementation footguns kept separately:
[[go-syscall-gotchas]]. Related: [[binary-formats-my-summary]] (execve loads a Mach-O image;
`libSystem` = the syscall gateway).

## Environment variables & process creation (§3)

Environment variables are passed as a **list** (`[]string` of `KEY=value`) to a process. You can
add them in the shell with `export A=value`; the **shell de-dups**, so setting a key again
overwrites it — no duplicates when you go through the shell. But you **can** end up with duplicate
env vars if a **non-shell** parent hand-builds the env when creating a new process.

Running a program (say a Python script) is a **new process**, and the shell makes it happen with
**fork + execve**:
- **fork** clones the current process (the shell itself). You fork *first* so the **shell
  survives** — otherwise it would become the new program and never return.
- **execve** then replaces the clone's whole image: it tears down the address space and **loads
  the new program from disk** (maps its code/data, fresh heap + stack), and finally points the
  instruction pointer at the new program's entry. The shell passes its env vars *into* execve.

## System calls (§4)

When a program needs to do I/O (or anything privileged), it makes a **system call**. The reason is
**security / isolation**: the CPU has two privilege levels — *user mode* (restricted: can't touch
hardware or other processes' memory) and *kernel mode* (full access). A syscall is the **guarded
doorway** between them, so a process can't reach outside its own memory; the kernel does the
operation on its behalf and enforces the boundary.

On a syscall we leave user mode and hand control to the kernel: **save where we are** (registers,
stack, etc.), switch to kernel mode, the OS does the work (open a file → returns a file
descriptor; read → fills the buffer we passed as a pointer, up to the size we asked), then control
returns to user mode.

**Passing arguments (`uintptr`):** everything we pass to a syscall is either an **integer** or an
**address/pointer** — and an address is just a number too. `uintptr` is the unsigned, pointer-sized
integer type that holds **both**. All syscall args are this type because there's **one common
`syscall.Syscall` function for every syscall**, so the slots use the one type that covers
everything (same at the ABI level regardless of language). To pass a **string/array**, take its
address via `unsafe.Pointer`, then convert that to `uintptr`.

**File descriptors** are per-process. Defaults: **0 = stdin, 1 = stdout, 2 = stderr**. Anything you
open yourself gets a new fd (3+) — e.g. `SYS_OPEN` returns one, which you then pass to `read`. An fd
is an **index into the kernel's per-process table**, not a memory address.

**The cat loop:** open a file, `read` ~1024 bytes into a buffer, `write` them to stdout, loop until
EOF. This is how you process a 16GB file without 16GB of RAM — you **stream** it in fixed chunks, so
memory stays constant. Buffer size is a performance knob: bigger = **fewer syscalls = faster**
(syscalls are expensive), ideally around the **page size (4096)**; too small = too many syscalls.

**Blocking read:** a blocking `read` **sleeps that thread** until the data comes back (the scheduler
runs other threads meanwhile). There is **one kernel shared across all processes** — it runs on
whichever core made the syscall.

## Permissions (§5)

Each file stores a **mode integer** plus **who owns it (UID)** and **what group it belongs to
(GID)** — so at permission-check time the kernel picks **one** triple by priority: **owner if your
UID matches → else group if your GID matches → else other**.

The permission part of the mode is **octal**: a 3-digit number where each digit is one group's
permissions — **1st digit = owner, 2nd = group, 3rd = other** — each `0–7`. Octal fits because
permissions are grouped in **3s**, and **3 bits = exactly 1 octal digit = 1 `rwx` triple** (9 bits
total for the three groups).

Each digit in binary gives the three actions **r, w, x**: `7 = 111 = rwx`, `5 = 101 = r-x`
(read + execute, no write), `0 = ---` (nothing). So `755` → `rwxr-xr-x`.

Above the permission digits, the **leading part of the mode is the file type** (e.g. a regular file
with 755 is `0o100755`). There are several file types (regular, directory, symlink, char/block
device, FIFO, socket) — the keeper is the principle **"everything is a file"**: Unix exposes
devices/pipes/sockets *as files* you use with the same syscalls. (Didn't memorize the full list —
it's lookup trivia; the principle + "first part = file type" is enough.)

**Decoding the mode (extract):** to get a group's digit, **shift right** to bring it to the bottom,
then **AND with 7** (`111`) to keep just those 3 bits, then read the value:
- owner: `mode >> 6 & 7`
- group: `mode >> 3 & 7`
- other: `mode >> 0 & 7`

**AND extracts, OR sets.** The inverse direction — **building** the mode integer from the digits —
uses **shift-OR**: `(owner << 6) | (group << 3) | other` (same idea as parsing an octal string:
`result = result<<3 | digit`). Shift-right + AND reads a field; shift-left + OR writes one.
