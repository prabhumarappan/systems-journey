package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// fatalf writes a formatted string to standard error and exits with status 1. LIBRARY
func fatalf(format string, args ...interface{}) {
	buf := []byte(fmt.Sprintf(format, args...))
	syscall.Syscall(syscall.SYS_WRITE, uintptr(syscall.Stderr), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))) // no point in checking the error here; we're about to exit.
	syscall.Syscall(syscall.SYS_EXIT, 1, 0, 0)                                                                       // exit with status 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage : %s <file>", os.Args[0])
		os.Exit(1)
	}

	//1. open the file with SYS_OPEN
	filename := []byte(os.Args[1])
	const MODE = syscall.O_RDONLY
	const FLAGS = 0

	fileDescriptor, _, err := syscall.Syscall(syscall.SYS_OPEN, uintptr(unsafe.Pointer(&filename[0])), MODE, FLAGS)
	if err != 0 {
		fatalf("open :%v\n", err)
	}

	//2. read the chunks into memroy with SYS_READ
	var buf [1024]byte

READ:
	for {
		ptr := &buf[0]
		n, _, readErr := syscall.Syscall(syscall.SYS_READ, fileDescriptor, uintptr(unsafe.Pointer(ptr)), uintptr(len(buf)))
		if readErr != 0 {
			fatalf("read: %v\n", readErr)
		}

		for offset := uintptr(0); offset < n; {
			ptr := &buf[offset]
			m, _, writeErr := syscall.Syscall(syscall.SYS_WRITE, uintptr(syscall.Stdout), uintptr(unsafe.Pointer(ptr)), uintptr(n-offset))
			if writeErr != 0 {
				fatalf("write : %v\n", writeErr)
			}

			if m == n {
				continue READ
			}
			offset += m
		}
		if n == 0 {
			break READ
		}
	}

	_, _, _ = syscall.Syscall(syscall.SYS_FSYNC, uintptr(syscall.Stdout), 0, 0)

	syscall.Syscall(syscall.SYS_CLOSE, fileDescriptor, 0, 0)

	syscall.Syscall(syscall.SYS_EXIT, 0, 0, 0)

}
