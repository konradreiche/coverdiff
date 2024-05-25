# 📑 coverdiff
[![ci](https://github.com/konradreiche/coverdiff/actions/workflows/ci.yaml/badge.svg)](https://github.com/konradreiche/coverdiff/actions) [![codecov](https://codecov.io/gh/konradreiche/coverdiff/graph/badge.svg?token=kXoAXWhLJS)](https://codecov.io/gh/konradreiche/coverdiff)

Print your Go test coverage line-by-line in the form of a code diff, highlighting each line. This tool takes inspiration from Go's HTML presentation of test coverage but brings it to the terminal instead.

[![asciicast](https://asciinema.org/a/627967.svg)](https://asciinema.org/a/627967)

## Getting Started

```
go install github.com/konradreiche/coverdiff@latest
```

## Usage

```
Usage:
	coverdiff test [packages]
	coverdiff [file]

Flags:
	-h, --help	print help text

Examples:

	coverdiff test ./...

	PAGER=delta coverdiff test

	go test -cover -coverprofile=coverage.out
	cat coverage.out | coverdiff

	go test -cover -coverprofile >(coverdiff)
```
