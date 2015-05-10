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

const GoldiVersion = "0.1.0"

var (
	app = kingpin.New("goldigen", "The goldi dependency injection container generator.\n\nSee https://github.com/fgrosse/goldi for further information.")

	inputFile     = app.Flag("in", "The input yaml file to generate type definitions from").Required().File()
	outputPath    = app.Flag("out", "The output file to save the generated go code").String()
	packageName   = app.Flag("package", "The name of the genarated package").Required().String()
	functionName  = app.Flag("function", fmt.Sprintf("The name of the generated function that must be called to register your types (default %q)", generator.DefaultFunctionName)).String()
	noInteraction = app.Flag("nointeraction", "Do not ask for any user input").Default("false").Bool()
)

func main() {
	app.Version(GoldiVersion)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := generator.NewConfig(*packageName, *functionName)
	gen := generator.New(config)
	output := &bytes.Buffer{}
	err := gen.Generate(*inputFile, output)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *outputPath == "" {
		fmt.Println(output.String())
		return
	}

	if _, err := os.Stat(*outputPath); err == nil {
		fmt.Printf("Output file %q does already exist. ", *outputPath)
		if *noInteraction {
			fmt.Println("")
			os.Exit(1)
		}

		fmt.Println("Do you want me to overwrite that file? [yN] ", *outputPath)
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

	err = ioutil.WriteFile(*outputPath, output.Bytes(), 0644)
	fmt.Printf("Successfully wrote %d bytes to %q\n", output.Len(), *outputPath)
}
