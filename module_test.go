package main

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"
)

func TestFindModuleInfo(t *testing.T) {
	moduleInfo, err := findModuleInfo()
	if err != nil {
		t.Fatal(err)
	}
	got := moduleInfo.modulePath
	want := "github.com/konradreiche/coverdiff"
	if got != want {
		t.Errorf("got %s, want: %s", got, want)
	}
}

func TestFindModuleDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := findModuleDir(".")
		if err != nil {
			t.Fatal(err)
		}
		wd := changeDir(t, "testdata")
		_, err = findModuleDir(wd)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty-dir-param", func(t *testing.T) {
		_, err := findModuleDir("")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		got := err.Error()
		want := "dir not set"
		if got != want {
			t.Errorf("got %s, want: %s", got, want)
		}
	})

	t.Run("goroot-boundary", func(t *testing.T) {
		wd := changeDir(t, "..")
		build.Default.GOROOT = filepath.Join(wd, "..")

		_, err := findModuleDir(wd)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		got := err.Error()
		want := "no go.mod file found"
		if got != want {
			t.Errorf("got %s, want: %s", got, want)
		}
	})
}

func changeDir(tb testing.TB, path string) string {
	tb.Helper()
	// reset to current working directory afterwards to avoid breaking
	// subsequent tests
	cwd, err := os.Getwd()
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			tb.Error(err)
		}
	})
	if err := os.Chdir(path); err != nil {
		tb.Fatal(err)
	}
	return filepath.Join(cwd, path)
}
