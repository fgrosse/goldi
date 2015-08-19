package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("FuncReferenceType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewFuncReferenceType("my_controller", "FancyAction")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewFuncReferenceType()", func() {
		It("should panic if the method name is not exported", func() {
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected NewFuncReferenceType to panic")
				Expect(r).To(MatchError(fmt.Errorf(`can not use unexported method "doStuff" as second argument to NewFuncReferenceType`)))
			}()

			goldi.NewFuncReferenceType("foo", "doStuff")
		})
	})

	Describe("Arguments()", func() {
		It("should return the referenced service ID", func() {
			typeDef := goldi.NewFuncReferenceType("my_controller", "FancyAction")
			Expect(typeDef.Arguments()).To(Equal([]interface{}{"@my_controller"}))
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

		It("should get the correct method of the referenced type", func() {
			container.Register("foo", goldi.NewStructType(testAPI.MockType{}, "I was created by @foo"))
			typeDef := goldi.NewFuncReferenceType("foo", "ReturnString")

			generated := typeDef.Generate(resolver)
			Expect(generated).To(BeAssignableToTypeOf(func(string) string { return "" }))

			Expect(generated.(func(string) string)("TEST")).To(Equal("I was created by @foo TEST"))
		})

		It("panic if the referenced type has no such method", func() {
			container.Register("foo", goldi.NewStructType(testAPI.MockType{}))
			typeDef := goldi.NewFuncReferenceType("foo", "ThisMethodDoesNotExist")

			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected Generate to panic")
				Expect(r).To(MatchError(fmt.Errorf("could not generate func reference type @foo::ThisMethodDoesNotExist method does not exist")))
			}()

			typeDef.Generate(resolver)
		})
	})
})
