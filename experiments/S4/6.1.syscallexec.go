package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	// check for arg size to make sure it includes another executable
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: %v <filepath>\n", os.Args[0])
		os.Exit(1)
	}

	// keep everything except the first, which is the path of the script itself
	goargs := os.Args[1:]

	// check if executable path is absolute
	if len(goargs) == 0 || goargs[0][0] != '/' {
		fmt.Fprintf(os.Stderr, "error: %s filepath needs to start with /", goargs[0])
		os.Exit(1)
	}

	err := exec(goargs, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with exec: %v", err)
	}
}

func pointerFromStr(s string) unsafe.Pointer {
	// create a byte array of size 1 more than the string itself for null termination
	newS := make([]byte, len(s)+1)
	copy(newS, s)
	return unsafe.Pointer(&newS[0])
}

func exec(args []string, env []string) error {
	// make a list of pointers with 1 extra for null pointer
	cargs := make([]unsafe.Pointer, len(args)+1)
	for i := range args {
		cargs[i] = pointerFromStr(args[i])
	}

	// same as above
	cenv := make([]unsafe.Pointer, len(env)+1)
	for i := range env {
		cenv[i] = pointerFromStr(env[i])
	}

	// convert the path of the executable to be pointer
	path := pointerFromStr(args[0])

	_, _, err := syscall.Syscall(syscall.SYS_EXECVE, uintptr(path), uintptr(unsafe.Pointer(&cargs[0])), uintptr(unsafe.Pointer(&cenv[0])))
	return err
}
