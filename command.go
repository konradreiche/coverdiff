package main

import (
	"io"
	"os"

	"golang.org/x/tools/cover"
)

func command(stdin io.Reader, stdout io.Writer) error {
	// parse coverage profiles from stdin or file path if provided
	profiles, err := parseCoverProfiles(stdin)
	if err != nil {
		return err
	}

	// find path that points to module root which will be needed to construct an
	// absolute file path to the Go source files to generate a diff for
	moduleInfo, err := findModuleInfo()
	if err != nil {
		return err
	}
	return printDiff(stdout, profiles, moduleInfo.modulePath)
}

func parseCoverProfiles(stdin io.Reader) ([]*cover.Profile, error) {
	if len(os.Args) > 1 && os.Args[1] != "-" {
		profiles, err := cover.ParseProfiles(os.Args[1])
		return profiles, err
	}
	profiles, err := cover.ParseProfilesFromReader(stdin)
	return profiles, err
}
