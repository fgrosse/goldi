package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
)

func ExampleNewProxyType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	container.Register("logger_provider", goldi.NewStructType(LoggerProvider{}))
	container.Register("logger", goldi.NewProxyType("logger_provider", "GetLogger", "My logger"))

	l := container.MustGet("logger").(*SimpleLogger)
	fmt.Printf("%s: %T", l.Name, l)
	// Output:
	// My logger: *goldi_test.SimpleLogger
}

// ExampleNewProxyType_ prevents godoc from printing the whole content of this file as example
func ExampleNewProxyType_() {}

var _ = Describe("proxyType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewProxyType("logger_provider", "GetLogger", "My logger")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewProxyType()", func() {
		It("should return an invalid type if the method name is not exported", func() {
			t := goldi.NewProxyType("logger_provider", "getLogger", "My logger")
			Expect(goldi.IsValid(t)).To(BeFalse())
			Expect(t).To(MatchError(`can not use unexported method "getLogger" as second argument to NewProxyType`))
		})
	})

	Describe("Arguments()", func() {
		It("should return the referenced service ID", func() {
			typeDef := goldi.NewProxyType("logger_provider", "GetLogger", "My logger")
			Expect(typeDef.Arguments()).To(Equal([]interface{}{"@logger_provider", "My logger"}))
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
			container.Register("logger_provider", goldi.NewStructType(LoggerProvider{}))
			typeDef := goldi.NewProxyType("logger_provider", "GetLogger", "My logger")

			generated, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generated).To(BeAssignableToTypeOf(&SimpleLogger{}))
			Expect(generated.(*SimpleLogger).Name).To(Equal("My logger"))
		})

		It("should return an error if the referenced type has no such method", func() {
			typeDef := goldi.NewProxyType("foobar", "DoStuff")

			_, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("could not generate proxy type @foobar::DoStuff : type foobar does not exist"))
		})

		It("should return an error if the referenced type has no such method", func() {
			container.Register("logger_provider", goldi.NewStructType(LoggerProvider{}))
			typeDef := goldi.NewProxyType("logger_provider", "ThisMethodDoesNotExist", "foobar")

			_, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("could not generate proxy type @logger_provider::ThisMethodDoesNotExist : method does not exist"))
		})
	})
})
