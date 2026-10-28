package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Lookup("root")
	if err != nil {
		fmt.Fprintf(os.Stderr, "user read error %v", err)
		os.Exit(1)
	}

	current, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "current user read error %v", err)
		os.Exit(1)
	}

	fmt.Printf("%v %v %v\n", usr.Name, usr.Gid, usr.Uid)
	fmt.Printf("%v %v %v\n", current.Name, current.Gid, current.Uid)
}
