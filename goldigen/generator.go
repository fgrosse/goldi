package goldigen

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

// The Generator is used to generate compilable go code from a yaml configuration
type Generator struct {
	Config GeneratorConfig
}

// NewGenerator creates a new Generator instance
func NewGenerator(config GeneratorConfig) *Generator {
	return &Generator{config}
}

// Generate reads a yaml type configuration from the `input` and writes the corresponding go code to the `output`.
// You have to provide the package name that should be used when generating the output.
func (g *Generator) Generate(input io.Reader, output io.Writer) error {
	conf, err := g.parseInput(input)
	if err != nil {
		return fmt.Errorf("could not parse type definition: %s", err)
	}

	err = conf.Validate()
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "package %s\n\n", g.Config.PackageName)
	g.generateImports(conf, output)
	g.typeRegistrationFunction(conf, output)
	return nil
}

func (g *Generator) parseInput(input io.Reader) (*TypesConfiguration, error) {
	inputData, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	inputData = g.sanitizeInput(inputData)

	var config TypesConfiguration
	err = yaml.Unmarshal(inputData, &config)
	return &config, err
}

func (g *Generator) sanitizeInput(input []byte) []byte {
	sanitizedInput := &bytes.Buffer{}
	line := &bytes.Buffer{}
	lineBeginning := true
	for _, c := range input {
		switch c {
		case '\n':
			if strings.TrimSpace(line.String()) != "" {
				sanitizedInput.Write(line.Bytes())
				sanitizedInput.WriteByte('\n')
				line.Reset()
				lineBeginning = true
			}
		case '\t':
			if lineBeginning {
				line.WriteString("    ")
			} else {
				line.WriteByte(c)
			}
		case ' ':
			line.WriteByte(c)
		default:
			lineBeginning = false
			line.WriteByte(c)
		}
	}

	sanitizedInput.Write(line.Bytes())
	return sanitizedInput.Bytes()
}

func (g *Generator) generateImports(conf *TypesConfiguration, output io.Writer) {
	packages := conf.Packages("github.com/fgrosse/goldi")

	fmt.Fprint(output, "import (\n")
	for _, pkg := range packages {
		fmt.Fprintf(output, "\t%q\n", pkg)
	}

	fmt.Fprint(output, ")\n\n")
}

func (g *Generator) typeRegistrationFunction(conf *TypesConfiguration, output io.Writer) {
	fmt.Fprintf(output, "func %s(types goldi.TypeRegistry) {\n", g.Config.FunctionName)
	typeIDs := make([]string, len(conf.Types))
	i := 0
	for typeID, _ := range conf.Types {
		typeIDs[i] = typeID
		i++
	}
	sort.Strings(typeIDs)

	for _, typeID := range typeIDs {
		typeDef := conf.Types[typeID]
		arguments := []string{
			fmt.Sprintf("%q", typeID),
			fmt.Sprintf("%s.%s", typeDef.PackageName(), typeDef.FactoryMethod),
		}
		arguments = append(arguments, typeDef.Arguments()...)
		fmt.Fprint(output, "\ttypes.RegisterType(")
		fmt.Fprint(output, strings.Join(arguments, ", "))
		fmt.Fprint(output, ")\n")
	}
	fmt.Fprint(output, "}\n")
}
