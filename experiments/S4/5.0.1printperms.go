package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	mode := info.Mode()
	fmt.Printf("Symbolic File Permissions: %s\n", mode.Perm())
}
