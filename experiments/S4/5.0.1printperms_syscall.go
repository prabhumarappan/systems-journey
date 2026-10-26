package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	filePath := os.Args[1]
	var st syscall.Stat_t
	syscall.Stat(filePath, &st)
	mode := st.Mode
	rwr := "rwxrwxrwx" // will be used to write symbloic representation
	var out [9]byte
	// for each bit in the permission bit
	// starts from higher order bits
	for i := 0; i < 9; i++ {
		// check if 8-ith position bit is set
		if mode&(1<<(8-i)) != 0 {
			// if it is set, that means that the permission is set, consider whatever value rwr has for that position
			out[i] = rwr[i]
		} else {
			// if not set, that means its empty/unset
			out[i] = '-'
		}
	}
	fmt.Printf("%s\n", out)
}
