package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.CommandLine.Usage = printUsage(os.Stderr)
	flag.Parse()

	if err := command(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("coverdiff: %w", err))
		fmt.Fprintln(os.Stderr, "Use coverdiff --help to display help text")
		os.Exit(1)
	}
}
