package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/cover"
)

func printDiff(stdout io.Writer, profiles []*cover.Profile, modulePath string) error {
	for _, profile := range profiles {
		// create file path using absolute path to module
		filepath := strings.ReplaceAll(profile.FileName, modulePath+"/", "")

		b, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("os.ReadFile: %w", err)
		}
		lines := strings.Split(string(b), "\n")

		// track coverage mapping source code line to profile block
		blockByLine := make(map[int][]cover.ProfileBlock)

		for _, block := range profile.Blocks {
			for i := block.StartLine; i <= block.EndLine; i++ {
				// handle coverage of blocks by skipping coverage for only one line
				// otherwise we will print duplicate lines
				if len(blockByLine[i]) == 1 {
					continue
				}
				blockByLine[i] = append(blockByLine[i], block)
			}
		}

		// print diff file headers
		fmt.Fprintf(stdout, "diff --git a/%s b/%s\n", filepath, filepath)
		fmt.Fprintf(stdout, "--- a/%s\n", filepath)
		fmt.Fprintf(stdout, "+++ b/%s\n", filepath)

		// print diff index header
		fmt.Fprintf(stdout, "@@ -%d,%d +%d,%d @@ %s\n", 0, 0, len(lines), 0, lines[0])

		// print all lines regardless of coverage
		for i, line := range lines[1:] {
			blocks, ok := blockByLine[i+2]
			if !ok {
				fmt.Fprintln(stdout, line)
				continue
			}
			for _, block := range blocks {
				if block.Count == 1 {
					fmt.Fprintf(stdout, "+%s\n", line)
				} else {
					fmt.Fprintf(stdout, "-%s\n", line)
				}
			}
		}
	}
	return nil
}
