package goldigen

import (
	. "github.com/fgrosse/goldi/tests/testAPI/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"github.com/fgrosse/goldi/generator"
	"strings"
)

var _ = Describe("Generator", func() {
	var (
		gen               *generator.Generator
		output            *bytes.Buffer
		inputName         = "input_file.yml"
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
		`
	)

	BeforeEach(func() {
		config := generator.NewConfig("foobar", "RegisterTypes")
		gen = generator.New(config)
		output = &bytes.Buffer{}
	})

	It("should generate valid go code", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
		Expect(output).To(BeValidGoCode())
	})

	It("should use the given package name", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
		Expect(output).To(DeclarePackage("foobar"))
	})

	Describe("generating import statements", func() {
		It("should import the goldi package", func() {
			Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
			Expect(output).To(BeValidGoCode())
			Expect(output).To(ImportPackage("github.com/fgrosse/goldi"))
		})

		It("should import the type packages", func() {
			Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
			Expect(output).To(BeValidGoCode())
			Expect(output).To(ImportPackage("github.com/fgrosse/graphigo"))
		})

		It("should not import the output package", func() {
			Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
			Expect(output).To(BeValidGoCode())
			Expect(output).NotTo(ImportPackage("github.com/fgrosse/some/thing"))
		})
	})

	It("should define the types in a global function", func() {
		Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
		// Note that NewFoo has no explicit package name since it is defined within the given outputPackageName
		Expect(output).To(ContainCode(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.RegisterType("goldi.test.foo", NewFoo)
				types.RegisterType("graphigo.client", graphigo.NewClient)
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
			Expect(gen.Generate(strings.NewReader(exampleYaml), output, inputName, outputPackageName)).To(Succeed())
			Expect(output).To(ContainCode(`
				func RegisterTypes(types goldi.TypeRegistry) {
					types.RegisterType("graphigo.client", graphigo.NewClient, "%graphigo.base_url%", 100)
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
		Expect(gen.Generate(strings.NewReader(invalidInput), output, inputName, outputPackageName)).
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
		Expect(gen.Generate(strings.NewReader(input), output, inputName, outputPackageName)).To(Succeed())
		Expect(output).To(ContainCode(fmt.Sprintf(`
			func RegisterTypes(types goldi.TypeRegistry) {
				types.RegisterType("test", bar.NewFoo, "%s")
			}
		`, "Hello\t\t\tWorld")))
	})
})
