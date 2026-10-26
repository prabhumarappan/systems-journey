package main

import (
	"syscall"
	"unsafe"
)

func gowrite(fd int, buf []byte) (int, error) {
	n, _, errno := syscall.Syscall(syscall.SYS_WRITE, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(n), errno
}

func main() {

}
