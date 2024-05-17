package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/cover"
)

func command(stdin io.Reader, stdout io.Writer) error {
	var fileName string
	switch flag.CommandLine.Arg(0) {
	case "test":
		output, err := runGoTests()
		if err != nil {
			return err
		}
		stdin = output
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

	if pager := os.ExpandEnv("$PAGER"); pager != "" {
		return usePager(pager, stdout, func(w io.Writer) error {
			return printDiff(w, profiles, moduleInfo)
		})
	}
	return printDiff(stdout, profiles, moduleInfo)
}

func usePager(
	name string,
	stdout io.Writer,
	printDiff func(io.Writer) error,
) error {
	cmd := exec.Command(name)
	cmd.Stdout = stdout
	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer w.Close()
		printDiff(w)
	}()

	return cmd.Run()
}

func runGoTests() (*bytes.Buffer, error) {
	cmd := exec.Command(
		"go",
		"test",
		"./...",
		"-json",
		"-cover",
		"-coverprofile=/dev/stdout",
	)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(string(b), "\n") {
		// filter Go JSON test output
		if strings.HasPrefix(line, `{"Time"`) {
			continue
		}
		lines = append(lines, line)
	}
	coverProfile := strings.Join(lines, "\n")
	return bytes.NewBufferString(coverProfile), nil
}

func parseCoverProfiles(fileName string, stdin io.Reader) ([]*cover.Profile, error) {
	if fileName != "" {
		profiles, err := cover.ParseProfiles(fileName)
		return profiles, err
	}
	profiles, err := cover.ParseProfilesFromReader(stdin)
	return profiles, err
}
