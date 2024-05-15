package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/tools/cover"
)

func command(stdin io.Reader, stdout io.Writer) error {
	var (
		fileName string
		err      error
	)
	switch flag.CommandLine.Arg(0) {
	case "test":
		fileName, err = runGoTests()
		if err != nil {
			return err
		}
		defer deleteTempFile(fileName)
	case "-":
	// user explicitly specified to use stdin
	default:
		fileName = flag.CommandLine.Arg(0)
	}

	// parse coverage profiles from stdin or file path if provided
	profiles, err := parseCoverProfiles(fileName, stdin)
	if err != nil {
		return err
	}

	// find path that points to module root which will be needed to construct an
	// absolute file path to the Go source files to generate a diff for
	moduleInfo, err := findModuleInfo()
	if err != nil {
		return err
	}
	return printDiff(stdout, profiles, moduleInfo)
}

func runGoTests() (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "coverdiff-*")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("go", "test", "./...", "-cover", "-coverprofile="+f.Name())
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func deleteTempFile(fileName string) {
	if err := os.Remove(fileName); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("coverdiff: %w", err))
		os.Exit(1)
	}
}

func parseCoverProfiles(fileName string, stdin io.Reader) ([]*cover.Profile, error) {
	if fileName != "" {
		profiles, err := cover.ParseProfiles(fileName)
		return profiles, err
	}
	profiles, err := cover.ParseProfilesFromReader(stdin)
	return profiles, err
}
