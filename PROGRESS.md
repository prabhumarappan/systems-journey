# Progress: OS → Databases → Distributed Systems

**How this works:** one checkbox = one ~45-min session. No dates, no "behind schedule" — just do the next unchecked box. Minimum viable session = 20 min (read half the assignment, or just review diagrams). Never miss twice in a row.

**To start any session:** open Claude here and say **"next session"**.

---

## Phase 1 — OS foundations (targeted OSTEP track)

### M0 · Setup + Go syscalls warmup
- [x] S1 — Repo scaffolded; `gokv` Go module created *(done with Claude, 2026-07-15)*
- [ ] S2 — Read eblog "Starting Systems Programming" pt 1 + pt 2 → bullets in `notes/os/eblog-syscalls.md` → ritual with Claude
- [ ] S3 — Build: `experiments/syscalls/` — Go program doing raw open/read/write; watch its syscalls (`dtruss`/`fs_usage` on macOS)

### M1 · Virtualization refresh (OSTEP ch 4–7)
- [ ] S4 — Read ch 4 (Processes) + ch 5 (Process API) → notes → ritual
- [ ] S5 — Read ch 6 (Limited Direct Execution), skim ch 7 (Scheduling) → notes → ritual
- [ ] S6 — Build: mini shell in Go — run commands, wait, exit status
- [ ] S7 — Build: mini shell — pipes (`ls | wc -l`) and redirects

### M2 · Memory (OSTEP ch 13–15, 18–19, 21–22)
- [ ] S8 — Read ch 13 (Address Spaces) + ch 14 (Memory API) → ritual
- [ ] S9 — Read ch 15 (Address Translation) → ritual
- [ ] S10 — Read ch 18 (Paging) → ritual
- [ ] S11 — Read ch 19 (TLBs) → ritual
- [ ] S12 — Read ch 21–22 (Swapping: mechanisms + policies) → ritual
- [ ] S13 — Build: `experiments/memory/` — pointer addresses across runs (ASLR), `syscall.Mmap` demo; start `gokv` v0: in-memory map + REPL (get/set/del)

### M3 · Concurrency (OSTEP ch 26–30, 32)
- [ ] S14 — Read ch 26 (Concurrency intro) → ritual
- [ ] S15 — Read ch 27 (Thread API) + map to goroutines → ritual
- [ ] S16 — Read ch 28 (Locks) → ritual
- [ ] S17 — Read ch 29 (Locked Data Structures) → ritual
- [ ] S18 — Read ch 30 (Condition Variables) → ritual
- [ ] S19 — Read ch 32 (Concurrency Bugs) → ritual
- [ ] S20 — Build: `gokv` v1 — TCP server, goroutine per connection, deliberately racy
- [ ] S21 — Build: `gokv` v1 — catch the race with `-race`, fix with RWMutex; benchmark Mutex vs RWMutex vs sharded map

### M4 · Persistence (OSTEP ch 36–37, 39–40, 42–44)
- [ ] S22 — Read ch 36 (I/O Devices) + ch 37 (Hard Disks) → ritual
- [ ] S23 — Read ch 39 (Files & Directories) → ritual
- [ ] S24 — Read ch 40 (FS Implementation) → ritual
- [ ] S25 — Read ch 42 (Crash Consistency + Journaling) → ritual ★ core DB concept
- [ ] S26 — Read ch 43 (Log-structured FS) → ritual ★ this is the bridge to Bitcask/LSM
- [ ] S27 — Read ch 44 (Flash/SSDs) → ritual
- [ ] S28 — Build: `gokv` v2 — append-only log file, replay on startup
- [ ] S29 — Build: `gokv` v2 — fsync experiments (throughput with/without), kill -9 durability test
- [ ] S30 — **Phase 1 exit demo:** concurrent crash-safe KV server; review all Phase 1 diagrams with Claude

---

## Phase 2 — Database internals (checklist expanded when Phase 1 nears done)

- **M5** · Bitcask storage engine (~4 wks) — bitcask-intro.pdf, DDIA ch 3, 15-445 lectures 1–4; build keydir + segments + compaction
- **M6** · Pages, buffer pools, B+trees (~7 wks, hardest stretch) — 15-445 + Database Internals ch 2–6 + build-your-own.org; disk B+tree → `godb`
- **M7** · Transactions, WAL, MVCC (~4 wks) — 15-445 + DDIA ch 7
- **M8** · Tiny SQL layer (optional, ~3 wks)

## Phase 3 — Distributed systems (checklist expanded when Phase 2 nears done)

- **M9** · Fundamentals — Kleppmann lectures + DDIA replication/partitioning; Fly.io dist-sys challenges 1–3
- **M10** · Raft — paper + 6.824 lectures; Fly.io challenges 4–6; own Raft in Go
- **M11** · Capstone — Raft-replicated `godb` (toydb architecture as reference)

---

## Session log

| # | Date | What happened | Next step |
|---|------|---------------|-----------|
| S1 | 2026-07-15 | Scaffolded repo, plan approved, gokv module created | S2: read eblog Starting Systems Programming pt 1–2 |
