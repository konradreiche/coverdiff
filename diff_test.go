package main

import (
	"bytes"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"testing"

	"golang.org/x/tools/cover"
)

var (
	_, b, _, _  = runtime.Caller(0)
	projectPath = filepath.Dir(b)
)

func TestPrintDiff(t *testing.T) {
	profiles, err := cover.ParseProfiles("testdata/coverage.out")
	if err != nil {
		t.Fatal(err)
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		t.Fatal("binary build information not available")
	}

	var b bytes.Buffer
	if err := printDiff(&b, profiles, bi.Path); err != nil {
		t.Fatal(err)
	}
}
