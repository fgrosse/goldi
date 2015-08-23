package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("ConfiguredType", func() {
	var embeddedType goldi.TypeFactory
	BeforeEach(func() {
		embeddedType = goldi.NewStructType(tests.MockType{})
	})

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewConfiguredType(embeddedType, "configurator_type", "Configure")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewConfiguredType()", func() {
		Context("with invalid argument", func() {
			It("should panic if the embedded type is nil", func() {
				Expect(func() { goldi.NewConfiguredType(nil, "configurator_type", "Configure") }).To(Panic())
			})

			It("should panic if either the configurator ID or method is empty", func() {
				Expect(func() { goldi.NewConfiguredType(embeddedType, "", "") }).To(Panic())
				Expect(func() { goldi.NewConfiguredType(embeddedType, "configurator_type", "") }).To(Panic())
				Expect(func() { goldi.NewConfiguredType(embeddedType, "", "configure") }).To(Panic())
			})

			It("should panic if the configurator method is not exported", func() {
				Expect(func() { goldi.NewConfiguredType(embeddedType, "configurator_type", "configure") }).To(Panic())
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
			embeddedType = goldi.NewStructType(tests.MockType{}, "%param_of_embedded%", "another param")
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
			container.Register("configurator_type", goldi.NewType(tests.NewMockTypeConfigurator, "~~ configured ~~"))

			generatedType, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generatedType).NotTo(BeNil())
			Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))
			Expect(generatedType.(*tests.MockType).StringParameter).To(Equal("~~ configured ~~"))
		})
	})
})
