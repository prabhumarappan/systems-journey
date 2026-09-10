// want to see if i can do the backward thing with lines too
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func fetchNLines(lines int, f *os.File) []byte {
	reader := bufio.NewReader(f)
	var out []byte

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "read: %v\n", err)
			os.Exit(1)
		}
		if err == io.EOF {
			lines = 0
		}
		out = append(out, line...)
		lines = lines - 1
		if lines <= 0 {
			break
		}

	}

	return out
}

func fetchNBytes(before int, after int, f *os.File) []byte {
	n := before + after
	if n == 0 {
		return nil
	}
	// // make buffer to fit the bytes in
	buf := make([]byte, n)

	// read from start till the buffer is filled up
	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}
	buf = buf[:n]

	return buf
}

func readFromFile(before int, after int, lines int, readbybytes bool, f *os.File) []byte {
	if readbybytes {
		return fetchNBytes(before, after, f)
	} else {
		return fetchNLines(lines, f)
	}
}

func main() {
	var offset, before, after, lines, start int
	var from string
	var newline bool
	readbybytes := true

	{ // 1. define and parse command line flags
		flag.IntVar(&offset, "offset", -1, "offset to read from: must be specified")                // -offset is required
		flag.IntVar(&before, "before", -1, "bytes to read before offset: will be clamped to 0")     // defaults to -before -1
		flag.IntVar(&after, "after", -1, "bytes to read after offset: will be clamped to 0")        // defaults to -after -1
		flag.StringVar(&from, "from", "", "file to read from: if empty, reads from standard input") // -from is required
		flag.IntVar(&lines, "lines", -1, "number of lines to read from the offset")
		flag.BoolVar(&newline, "newline", false, "append a newline to the output") // -newline is optional
		flag.Parse()
	}

	// figure out which type of reading to use
	{
		if lines > 0 {
			readbybytes = false
		}
		if offset < 0 {
			fmt.Fprintf(os.Stderr, "missing or invalid -offset\n")
			os.Exit(1)
		}
	}

	// bounds checking and normalization
	if readbybytes {
		before = max(before, 0)      // can't be negative
		before = min(before, offset) // can't go past the beginning
		after = max(after, 0)
		start = offset - before
	} else {
		start = offset
	}

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

	//
	buf := readFromFile(before, after, lines, readbybytes, f)

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
