package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeAlias", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewTypeAlias("foo")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("Arguments()", func() {
		It("should return the aliased service ID", func() {
			typeDef := goldi.NewTypeAlias("foo")
			Expect(typeDef.Arguments()).To(Equal([]interface{}{"@foo"}))
		})
	})

	Describe("Generate()", func() {
		var (
			container *goldi.Container
			resolver  *goldi.ParameterResolver
		)

		BeforeEach(func() {
			config := map[string]interface{}{}
			container = goldi.NewContainer(goldi.NewTypeRegistry(), config)
			resolver = goldi.NewParameterResolver(container)
		})

		It("should act as alias for the actual type", func() {
			container.Register("foo", goldi.NewStructType(testAPI.MockType{}, "I was created by @foo"))
			alias := goldi.NewTypeAlias("foo")

			generated := alias.Generate(resolver)
			Expect(generated).To(BeAssignableToTypeOf(&testAPI.MockType{}))
			Expect(generated.(*testAPI.MockType).StringParameter).To(Equal("I was created by @foo"))
		})

		It("should work with func reference types", func() {
			container.Register("foo", goldi.NewStructType(testAPI.MockType{}, "I was created by @foo"))
			alias := goldi.NewTypeAlias("foo::ReturnString")

			generated := alias.Generate(resolver)
			Expect(generated).To(BeAssignableToTypeOf(func(string) string { return "" }))
			Expect(generated.(func(string) string)("TEST")).To(Equal("I was created by @foo TEST"))
		})
	})
})
