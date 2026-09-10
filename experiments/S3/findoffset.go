package main

import (
	"fmt"
	"os"
)


func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: findoffset <filename> <string>")
		os.Exit(1)
	}

	filepath, pattern := os.Args[1], os.Args[2]

	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s : %v", filepath, err)
		os.Exit(1)
	}

	for i:=0; i<len(b) - len(pattern); i++ {
		for j := range pattern {
			if b[i+j] != pattern[j] {
				break
			}

			if j == len(pattern) - 1 {
				fmt.Fprintf(os.Stdout, "%d\n", i)
				os.Exit(0)
			}
		}
	}
	os.Exit(1)
}
