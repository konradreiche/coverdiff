package main

import (
	"fmt"
	"io"
)

const usage = `📑 coverdiff - print Go test coverage as diff

Usage:
	coverdiff [file]

Flags:
	-h, --help	print help text

coverdiff is a tool designed to process the cover profile output of Go and
display the coverage as a diff, similar to Go's -html option but optimized for
terminal convenience. You can provide the cover profile as a file or pass it
through standard input.

Examples:

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
