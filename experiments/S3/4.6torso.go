package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var offset, before, after int
	var from string
	var newline bool

	{ // 1. define and parse command line flags
		flag.IntVar(&offset, "offset", -1, "offset to read from: must be specified")                // -offset is required
		flag.IntVar(&before, "before", 128, "bytes to read before offset: will be clamped to 0")    // defaults to -before 128
		flag.IntVar(&after, "after", 128, "bytes to read after offset: will be clamped to 0")       // defaults to -after 128
		flag.StringVar(&from, "from", "", "file to read from: if empty, reads from standard input") // -from is required
		flag.BoolVar(&newline, "newline", false, "append a newline to the output")                  // -newline is optional
		flag.Parse()
	}

	// bounds checking and normalization
	{
		before = max(before, 0)      // can't be negative
		before = min(before, offset) // can't go past the beginning
		after = max(after, 0)
		if offset < 0 {
			fmt.Fprintf(os.Stderr, "missing or invalid -offset\n")
			os.Exit(1)
		}
	}

	// making sure the required lenght is not zero
	start := offset - before
	n := before + after
	if n == 0 {
		return
	}

	// make buffer to fit the bytes in
	buf := make([]byte, n)

	// read file
	f, err := os.Open(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %s: %v \n", from, err)
		os.Exit(1)
	}

	// seek to the specified start byte
	_, err = f.Seek(int64(start), io.SeekStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seek: %s: %v\n", from, err)
		f.Close()
		os.Exit(1)
	}

	// read from start till the buffer is filled up
	n, err = io.ReadFull(f, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		fmt.Fprintf(os.Stderr, "read: %s: %v\n", from, err)
		os.Exit(1)
	}
	buf = buf[:n]

	// write from buffer to stdout
	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		f.Close()
		os.Exit(1)
	}

	if newline {
		fmt.Println()
	}

	f.Close()
}
