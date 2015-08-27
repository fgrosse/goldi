package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
)

func ExampleNewInstanceType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	myInstance := new(SimpleLogger)
	myInstance.Name = "Foobar" // you can configure the instance in your code

	// now register this instance as a type
	container.Register("logger", goldi.NewInstanceType(myInstance))

	// each reference to the "logger" type will now be resolved to that instance
	fmt.Println(container.MustGet("logger").(*SimpleLogger).Name)
	// Output:
	// Foobar
}

// ExampleNewInstanceType_ prevents godoc from printing the whole content of this file as example
func ExampleNewInstanceType_() {}

var _ = Describe("instanceType", func() {
	var resolver *goldi.ParameterResolver

	BeforeEach(func() {
		container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})
		resolver = goldi.NewParameterResolver(container)
	})

	It("should return an invalid type if NewInstanceType is called with nil", func() {
		Expect(goldi.IsValid(goldi.NewInstanceType(nil))).To(BeFalse())
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := goldi.NewInstanceType(NewFoo())
			Expect(typeDef.Arguments()).To(BeEmpty())
		})
	})

	Describe("Generate", func() {
		It("should always return the given instance", func() {
			instance := NewFoo()
			factory := goldi.NewInstanceType(instance)

			for i := 0; i < 3; i++ {
				generateResult, err := factory.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generateResult == instance).To(BeTrue(),
					fmt.Sprintf("generateResult (%p) should point to the same instance as instance (%p)", generateResult, instance),
				)
			}
		})
	})

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewInstanceType("foo")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})
})
