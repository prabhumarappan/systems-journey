package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// fatalf prints an error message to stderr with fmt.Fprintf, then exits with status 1.
func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file> <offset> <replacement>", os.Args[0])
		os.Exit(1)
	}

	var (
		file        = os.Args[1]
		offset, err = strconv.ParseInt(os.Args[2], 0, 64)
		replacement = os.Args[3]
	)

	if err != nil || offset < 0 {
		fatalf("invalid offset: %v\nUsage: %s <file> <offset> <replacement>", err, os.Args[0])
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		fatalf("open %s: %v\n", file, err)
	}
	defer f.Close()

	_, err = io.CopyN(os.Stdout, f, offset)
	if err != nil {
		fatalf("copy: %v\n", err)
	}

	_, err = os.Stdout.Write([]byte(replacement))
	if err != nil {
		fatalf("write: %v\n", err)
	}

	if _, err := io.CopyN(io.Discard, f, int64(len(replacement))); err != nil {
		fatalf("copy: %v\n", err)
	}

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		fatalf("copy: %v\n", err)
	}

}
