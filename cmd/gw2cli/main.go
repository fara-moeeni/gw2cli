package main

import (
	"fmt"
	"os"
)

const Version = "2.8.2"

func main() {
	if err := run(os.Args[1:], os.Getenv("GW2_API_KEY")); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
