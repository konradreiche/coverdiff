package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/cover"
)

func printDiff(stdout io.Writer, profiles []*cover.Profile, moduleInfo moduleInfo) error {
	for _, profile := range profiles {
		// use absolute path to module to read source code
		absPath := strings.ReplaceAll(profile.FileName, moduleInfo.modulePath, moduleInfo.moduleDir)
		b, err := os.ReadFile(absPath)
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

		// use relative path to source code for diff headers
		relPath := strings.ReplaceAll(profile.FileName, moduleInfo.modulePath+"/", "")

		// print diff file headers
		fmt.Fprintf(stdout, "diff --git a/%s b/%s\n", relPath, relPath)
		fmt.Fprintf(stdout, "--- a/%s\n", relPath)
		fmt.Fprintf(stdout, "+++ b/%s\n", relPath)

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
