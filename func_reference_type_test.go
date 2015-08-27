package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
)

func ExampleNewFuncReferenceType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	logger := new(SimpleLogger)
	container.Register("logger", goldi.NewInstanceType(logger))
	container.Register("log_func", goldi.NewFuncReferenceType("logger", "DoStuff"))

	f := container.MustGet("log_func").(func(string) string)
	fmt.Println(f("Hello World")) // executes logger.DoStuff
	// Output:
	// Hello World
}

// ExampleNewFuncReferenceType_ prevents godoc from printing the whole content of this file as example
func ExampleNewFuncReferenceType_() {}

var _ = Describe("funcReferenceType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewFuncReferenceType("my_controller", "FancyAction")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewFuncReferenceType()", func() {
		It("should return an invalid type if the method name is not exported", func() {
			t := goldi.NewFuncReferenceType("foo", "doStuff")
			Expect(goldi.IsValid(t)).To(BeFalse())
			Expect(t).To(MatchError(`can not use unexported method "doStuff" as second argument to NewFuncReferenceType`))
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
			container.Register("foo", goldi.NewStructType(Foo{}, "I was created by @foo"))
			typeDef := goldi.NewFuncReferenceType("foo", "ReturnString")

			generated, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generated).To(BeAssignableToTypeOf(func(string) string { return "" }))
			Expect(generated.(func(string) string)("TEST")).To(Equal("I was created by @foo TEST"))
		})

		It("should return an error if the referenced type can not be generated", func() {
			container.Register("foo", goldi.NewStructType(nil))
			typeDef := goldi.NewFuncReferenceType("foo", "DoStuff")

			_, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError(`could not generate func reference type @foo::DoStuff : goldi: error while generating type "foo": the given struct is nil`))
		})

		It("should return an error if the referenced type has no such method", func() {
			container.Register("foo", goldi.NewStructType(Foo{}))
			typeDef := goldi.NewFuncReferenceType("foo", "ThisMethodDoesNotExist")

			_, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("could not generate func reference type @foo::ThisMethodDoesNotExist : method does not exist"))
		})
	})
})
