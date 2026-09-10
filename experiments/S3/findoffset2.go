package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage findoffset2 <filename> <string> <occurence>")
	}
	filepath, pattern, s_occurences := os.Args[1], os.Args[2], os.Args[3]

	occurences, err := strconv.Atoi(s_occurences)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting %s to int", s_occurences)
		os.Exit(1)
	}

	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s: %v", filepath, err)
		os.Exit(1)
	}

	occured := 0
	for i := 0; i < len(b)-len(pattern); i++ {
		for j := range pattern {
			if b[i+j] != pattern[j] {
				break
			}

			if j == len(pattern)-1 {
				occured += 1
				if occured == occurences {
					fmt.Fprintf(os.Stdout, "%d\n", i)
					os.Exit(0)
				}
			}
		}
	}
	os.Exit(1)
}
