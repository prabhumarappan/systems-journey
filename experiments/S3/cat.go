package main

import (
	"os"
	"fmt"
	"io"
)

func main() {
	for _, file := range os.Args[1:] {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v", file, err)
			os.Exit(1)
		}

		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v", file, err)
			os.Exit(1)
		}

		os.Stdout.Write(b)
	}
}
