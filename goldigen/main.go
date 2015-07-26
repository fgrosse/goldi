package main

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/goldi/generator"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
)

var (
	app = kingpin.New("goldigen", "The goldi dependency injection container generator.\n\nSee https://github.com/fgrosse/goldi for further information.")

	inputFile     = app.Flag("in", "The input yaml file to generate type definitions from").Required().File()
	outputPath    = app.Flag("out", "The output file to save the generated go code").String()
	packageName   = app.Flag("package", "The name of the genarated package").String()
	functionName  = app.Flag("function", fmt.Sprintf("The name of the generated function that must be called to register your types (default %q)", generator.DefaultFunctionName)).String()
	noInteraction = app.Flag("nointeraction", "Do not ask for any user input").Default("false").Bool()
	verbose       = app.Flag("verbose", "Print verbose output").Default("false").Bool()
	overwrite     = app.Flag("overwrite", "Overwrite any existing files").Default("false").Short('y').Bool()
	forceStdOut   = app.Flag("echo", "Echo the generated code to std out even if a output path is given").Default("false").Bool()
)

func main() {
	defer panicHandler()
	app.Version(generator.Version)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	inputPath, _ := filepath.Abs((*inputFile).Name())
	if *outputPath != "" {
		*outputPath, _ = filepath.Abs(*outputPath)
	}

	outputPackageName := determineOutputPackageName()
	config := generator.NewConfig(outputPackageName, *functionName, inputPath, *outputPath)
	gen := generator.New(config)
	output := &bytes.Buffer{}

	if *verbose {
		gen.Debug = true
	}

	logVerboseGeneratorConfig(inputPath, outputPackageName)
	err := gen.Generate(*inputFile, output)
	if err != nil {
		log(err.Error())
		os.Exit(1)
	}

	if *outputPath == "" || *forceStdOut {
		logVerbose("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(output.String())
		return
	}

	writeOutputFile(output)
}

func panicHandler() {
	if r := recover(); r != nil {
		log("FATAL ERROR: %s", r)
		os.Exit(1)
	}
}

func determineOutputPackageName() string {
	outputPackageName := *packageName
	if outputPackageName != "" {
		return outputPackageName
	}

	goPathChecker := generator.NewGoPathChecker(*verbose)
	outputPackageName = goPathChecker.PackageName(*outputPath)
	logVerbose("Package name for output path %q is %q", *outputPath, outputPackageName)

	if outputPackageName != "" {
		return outputPackageName
	}

	if *outputPath != "" {
		log("Could not determine the output package name for %q", *outputPath)
	}

	return ask("Output package name: ")
}

func ask(question string) string {
	if *noInteraction {
		os.Exit(1)
	}

	log(question)
	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(answer)
}

func logVerboseGeneratorConfig(inputPath, outputPackageName string) {
	logVerbose("Generating output from file %q", inputPath)
	if *outputPath != "" {
		logVerbose("Output will be saved to %q", *outputPath)
	}

	if outputPackageName == "" {
		logVerbose("Output package name is empty")
	} else {
		logVerbose("Output package name is %q", outputPackageName)
	}
}

func logVerbose(message string, args ...interface{}) {
	if *verbose == false {
		return
	}

	log(message, args...)
}

func log(message string, args ...interface{}) {
	writer := os.Stdout
	if *outputPath == "" {
		// since we already output the generated code on stdout we print messages on stderr
		writer = os.Stderr
	}

	fmt.Fprintf(writer, message+"\n", args...)
}

func writeOutputFile(output *bytes.Buffer) {
	if _, err := os.Stat(*outputPath); err == nil {
		checkUserWantsToOverwriteFile()
	}

	err := ioutil.WriteFile(*outputPath, output.Bytes(), 0644)
	if err != nil {
		log("Error while writing output file: %s", err)
		os.Exit(1)
	}
	log("Successfully wrote %d bytes to %q", output.Len(), *outputPath)
}

func checkUserWantsToOverwriteFile() {
	if *overwrite {
		return
	}

	log("Output file %q does already exist.", *outputPath)
	answer := ask("Do you want me to overwrite that file? [yN] ")
	answer = strings.ToLower(answer)
	if answer == "" || answer == "n" {
		log("Output has NOT been saved")
		os.Exit(1)
	}
}
