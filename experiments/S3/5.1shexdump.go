package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	var src io.Reader
	switch len(os.Args) {
	case 1:
		src = os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v", os.Args[1], err)
			os.Exit(1)
		}
		defer f.Close()
		src = f
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s [filename]", os.Args[0])
		os.Exit(1)
	}

	err := hexdump(os.Stdout, src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hexdump: %v", err)
		os.Exit(1)
	}
}

func hexdump(dst io.Writer, src io.Reader) error {
	r := bufio.NewReader(src) //make a buffered reader and read batched io calls

	for {
		var raw [16]byte                 // make a 16 byte array to read into
		encoded := make([]byte, 0, 16*3) // 16 hex per line, each like "ab ", so 16 * 3
		n, err := io.ReadFull(r, raw[:]) // read 16 bytes into raw

		const hex = "0123456789abcdef" // used for hex rep
		if n != 0 {
			// for the 16 characters, generate the hex and the space
			for i := range min(n, 16) {
				encoded = append(encoded, hex[raw[i]>>4], hex[raw[i]&0x0f], ' ')
			}

			// add a newline
			encoded[len(encoded)-1] = '\n'

			// write everything to destination, std.out in this case
			if _, err := dst.Write(encoded); err != nil {
				return err
			}
		}

		// reached the EOF/partial byte, return nil, otherwise return the error
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}
