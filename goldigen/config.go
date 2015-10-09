package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultFunctionName is the name of the registration function that is used if nothing else has been specified.
const DefaultFunctionName = "RegisterTypes"

// Config is the goldigen configuration.
type Config struct {
	Package      string
	FunctionName string
	InputPath    string
	OutputPath   string
}

// NewConfig creates a new Config with the given parameters.
// This function will panic if completePackage is empty.
// If you pass an empty function name the default function name will be assumed
func NewConfig(completePackage, functionName, inputPath, outputPath string) Config {
	if completePackage == "" {
		panic(fmt.Errorf("Output package name can not be empty"))
	}

	if functionName == "" {
		functionName = DefaultFunctionName
	}

	return Config{completePackage, functionName, inputPath, outputPath}
}

// PackageName returns the name of the configured package.
func (c Config) PackageName() string {
	packageParts := strings.Split(c.Package, "/")

	return packageParts[len(packageParts)-1]
}

// OutputName returns the base name of the configured output path
func (c Config) OutputName() string {
	return filepath.Base(c.OutputPath)
}

// InputName returns the input file path relative to the output directory.
func (c Config) InputName() string {
	inputFile, err := filepath.Rel(filepath.Dir(c.OutputPath), c.InputPath)
	if err != nil {
		panic(err)
	}

	return inputFile
}
