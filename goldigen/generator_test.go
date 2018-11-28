package main_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fgrosse/goldi/goldigen"
	. "github.com/fgrosse/gomega-matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error from errorReader")
}

var _ = Describe("Generator", func() {
	var (
		gen               *main.Generator
		output            *bytes.Buffer
		inputPath         = "/absolute/path/conf/servo_types.yml"
		outputPath        = "/absolute/path/servo_types.go"
		outputPackageName = "github.com/fgrosse/some/thing"
		exampleYaml       = `
			types:
				goldi.test.foo:
					package: github.com/fgrosse/some/thing
					type:    Foo
					factory: NewFoo

				graphigo.client:
					package: github.com/fgrosse/graphigo
					type:    Graphigo
					factory: NewClient

				simple.struct:
					package: github.com/fgrosse/servo/example
					type:    MyStruct

				http_handler:
					package: github.com/fgrosse/servo/example
					func:    HandleHTTP

				logger:
					package: github.com/mgutz/logxi.v1
					package-name: log
					type:    Logger
					factory: New
					arguments: [ test ]
		`
	)

	BeforeEach(func() {
		config := main.NewConfig(outputPackageName, "RegisterTypes", inputPath, outputPath)
		gen = main.NewGenerator(config)
		output = &bytes.Buffer{}
	})

	It("should return error for invalid definition", func() {
		yaml := `
			types:
				goldi.test.foo:
					package: test
					factory: @foo_provider::NewFoo
					args:
						invalid yml`

		err := gen.Generate(strings.NewReader(yaml), output)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(HavePrefix("could not parse type definition: yaml: unmarshal errors:"))
	})

	It("should return error when cannot read from input", func() {
		Expect(gen.Generate(errorReader{}, output)).To(MatchError("could not parse type definition: error from errorReader"))
	})

	It("should not be necessary to quote type references because of the @", func() {
		// the @ has some special significance in yaml which we are going to ignore in goldigen
		yaml := `
			types:
				goldi.test.foo:
					package: test
					factory: @foo_provider::NewFoo
					args:
						- @bar
						- john.doe@example.com
						- 'alice@example.com'
						- "mallory@example.com"
						- There is an @ here
						- 1`

		Expect(gen.Generate(strings.NewReader(yaml), output)).To(Succeed())
		Expect(output).To(ContainCode(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.Register("goldi.test.foo", goldi.NewProxyType("foo_provider", "NewFoo", "@bar", "john.doe@example.com", "alice@example.com", "mallory@example.com", "There is an @ here", 1))
			}
		`))
	})

	It("should generate valid go code", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
		Expect(output).To(BeValidGoCode())
	})

	It("should use the given package name", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
		Expect(output).To(DeclarePackage("thing"))
	})

	Describe("generating import statements", func() {
		BeforeEach(func() {
			Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
			Expect(output).To(BeValidGoCode())
		})

		It("should import the goldi package", func() {
			Expect(output).To(ImportPackage("github.com/fgrosse/goldi"))
		})

		It("should import the type packages", func() {
			Expect(output).To(ImportPackage("github.com/fgrosse/graphigo"))
		})

		It("should not import the output package", func() {
			Expect(output).NotTo(ImportPackage(outputPackageName))
		})

		It("should not import type packages multuple times", func() {
			Expect(output).To(ImportPackage("github.com/fgrosse/servo/example"))
		})

		It("should import packages that contain a version", func() {
			Expect(output).To(ImportPackage("github.com/mgutz/logxi.v1"))
		})
	})

	It("should define the types in a global function", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
		// Note that NewFoo has no explicit package name since it is defined within the given outputPackageName
		Expect(output).To(ContainCode(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.RegisterAll(map[string]goldi.TypeFactory{
					"goldi.test.foo":  goldi.NewType(NewFoo),
					"graphigo.client": goldi.NewType(graphigo.NewClient),
					"http_handler":    goldi.NewFuncType(example.HandleHTTP),
					"logger":          goldi.NewType(log.New, "test"),
					"simple.struct":   goldi.NewStructType(new(example.MyStruct)),
				})
			}
		`))
	})

	Context("with parameters", func() {
		BeforeEach(func() {
			exampleYaml = `
			parameters:
				graphigo.base_url: https://example.com/graphigo:8443

			types:
				graphigo.client:
					package: github.com/fgrosse/graphigo
					type:    Graphigo
					factory: NewClient
					arguments:
						- "%graphigo.base_url%"
						- 100
		`
		})

		It("should define the types in a global function", func() {
			Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
			Expect(output).To(ContainCode(`
				func RegisterTypes(types goldi.TypeRegistry) {
					types.Register("graphigo.client", goldi.NewType(graphigo.NewClient, "%graphigo.base_url%", 100))
				}
			`))
		})
	})

	It("should validate the input", func() {
		invalidInput := `
			types:
				ok:
					package: some/package
					factory: NewFoo

				bad:
					type: TypeRegistry
					factory: NewTypeRegistry
		`
		Expect(gen.Generate(strings.NewReader(invalidInput), output)).
			To(MatchError(`type definition of "bad" is missing the required "package" key`))
	})

	It("should not replace tab characters in the middle of any value", func() {
		input := fmt.Sprintf(`
			types:
				test:
					package: foo/bar
					factory: NewFoo
					arguments:
            			- "%s"
		`, "Hello\t\t\tWorld")
		Expect(gen.Generate(strings.NewReader(input), output)).To(Succeed())
		Expect(output).To(ContainCode(fmt.Sprintf(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.Register("test", goldi.NewType(bar.NewFoo, "%s"))
			}
		`, "Hello\t\t\tWorld")))
	})

	It("should include the go generate code which was used to create this file", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output)).To(Succeed())
		Expect(output).To(ContainCode(fmt.Sprintf(
			`//go:generate goldigen --in "conf/servo_types.yml" --out "servo_types.go" --package %s --function RegisterTypes --overwrite --nointeraction`,
			outputPackageName,
		)))
	})

	It("should allow specifying configuration types", func() {
		input := `
			types:
				test:
					package: foo/bar
					factory: NewFoo
					configurator: ["@confoogurator", Configure]
		`
		Expect(gen.Generate(strings.NewReader(input), output)).To(Succeed())
		Expect(output).To(ContainCode(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.Register("test", goldi.NewConfiguredType(
					goldi.NewType(bar.NewFoo),
					"confoogurator", "Configure",
				))
			}
		`))
	})

	It("should log message in debug mode", func() {
		logger := new(bytes.Buffer)
		gen.Debug = true
		gen.Logger = logger
		gen.Generate(strings.NewReader(exampleYaml), output)
		Expect(logger.String()).NotTo(BeEmpty())
		gen.Debug = false
	})

})
