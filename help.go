package main

import (
	"fmt"
	"io"
)

const usage = `📑 coverdiff - print Go test coverage as diff

Usage:
	coverdiff test [packages]
	coverdiff [file]

Flags:
	-h, --help	print help text

coverdiff is a tool that processes Go cover profiles and displays coverage
differences directly on the terminal, similar to Go's -html option. You can
provide the cover profile as a file, pass it through standard input, or let
coverdiff run the Go tests.

To handle larger outputs and add syntax highlighting automatically, set the
PAGER environment variable to your preferred pager, such as delta.

Examples:

	coverdiff test ./...

	PAGER=delta coverdiff test

	go test -cover -coverprofile=coverage.out
	cat coverage.out | coverdiff

	go test -cover -coverprofile >(coverdiff)

`

func printUsage(stderr io.Writer) func() {
	// return func() to conform with the required type of flag.CommandLine.Usage
	return func() {
		fmt.Fprint(stderr, usage)
	}
}
