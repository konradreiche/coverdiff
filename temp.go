package main

import (
	"fmt"
	"os"
)

func createTempFile() (*os.File, error) {
	dir := os.TempDir()
	file, err := os.CreateTemp(dir, "coverdiff.test.*.out")
	if err != nil {
		return nil, fmt.Errorf("createTempFile: %w", err)
	}
	return file, nil
}
