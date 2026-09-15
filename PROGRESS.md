# Progress: OS → Databases → Distributed Systems

**How this works:** one checkbox = one ~45-min session. No dates, no "behind schedule" — just do the next unchecked box. Minimum viable session = 20 min (read half the assignment, or just review diagrams). Never miss twice in a row.

**Sizing rule:** ~3 blog chapters or one OSTEP chapter per session. If an assignment doesn't fit one slot, the session gets split — the session bends, never the slot. Flag it to Claude and it gets resized.

**To start any session:** open Claude here and say **"next session"**.

---

## Phase 1 — OS foundations (targeted OSTEP track)

### M0 · Setup + Go syscalls warmup (eblog "Starting Systems Programming")
- [x] S1 — Repo scaffolded; `gokv` Go module created *(done with Claude, 2026-07-15)*
- [x] S2 — Read pt 1, chapters 1–3 → bullets in `notes/os/eblog-programs.md` → ritual with Claude *(2026-07-16)*
- [x] S3 — Read pt 1, §4.5–4.7 + data/code segment investigations + §5.1 `shexdump`; built `experiments/S3/4.6torsolines.go` (torso→lines exercise); notes in `notes/os/eblog-programs.md` + new `notes/os/bytes-and-numbers.md` *(2026-09-09)*
- [ ] S4 — Read pt 2 (syscalls & files) §5.1 Users onward → bullets in `notes/os/eblog-syscalls.md` → ritual
- [ ] S5 — Read pt 2, chapters 4–6 → append bullets → ritual
- [ ] S6 — Read pt 2, chapters 7–9 → append bullets → ritual
- [ ] S7 — Build: `experiments/syscalls/` — Go program doing raw open/read/write; watch its syscalls (`dtruss`/`fs_usage` on macOS)

### M1 · Virtualization refresh (OSTEP ch 4–7)
- [ ] S8 — Read ch 4 (Processes) + ch 5 (Process API) → notes → ritual
- [ ] S9 — Read ch 6 (Limited Direct Execution), skim ch 7 (Scheduling) → notes → ritual
- [ ] S10 — Build: mini shell in Go — run commands, wait, exit status
- [ ] S11 — Build: mini shell — pipes (`ls | wc -l`) and redirects

### M2 · Memory (OSTEP ch 13–15, 18–19, 21–22)
- [ ] S12 — Read ch 13 (Address Spaces) + ch 14 (Memory API) → ritual
- [ ] S13 — Read ch 15 (Address Translation) → ritual
- [ ] S14 — Read ch 18 (Paging) → ritual
- [ ] S15 — Read ch 19 (TLBs) → ritual
- [ ] S16 — Read ch 21–22 (Swapping: mechanisms + policies) → ritual
- [ ] S17 — Build: `experiments/memory/` — pointer addresses across runs (ASLR), `syscall.Mmap` demo; start `gokv` v0: in-memory map + REPL (get/set/del)

### M3 · Concurrency (OSTEP ch 26–30, 32)
- [ ] S18 — Read ch 26 (Concurrency intro) → ritual
- [ ] S19 — Read ch 27 (Thread API) + map to goroutines → ritual
- [ ] S20 — Read ch 28 (Locks) → ritual
- [ ] S21 — Read ch 29 (Locked Data Structures) → ritual
- [ ] S22 — Read ch 30 (Condition Variables) → ritual
- [ ] S23 — Read ch 32 (Concurrency Bugs) → ritual
- [ ] S24 — Build: `gokv` v1 — TCP server, goroutine per connection, deliberately racy
- [ ] S25 — Build: `gokv` v1 — catch the race with `-race`, fix with RWMutex; benchmark Mutex vs RWMutex vs sharded map

### M4 · Persistence (OSTEP ch 36–37, 39–40, 42–44)
- [ ] S26 — Read ch 36 (I/O Devices) + ch 37 (Hard Disks) → ritual
- [ ] S27 — Read ch 39 (Files & Directories) → ritual
- [ ] S28 — Read ch 40 (FS Implementation) → ritual
- [ ] S29 — Read ch 42 (Crash Consistency + Journaling) → ritual ★ core DB concept
- [ ] S30 — Read ch 43 (Log-structured FS) → ritual ★ this is the bridge to Bitcask/LSM
- [ ] S31 — Read ch 44 (Flash/SSDs) → ritual
- [ ] S32 — Build: `gokv` v2 — append-only log file, replay on startup
- [ ] S33 — Build: `gokv` v2 — fsync experiments (throughput with/without), kill -9 durability test
- [ ] S34 — **Phase 1 exit demo:** concurrent crash-safe KV server; review all Phase 1 diagrams with Claude

---

## Phase 2 — Database internals (checklist expanded when Phase 1 nears done)

