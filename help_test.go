package main

import (
	"bytes"
	"testing"
)

func TestPrintUsage(t *testing.T) {
	var stderr bytes.Buffer
	printUsage(&stderr)()

	got := stderr.String()
	want := usage
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
