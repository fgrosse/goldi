package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GoPathChecker struct {
	Verbose bool
}

func NewGoPathChecker(isVerbose bool) *GoPathChecker {
	return &GoPathChecker{isVerbose}
}

func (c *GoPathChecker) PackageName(outputPath string) string {
	c.log("GoPathChecker is determining package name for output path %q", outputPath)

	goPaths := os.Getenv("GOPATH")
	if outputPath == "" || goPaths == "" {
		c.log("output path or GOPATH is empty")
		return ""
	}

	outputPath, err := filepath.Abs(outputPath)
	if err != nil {
		// this can only happen if go has trouble determining the current working directory on this OS
		panic(fmt.Errorf("Could not get absolut file path from %q: %s", outputPath, err))
	}
	c.log("absolute output path is %q", outputPath)

	outputDir := filepath.Dir(outputPath)
	c.log("output dir is %q", outputDir)
	for _, goPath := range strings.Split(goPaths, ":") {
		goPath = goPath + "/src/"
		if strings.Contains(outputDir, goPath) {
			c.log("Found %q in GOPATH %q", outputDir, goPath)
			packageName := strings.TrimPrefix(outputDir, goPath)
			return packageName
		}

		c.log("%q is not contained in GOPATH %q", outputDir, goPath)
	}

	return ""
}

func (c *GoPathChecker) log(message string, args ...interface{}) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, message+"\n", args...)
	}
}
