package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GoPathChecker struct{}

func NewGoPathChecker() *GoPathChecker {
	return &GoPathChecker{}
}

func (c *GoPathChecker) PackageName(outputPath string) string {
	goPaths := os.Getenv("GOPATH")
	if outputPath == "" || goPaths == "" {
		return ""
	}

	outputPath, err := filepath.Abs(outputPath)
	if err != nil {
		panic(fmt.Errorf("Could not get absolut file path from %q", outputPath))
	}
	outputDir := filepath.Dir(outputPath)
	for _, goPath := range strings.Split(goPaths, ":") {
		goPath = goPath + "/src/"
		if strings.Contains(outputDir, goPath) {
			packageName := strings.TrimPrefix(outputDir, goPath)
			return packageName
		}
	}

	return ""
}
