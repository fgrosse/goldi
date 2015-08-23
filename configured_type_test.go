package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("ConfiguredType", func() {
	var embeddedType TypeFactory
	BeforeEach(func() {
		embeddedType = NewStructType(tests.MockType{})
	})

	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = NewConfiguredType(embeddedType, "configurator_type", "Configure")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewConfiguredType()", func() {
		Context("with invalid argument", func() {
			It("should return an invalid type if the embedded type is nil", func() {
				typeDef := NewConfiguredType(nil, "configurator_type", "Configure")
				Expect(IsValid(typeDef)).To(BeFalse())
			})

			It("should return an invalid type if either the configurator ID or method is empty", func() {
				Expect(IsValid(NewConfiguredType(embeddedType, "", ""))).To(BeFalse())
				Expect(IsValid(NewConfiguredType(embeddedType, "configurator_type", ""))).To(BeFalse())
				Expect(IsValid(NewConfiguredType(embeddedType, "", "configure"))).To(BeFalse())
			})

			It("should return an invalid type if the configurator method is not exported", func() {
				Expect(IsValid(NewConfiguredType(embeddedType, "configurator_type", "configure"))).To(BeFalse())
			})
		})

		Context("with valid arguments", func() {
			It("should create the type", func() {
				typeDef := NewConfiguredType(embeddedType, "configurator_type", "Configure")
				Expect(typeDef).NotTo(BeNil())
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return the arguments of the embedded type and the configurator as type ID", func() {
			embeddedType = NewStructType(tests.MockType{}, "%param_of_embedded%", "another param")
			typeDef := NewConfiguredType(embeddedType, "configurator_type", "Configure")
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
			container *Container
			resolver  *ParameterResolver
		)

		BeforeEach(func() {
			container = NewContainer(NewTypeRegistry(), config)
			resolver = NewParameterResolver(container)
		})

		It("should get the embedded type and configurator and configure it", func() {
			typeDef := NewConfiguredType(embeddedType, "configurator_type", "Configure")
			container.Register("configurator_type", NewType(tests.NewMockTypeConfigurator, "~~ configured ~~"))

			generatedType, err := typeDef.Generate(resolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(generatedType).NotTo(BeNil())
			Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))
			Expect(generatedType.(*tests.MockType).StringParameter).To(Equal("~~ configured ~~"))
		})
	})
})
