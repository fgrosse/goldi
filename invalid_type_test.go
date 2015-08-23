package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"fmt"
)

var _ = Describe("InvalidType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewInvalidType(fmt.Errorf("Something bad happened"))
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := goldi.NewInvalidType(fmt.Errorf("Something bad happened"))
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(BeEmpty())
		})
	})

	Describe("Generate()", func() {
		It("should return its error", func() {
			config := map[string]interface{}{}
			container := goldi.NewContainer(goldi.NewTypeRegistry(), config)
			resolver := goldi.NewParameterResolver(container)

			e := fmt.Errorf("Something bad happened")
			typeDef := goldi.NewInvalidType(e)

			generated, err := typeDef.Generate(resolver)
			Expect(generated).To(BeNil())
			Expect(err).To(Equal(e))
		})
	})
})
