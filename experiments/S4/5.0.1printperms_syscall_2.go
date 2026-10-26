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

	// based on octet position 0...7, figure out what the value will be
	// follows read write execute
	// 0 => ---
	// 1 => 001 => --x
	// 2 => 010 => -w-
	// 3 => 011 => -wx
	var table = [8]string{"---", "--x", "-w-", "-wx", "r--", "r-x", "rw-", "rwx"}
	// output string
	var out string
	// do 3 groups
	// user group other

	// find out group value from table using bit manipulation
	/*
		 	mode = 33188
			find value for first group
				33188 >> 6 => 1000000110 = 518

				518 & 7 =
					1000000110
					0000000111
			   	   &0000000110
			       =       (6)

			find value for second group
				33188 >> 3 => 1000000110100 = 4148

				4148 & 7 =
					1000000110100
					0000000000111
			        &0000000000100
			        =          (4)

			find value for third group
				33188 >> 0 => 1000000110100100 = 33188

				33188 & 7 =
					1000000110100100
					0000000000000111
			        &0000000000000100
			        =             (4)
	*/
	for i := 0; i < 3; i++ {
		pos := (mode >> (9 - 3*(i+1))) & 7
		out += table[pos]
	}
	fmt.Printf("%s\n", out)
}
