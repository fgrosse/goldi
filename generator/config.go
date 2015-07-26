package generator

import (
	"strings"
	"path/filepath"
	"fmt"
)

const DefaultFunctionName = "RegisterTypes"

type Config struct {
	Package      string
	FunctionName string
	InputPath    string
	OutputPath   string
}

func NewConfig(completePackage, functionName, inputPath, outputPath string) Config {
	if completePackage == "" {
		panic(fmt.Errorf("Output package name can not be empty"))
	}

	if functionName == "" {
		functionName = DefaultFunctionName
	}

	return Config{completePackage, functionName, inputPath, outputPath}
}

func (c Config) PackageName() string {
	packageParts := strings.Split(c.Package, "/")

	return packageParts[len(packageParts)-1]
}

func (c Config) InputName() string {
	return filepath.Base(c.InputPath)
}
