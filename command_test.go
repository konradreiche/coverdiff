package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestCommand(t *testing.T) {
	t.Run("from-stdin-pipe", func(t *testing.T) {
		args := os.Args
		t.Cleanup(func() { os.Args = args })

		// override os.Args which contains flags from the test binary
		os.Args = []string{
			"coverdiff",
		}
		flag.Parse()

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString(readFile(t, "testdata/coverage.out"))
		if err := command(stdin, &stdout); err != nil {
			t.Fatal(err)
		}

		got := stdout.String()
		want := readFile(t, "testdata/coverdiff.out")
		if got != want {
			t.Errorf("got len=%d, want len=%d", len(got), len(want))
		}
	})

	t.Run("from-stdin-dash", func(t *testing.T) {
		os.Args[1] = "-"
		flag.Parse()

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString(readFile(t, "testdata/coverage.out"))
		if err := command(stdin, &stdout); err != nil {
			t.Fatal(err)
		}

		got := stdout.String()
		want := readFile(t, "testdata/coverdiff.out")
		if got != want {
			t.Errorf("got len=%d, want len=%d", len(got), len(want))
		}
	})

	t.Run("from-file", func(t *testing.T) {
		os.Args[1] = filepath.Join(projectPath, "testdata/coverage.out")
		flag.Parse()

		var stdout bytes.Buffer
		if err := command(nil, &stdout); err != nil {
			t.Fatal(err)
		}

		got := stdout.String()
		want := readFile(t, "testdata/coverdiff.out")
		if got != want {
			t.Errorf("got len=%d, want len=%d", len(got), len(want))
		}
	})

	t.Run("from-file-in-subdirectory", func(t *testing.T) {
		if err := os.Chdir("testdata"); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.Chdir(projectPath); err != nil {
				t.Fatal(err)
			}
		})
		os.Args[1] = filepath.Join(projectPath, "testdata/coverage.out")
		flag.Parse()

		var stdout bytes.Buffer
		if err := command(nil, &stdout); err != nil {
			t.Fatal(err)
		}

		got := stdout.String()
		want := readFile(t, "coverdiff.out")
		if got != want {
			t.Errorf("got len=%d, want len=%d", len(got), len(want))
		}
	})

	t.Run("outside-go-module", func(t *testing.T) {
		changeDir(t, "..")

		var stdout bytes.Buffer
		err := command(nil, &stdout)

		got := err.Error()
		want := "findModuleDir: no go.mod file found"
		if got != want {
			t.Errorf("got %s, want: %s", got, want)
		}
	})

	t.Run("run-go-tests", func(t *testing.T) {
		os.Args[1] = "test"
		flag.Parse()

		changeDir(t, "testdata")

		var stdout bytes.Buffer
		if err := command(nil, &stdout); err != nil {
			t.Fatal(err)
		}

		got := stdout.String()
		want := readFile(t, "coverdiff.out")
		if got != want {
			t.Errorf("got len=%d, want len=%d", len(got), len(want))
		}
	})
}

func readFile(tb testing.TB, name string) string {
	tb.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		tb.Fatal(err)
	}
	return string(b)
}
