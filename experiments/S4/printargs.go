package main

import (
	"fmt"
	"os"
)

func main() {
	for i, arg := range os.Args {
		fmt.Fprintf(os.Stdout, "%d: %s\n", i, arg)
	}
}
