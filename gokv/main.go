// gokv — a key-value store that grows with the learning plan:
// Phase 1: in-memory map → concurrent TCP server → append-only log persistence
// Phase 2: Bitcask → disk B+tree → WAL + transactions (becomes godb)
// Phase 3: Raft-replicated distributed store
package main

import "fmt"

func main() {
	fmt.Println("gokv v0 — nothing here yet; first real code lands in session S14")
}