- **M5** · Bitcask storage engine (~4 wks) — bitcask-intro.pdf, DDIA 2e (2026) storage & retrieval chapter *[remap to 2e ch# when expanded]*, 15-445 lectures 1–4; build keydir + segments + compaction
- **M6** · Pages, buffer pools, B+trees (~7 wks, hardest stretch) — 15-445 + Database Internals ch 2–6 + build-your-own.org; disk B+tree → `godb`
- **M7** · Transactions, WAL, MVCC (~4 wks) — 15-445 + DDIA 2e (2026) transactions chapter *[remap to 2e ch# when expanded]*
- **M8** · Tiny SQL layer (optional, ~3 wks)

## Phase 3 — Distributed systems (checklist expanded when Phase 2 nears done)

- **M9** · Fundamentals — Kleppmann lectures + DDIA 2e (2026) replication & partitioning chapters; Fly.io "Gossip Glomers" dist-sys challenges 1–3 (echo → unique IDs → broadcast/gossip), in Go on the Maelstrom harness (injects partitions & message loss). https://fly.io/dist-sys/
  - 📖 *The lifecycle of a sharded Postgres query* (PlanetScale/Neki), https://planetscale.com/blog/the-lifecycle-of-a-sharded-postgres-query — ~8k words tracing ONE query end-to-end through a real sharded DB: proxy routing, wire protocol, distributed planning, hash-join/nested-loop across shards, scatter-gather + aggregation, connection pooling. Concrete companion to the DDIA partitioning theory; revisit at M12 to compare with MongoDB sharding.
- **M10** · Raft — paper + MIT 6.5840 (formerly 6.824) lectures + labs; Fly.io Gossip Glomers challenges 4–6 (grow-only counter → Kafka-style replicated log → totally-available txns — #6 is basically a mini distributed DB, on-ramp to M11); own Raft in Go. Use it *selectively* (not the whole ~20-lecture course): Raft lectures + Lab 2 (Raft, essential); Lab 3 (KV on Raft) → feeds M11; Lab 4 (sharded KV, optional) → bridges M12. https://pdos.csail.mit.edu/6.824/
- **M11** · Capstone — Raft-replicated `godb` (toydb architecture as reference)

## Phase 4 — Real-world case study: how MongoDB works *(added 2026-09-11, Prabhu's request)*

The synthesis finale — no new building. Apply everything learned to a real production
database and see how the primitives you built by hand assemble into a shipping system.
Comes last, after the KV/`godb` work is done, because it only makes sense once you've built
the pieces it's made of.

- **M12** · How MongoDB works — map its internals onto what you built:
  - **Storage engine (WiredTiger):** BSON document model, B+tree on-disk storage, MVCC +
    snapshots, journaling (WAL), checkpoints, compression. → builds on **M6** (B+tree/buffer
    pool) + **M7** (WAL/MVCC).
  - **Indexing:** B+tree indexes; compound, multikey (arrays), partial/sparse/TTL; the query
    planner and index selection. → builds on **M6**.
  - **Sharding & replication:** shard keys, chunks, balancer, config servers, `mongos`
    routing; replica sets with a Raft-like election + the oplog for replication; read/write
    concerns and consistency tradeoffs. → builds on **M9** (partitioning) + **M10** (Raft).
  - Resources: MongoDB Architecture Guide + internals docs, WiredTiger docs, Aphyr/Jepsen
    MongoDB analyses (consistency & sharding failure modes), DDIA cross-refs; re-read the
    PlanetScale sharded-Postgres-query article (M9) and contrast its sharding/routing with
    MongoDB's `mongos` + shard keys.
  - Output: a diagram set mapping *MongoDB component → the gokv/godb primitive it corresponds
    to*, plus an explain-back session with Claude.

---

## Session log

| # | Date | What happened | Next step |
|---|------|---------------|-----------|
| S1 | 2026-07-15 | Scaffolded repo, plan approved, gokv module created. Resized M0 twice: settled on ~3 blog chapters per session (start small) | S2: read eblog pt 1, chapters 1–3 |
| S2 | 2026-07-16 | Read eblog pt 1 ch 1–3, full ritual: summary gap-checked (OR not AND on the 3 traits; low-level = bytes/registers; mindset branch), 3/3 probing Qs answered, mind map saved to diagrams/ | S3: read eblog pt 1, chapters 4–5 (data & code segments) |
| S3 | 2026-09-09 | Read §4.5–4.7 + data/code segment investigations + §5.1 `shexdump`. Heavy Q&A + built `4.6torsolines.go` (torso→lines): mode-via-sentinel, ReadFull short reads, the `ReadBytes`/EOF gotcha (3 tries). Reasoned through backward-lines (Strategy A) but didn't implement. Big side-threads captured: ReadAll/mmap/buffer-pool, string blob & `__rodata`, bit shifts/serialization/endianness (→ new `bytes-and-numbers.md`), hex-vs-decimal, linear-scan-vs-index = DB index motivation. Platform gotcha: macOS code-signing SIGKILL on patched binaries (`codesign -f -s -`) | S4: read eblog **pt 2** (startingsystems2.html) — §5.1 Users onward (syscalls & files) |
