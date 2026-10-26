package main

import (
	"fmt"
	"os"
)

func main() {
	env := os.Environ()
	var printed int
	keys := os.Args[1:]

	for _, key := range keys {
		val, ok := lookupenv(env, key)
		if ok {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", key, val)
			printed++
		}
	}
	if printed == len(keys) {
		os.Exit(0) //
	}

	fmt.Fprintf(os.Stderr, "\nmissing %d/%d environment variables\n", len(os.Args)-1-printed, len(os.Args)-1)
	for _, key := range keys {
		if _, ok := lookupenv(env, key); !ok {
			fmt.Fprintf(os.Stderr, "%s\n", key)
		}
	}
	os.Exit(1)
}

func lookupenv(env []string, key string) (string, bool) {
	for i := range len(env) {
		environment_variable := env[i]
		// check if env variable is atleast as much length as the passed in key + one more for the = symbol
		if len(environment_variable) < len(key)+1 {
			continue
		}
		// check if the first n characters of the env variable are equal to the key
		if environment_variable[:len(key)] != key {
			continue
		}
		// making sure that there is a = right after the key in the env variable
		if environment_variable[len(key)] != '=' {
			continue
		}
		return environment_variable[len(key)+1:], true
	}
	return "", false
}
