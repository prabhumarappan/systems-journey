package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

var table = [8]string{"---", "--x", "-w-", "-wx", "r--", "r-x", "rw-", "rwx"}

func symbolicToOctet(symbolic string) (string, error) {
	if len(symbolic) != 9 {
		return "", fmt.Errorf("invalid symbolic passed in %q", symbolic)
	}
	var octet [3]int
	for i := 0; i < 3; i++ {
		start := 3 * i
		end := 3 * (i + 1)
		sub := symbolic[start:end]
		for j := 0; j < len(table); j++ {
			table_item := table[j]
			if table_item == sub {
				octet[i] = j
				break
			}
		}
	}
	return fmt.Sprintf("%d%d%d", octet[0], octet[1], octet[2]), nil
}

func parseOctetToDecimal(octet string) (int, error) {
	result := 0
	for i := 0; i < len(octet); i++ {
		current := int(octet[i] - '0')
		if current > 7 {
			return 0, fmt.Errorf("invalid octet passed in for parsing %q", octet)
		}
		result = result<<3 | current
	}
	return result, nil
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage %s <file> <permission>", os.Args[0])
		os.Exit(1)
	}
	filePath := append([]byte(os.Args[1]), 0)
	permission := os.Args[2]

	var octet string

	if !isInteger(permission) {
		var err error
		octet, err = symbolicToOctet(permission)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error raised while converting symbolic to otctet %s", err)
			os.Exit(1)
		}
	} else {
		octet = permission
	}

	decimalPermission, err := parseOctetToDecimal(octet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error raised while converting octet to decimal %s", err)
		os.Exit(1)
	}

	// now do the chmod using syscall
	_, _, errno := syscall.Syscall(syscall.SYS_CHMOD, uintptr(unsafe.Pointer(&filePath[0])), uintptr(decimalPermission), 0)
	if errno != 0 {
		fmt.Fprintf(os.Stderr, "error with syscall chmod %v", errno)
		os.Exit(1)
	}
}
