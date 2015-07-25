package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/generator"
)

var _ = Describe("GeneratorConfig", func() {
	Describe("NewGeneratorConfig", func() {
		It("should set the default type registration function name", func() {
			config := generator.NewConfig("package_name", "", "", "")
			Expect(config.FunctionName).To(Equal(generator.DefaultFunctionName))
		})
	})

	Describe("PackageName", func() {
		It("should only return the package name", func() {
			config := generator.NewConfig("github.com/fgrosse/servo", "", "", "")
			Expect(config.Package).To(Equal("github.com/fgrosse/servo"))
			Expect(config.PackageName()).To(Equal("servo"))
		})
	})

	Describe("InputName", func() {
		It("should only return the file name", func() {
			config := generator.NewConfig("github.com/fgrosse/servo", "", "/home/fgrosse/tmp/types.yml", "")
			Expect(config.InputName()).To(Equal("types.yml"))
		})
	})
})
