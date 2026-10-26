package main

import (
	"fmt"
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
	hello := []byte("hello world\n")

	ptr := unsafe.Pointer(&hello[0])

	n, _, errno := syscall.Syscall(syscall.SYS_WRITE, uintptr(syscall.Stdout), uintptr(ptr), uintptr(len(hello)))
	if errno != 0 {
		fatalf("write: %v\n", errno)
	}
	if int(n) != len(hello) {
		fatalf("write: wrote %d bytes, expected %d\n", n, len(hello))
	}
}
