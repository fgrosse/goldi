package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("FuncReferenceType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = NewFuncReferenceType("my_controller", "FancyAction")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewFuncReferenceType()", func() {
		It("should return an invalid type if the method name is not exported", func() {
			t := NewFuncReferenceType("foo", "doStuff")
			Expect(IsValid(t)).To(BeFalse())
			Expect(t.(*invalidType).Err).To(MatchError(fmt.Errorf(`can not use unexported method "doStuff" as second argument to NewFuncReferenceType`)))
		})
	})

	Describe("Arguments()", func() {
		It("should return the referenced service ID", func() {
			typeDef := NewFuncReferenceType("my_controller", "FancyAction")
			Expect(typeDef.Arguments()).To(Equal([]interface{}{"@my_controller"}))
		})
	})

	Describe("Generate()", func() {
		var (
			container *Container
			resolver  *ParameterResolver
		)

		BeforeEach(func() {
			config := map[string]interface{}{}
			container = NewContainer(NewTypeRegistry(), config)
			resolver = NewParameterResolver(container)
		})

		It("should get the correct method of the referenced type", func() {
			container.Register("foo", NewStructType(tests.MockType{}, "I was created by @foo"))
			typeDef := NewFuncReferenceType("foo", "ReturnString")

			generated, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generated).To(BeAssignableToTypeOf(func(string) string { return "" }))
			Expect(generated.(func(string) string)("TEST")).To(Equal("I was created by @foo TEST"))
		})

		It("should return an error if the referenced type has no such method", func() {
			container.Register("foo", NewStructType(tests.MockType{}))
			typeDef := NewFuncReferenceType("foo", "ThisMethodDoesNotExist")

			_, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("could not generate func reference type @foo::ThisMethodDoesNotExist : method does not exist"))
		})
	})
})
