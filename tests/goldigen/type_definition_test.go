package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/goldigen"
)

var _ = Describe("TypeDefinition", func() {
	It("should return an error if a type definition is missing the factory method", func() {
		t := goldigen.TypeDefinition{
			Package:       "foo/bar",
			FactoryMethod: "NewBaz",
		}
		Expect(t.PackageName()).To(Equal("bar"))
	})

	It("should return all parameters such that they can be used in go code directly", func() {
		t := goldigen.TypeDefinition{
			Package:       "foo/bar",
			FactoryMethod: "NewBaz",
			RawArguments: []interface{}{
				"Hello World",
				true,
				42,
				3.1415,
				"%some_parameter%",
				"Hello\t\tWorld",
			},
		}

		arguments := t.Arguments()
		Expect(arguments).To(HaveLen(6))
		Expect(arguments[0]).To(Equal(`"Hello World"`))
		Expect(arguments[1]).To(Equal(`true`))
		Expect(arguments[2]).To(Equal(`42`))
		Expect(arguments[3]).To(Equal(`3.1415`))
		Expect(arguments[4]).To(Equal(`"%some_parameter%"`))
		Expect(arguments[5]).To(Equal("\"Hello\t\tWorld\""))
	})
})
