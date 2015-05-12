package main

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/goldi/generator"
	"gopkg.in/alecthomas/kingpin.v1"
	"io/ioutil"
	"os"
	"strings"
)

var (
	app = kingpin.New("goldigen", "The goldi dependency injection container generator.\n\nSee https://github.com/fgrosse/goldi for further information.")

	inputFile     = app.Flag("in", "The input yaml file to generate type definitions from").Required().File()
	outputPath    = app.Flag("out", "The output file to save the generated go code").String()
	packageName   = app.Flag("package", "The name of the genarated package").Required().String()
	functionName  = app.Flag("function", fmt.Sprintf("The name of the generated function that must be called to register your types (default %q)", generator.DefaultFunctionName)).String()
	noInteraction = app.Flag("nointeraction", "Do not ask for any user input").Default("false").Bool()
	verbose       = app.Flag("verbose", "Print verbose output").Default("false").Bool()
	yes           = app.Flag("yes", "Answer all questions with yes").Default("false").Short('y').Bool()
)

func main() {
	defer panicHandler()
	app.Version(generator.Version)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := generator.NewConfig(*packageName, *functionName)
	gen := generator.New(config)
	output := &bytes.Buffer{}

	goPathChecker := generator.NewGoPathChecker()
	outputPackageName := goPathChecker.PackageName(*outputPath)

	inputStat, _ := (*inputFile).Stat()
	inputFileName := inputStat.Name()

	logVerboseGeneratorConfig(inputFileName, outputPackageName)
	err := gen.Generate(*inputFile, output, inputFileName, outputPackageName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *outputPath == "" {
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

func logVerboseGeneratorConfig(inputFileName, outputPackageName string) {
	log("Generating output from file %q", inputFileName)
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
		fmt.Printf("Output file %q does already exist. ", *outputPath)
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
	if *noInteraction {
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Println("Do you want me to overwrite that file? [yN] ")
	if *yes {
		fmt.Println("yes")
		return
	}

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
