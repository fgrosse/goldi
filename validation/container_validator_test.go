package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/validation"
)

var _ = Describe("ContainerValidator", func() {
	var (
		registry  goldi.TypeRegistry
		config    map[string]interface{}
		container *goldi.Container
		validator *validation.ContainerValidator
	)

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
		config = map[string]interface{}{}
		container = goldi.NewContainer(registry, config)
		validator = validation.NewContainerValidator()
	})

	It("should return an error if an invalid type was registered", func() {
		registry.Register("main_type", goldi.NewFuncReferenceType("not_existent", "type"))

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when parameter has not been set", func() {
		typeDef := goldi.NewType(NewMockTypeWithArgs, "hello world", "%param%")
		registry.Register("main_type", typeDef)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a dependend type has not been registered", func() {
		typeDef := goldi.NewType(NewTypeForServiceInjection, "@injected_type")
		registry.Register("main_type", typeDef)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a direct circular type dependency exists", func() {
		injectedTypeID := "type_1"
		typeDef1 := goldi.NewType(NewTypeForServiceInjection, "@type_2")
		registry.Register(injectedTypeID, typeDef1)

		otherTypeID := "type_2"
		typeDef2 := goldi.NewType(NewTypeForServiceInjection, "@type_1")
		registry.Register(otherTypeID, typeDef2)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a transitive circular type dependency exists", func() {
		typeID1 := "type_1"
		typeDef1 := goldi.NewType(NewTypeForServiceInjection, "@type_2")
		registry.Register(typeID1, typeDef1)

		typeID2 := "type_2"
		typeDef2 := goldi.NewType(NewTypeForServiceInjection, "@type_3")
		registry.Register(typeID2, typeDef2)

		typeID3 := "type_3"
		typeDef3 := goldi.NewType(NewTypeForServiceInjection, "@type_1")
		registry.Register(typeID3, typeDef3)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should not return an error when everything is OK", func() {
		config["param"] = true
		registry.Register("injected_type",
			goldi.NewType(NewMockTypeWithArgs, "hello world", "%param%"),
		)

		registry.Register("main_type",
			goldi.NewType(NewTypeForServiceInjection, "@injected_type"),
		)

		registry.Register("foo_type",
			goldi.NewType(NewMockTypeWithArgs, "@injected_type::DoStuff", true),
		)

		Expect(validator.Validate(container)).To(Succeed())
	})

	Describe("MustValidate", func() {
		It("should panic if an error occurs", func() {
			typeDef := goldi.NewType(NewMockTypeWithArgs, "hello world", "%param%")
			registry.Register("main_type", typeDef)

			Expect(func() { validator.MustValidate(container) }).To(Panic())
		})

		It("should not panic if everything is ok", func() {
			config["param"] = true
			injectedTypeID := "injected_type"
			typeDef1 := goldi.NewType(NewMockTypeWithArgs, "hello world", "%param%")
			registry.Register(injectedTypeID, typeDef1)

			otherTypeID := "main_type"
			typeDef2 := goldi.NewType(NewTypeForServiceInjection, "@injected_type")
			registry.Register(otherTypeID, typeDef2)

			Expect(func() { validator.MustValidate(container) }).NotTo(Panic())
		})
	})
})
