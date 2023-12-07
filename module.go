package main

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type moduleInfo struct {
	modulePath string
	moduleDir  string
}

func findModuleInfo() (moduleInfo, error) {
	wd, err := os.Getwd()
	if err != nil {
		return moduleInfo{}, fmt.Errorf("os.Getwd: %w", err)
	}
	moduleDir, err := findModuleDir(wd)
	if err != nil {
		return moduleInfo{}, fmt.Errorf("findModuleDir: %w", err)
	}

	b, err := os.ReadFile(filepath.Join(moduleDir, "go.mod"))
	if err != nil {
		return moduleInfo{}, fmt.Errorf("os.ReadFile: %w", err)
	}
	modulePath := modfile.ModulePath(b)
	if moduleDir == "" {
		return moduleInfo{}, errors.New("no module path found")
	}
	return moduleInfo{
		modulePath: modulePath,
		moduleDir:  moduleDir,
	}, nil
}

func findModuleDir(dir string) (string, error) {
	if dir == "" {
		return "", errors.New("dir not set")
	}
	dir = filepath.Clean(dir)

	// look for enclosing go.mod
	for {
		f := filepath.Join(dir, "go.mod")
		if fi, err := os.Stat(f); err == nil && !fi.IsDir() {
			return dir, nil
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		if d == build.Default.GOROOT {
			// As a special case, don't cross GOROOT to find a go.work file.
			// The standard library and commands built in go always use the vendored
			// dependencies, so avoid using a most likely irrelevant go.work file.
			return "", errors.New("no go.mod file found")
		}
		dir = d
	}
	return "", errors.New("no go.mod file found")
}
