package main

import (
	"bytes"
	"path/filepath"
	"runtime"
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
	moduleInfo, err := findModuleInfo()
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := printDiff(&b, profiles, moduleInfo); err != nil {
		t.Fatal(err)
	}
}
