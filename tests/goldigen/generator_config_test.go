package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/generator"
)

var _ = Describe("GeneratorConfig", func() {
	Describe("NewGeneratorConfig", func() {
		It("should set the default type registration function name", func() {
			config := generator.NewConfig("package_name", "")
			Expect(config.FunctionName()).To(Equal(generator.DefaultFunctionName))
		})
	})
})
