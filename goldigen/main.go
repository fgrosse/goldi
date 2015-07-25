package main

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/goldi/generator"
	"gopkg.in/alecthomas/kingpin.v1"
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
	if *outputPath == "" {
		// TODO require a package name
	} else {
		*outputPath, _ = filepath.Abs(*outputPath)
	}

	outputPackageName := *packageName
	if outputPackageName == "" {
		goPathChecker := generator.NewGoPathChecker(*verbose)
		outputPackageName = goPathChecker.PackageName(*outputPath)
		log("Package name for output path %q is %q", *outputPath, outputPackageName)
	}

	config := generator.NewConfig(outputPackageName, *functionName, inputPath, *outputPath)
	gen := generator.New(config)
	output := &bytes.Buffer{}

	if *verbose {
		gen.Debug = true
	}

	logVerboseGeneratorConfig(inputPath, outputPackageName)
	err := gen.Generate(*inputFile, output)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *outputPath == "" || *forceStdOut {
		log("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(output.String())
		return
	}

	writeOutputFile(output)
}

func panicHandler() {
	if r := recover(); r != nil {
		fmt.Printf("FATAL ERROR: %s", r)
		os.Exit(1)
	}
}

func logVerboseGeneratorConfig(inputPath, outputPackageName string) {
	log("Generating output from file %q", inputPath)
	if *outputPath != "" {
		log("Output will be saved to %q", *outputPath)
	}

	if outputPackageName == "" {
		log("Output package name is empty")
	} else {
		log("Output package name is %q", outputPackageName)
	}
}

func log(message string, args ...interface{}) {
	if *verbose == false {
		return
	}

	fmt.Printf(message+"\n", args...)
}

func writeOutputFile(output *bytes.Buffer) {
	if _, err := os.Stat(*outputPath); err == nil {
		checkUserWantsToOverwriteFile()
	}

	err := ioutil.WriteFile(*outputPath, output.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error while writing output file: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully wrote %d bytes to %q\n", output.Len(), *outputPath)
}

func checkUserWantsToOverwriteFile() {
	if *overwrite {
		return
	}

	fmt.Printf("Output file %q does already exist.\n", *outputPath)
	if *noInteraction {
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Print("Do you want me to overwrite that file? [yN] ")
	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		panic(err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "" || answer == "n" {
		fmt.Println("Output has NOT been saved")
		os.Exit(1)
	}
}
