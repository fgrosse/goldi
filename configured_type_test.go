package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
)

func ExampleNewConfiguredType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	// this example configurator accepts a Foo type and will set its Value field to the given value
	configurator := &MyConfigurator{ConfiguredValue: "success!"}

	// register the configurator under a type ID
	container.Register("configurator_type", goldi.NewInstanceType(configurator))

	// create the type that should be configured
	embeddedType := goldi.NewStructType(Foo{})
	container.Register("foo", goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure"))

	fmt.Println(container.MustGet("foo").(*Foo).Value)
	// Output:
	// success!
}

// ExampleNewConfiguredType_ prevents godoc from printing the whole content of this file as example
func ExampleNewConfiguredType_() {}

var _ = Describe("configuredType", func() {
	var embeddedType goldi.TypeFactory
	BeforeEach(func() {
		embeddedType = goldi.NewStructType(Foo{})
	})

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewConfiguredType()", func() {
		Context("with invalid argument", func() {
			It("should return an invalid type if the embedded type is nil", func() {
				typeDef := goldi.NewConfiguredType(nil, "configurator_type", "Configure")
				Expect(goldi.IsValid(typeDef)).To(BeFalse())
			})

			It("should return an invalid type if either the configurator ID or method is empty", func() {
				Expect(goldi.IsValid(goldi.NewConfiguredType(embeddedType, "", ""))).To(BeFalse())
				Expect(goldi.IsValid(goldi.NewConfiguredType(embeddedType, "configurator_type", ""))).To(BeFalse())
				Expect(goldi.IsValid(goldi.NewConfiguredType(embeddedType, "", "configure"))).To(BeFalse())
			})

			It("should return an invalid type if the configurator method is not exported", func() {
				Expect(goldi.IsValid(goldi.NewConfiguredType(embeddedType, "configurator_type", "configure"))).To(BeFalse())
			})
		})

		Context("with valid arguments", func() {
			It("should create the type", func() {
				typeDef := goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
				Expect(typeDef).NotTo(BeNil())
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return the arguments of the embedded type and the configurator as type ID", func() {
			embeddedType = goldi.NewStructType(Foo{}, "%param_of_embedded%", "another param")
			typeDef := goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(HaveLen(3))
			Expect(typeDef.Arguments()).To(ContainElement("%param_of_embedded%"))
			Expect(typeDef.Arguments()).To(ContainElement("another param"))
			Expect(typeDef.Arguments()).To(ContainElement("@configurator_type"))
		})
	})

	Describe("Generate()", func() {
		var (
			config    = map[string]interface{}{}
			container *goldi.Container
			resolver  *goldi.ParameterResolver
		)

		BeforeEach(func() {
			container = goldi.NewContainer(goldi.NewTypeRegistry(), config)
			resolver = goldi.NewParameterResolver(container)
		})

		It("should get the embedded type and configurator and configure it", func() {
			typeDef := goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
			configurator := &MyConfigurator{ConfiguredValue: "success!"}
			container.Register("configurator_type", goldi.NewInstanceType(configurator))

			generatedType, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generatedType).NotTo(BeNil())
			Expect(generatedType).To(BeAssignableToTypeOf(&Foo{}))
			Expect(generatedType.(*Foo).Value).To(Equal("success!"))
		})

		It("should return an error if the embedded type can not be generated", func() {
			invalidType := goldi.NewStructType(nil)
			typeDef := goldi.NewConfiguredType(invalidType, "configurator_type", "Configure")
			configurator := &MyConfigurator{ConfiguredValue: "should not happen"}
			container.Register("configurator_type", goldi.NewInstanceType(configurator))

			generatedType, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("can not generate configured type: the given struct is nil"))
			Expect(generatedType).To(BeNil())
		})

		It("should return an error if the configurator returns an error", func() {
			typeDef := goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
			configurator := &MyConfigurator{ReturnError: true}
			container.Register("configurator_type", goldi.NewInstanceType(configurator))

			generatedType, err := typeDef.Generate(resolver)
			Expect(err).To(MatchError("can not configure type: this is the error message from the tests.MockTypeConfigurator"))
			Expect(generatedType).To(BeNil())
		})
	})
})
