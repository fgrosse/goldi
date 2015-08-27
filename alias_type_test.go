package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
)

func ExampleNewAliasType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	container.Register("logger", goldi.NewStructType(SimpleLogger{}))
	container.Register("default_logger", goldi.NewAliasType("logger"))
	container.Register("logging_func", goldi.NewAliasType("logger::DoStuff"))

	fmt.Printf("logger:         %T\n", container.MustGet("logger"))
	fmt.Printf("default_logger: %T\n", container.MustGet("default_logger"))
	fmt.Printf("logging_func:   %T\n", container.MustGet("logging_func"))
	// Output:
	// logger:         *goldi_test.SimpleLogger
	// default_logger: *goldi_test.SimpleLogger
	// logging_func:   func(string) string
}

// ExampleNewAliasType_ prevents godoc from printing the whole content of this file as example
func ExampleNewAliasType_() {}

var _ = Describe("aliasType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewAliasType("foo")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("Arguments()", func() {
		It("should return the aliased service ID", func() {
			typeDef := goldi.NewAliasType("foo")
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
			container.Register("foo", goldi.NewStructType(Foo{}, "I was created by @foo"))
			alias := goldi.NewAliasType("foo")

			generated, err := alias.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generated).To(BeAssignableToTypeOf(&Foo{}))
			Expect(generated.(*Foo).Value).To(Equal("I was created by @foo"))
		})

		It("should work with func reference types", func() {
			container.Register("foo", goldi.NewStructType(Foo{}, "I was created by @foo"))
			alias := goldi.NewAliasType("foo::ReturnString")

			generated, err := alias.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generated).To(BeAssignableToTypeOf(func(string) string { return "" }))
			Expect(generated.(func(string) string)("TEST")).To(Equal("I was created by @foo TEST"))
		})
	})
})
