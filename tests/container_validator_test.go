package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("ContainerValidator", func() {
	var (
		registry  goldi.TypeRegistry
		config    map[string]interface{}
		container *goldi.Container
		validator *goldi.ContainerValidator
	)

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
		config = map[string]interface{}{}
		container = goldi.NewContainer(registry, config)
		validator = goldi.NewContainerValidator()
	})

	It("should return an error when parameter has not been set", func() {
		typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, "hello world", "%param%")
		registry.Register("goldi.main_type", typeDef)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a dependend type has not been registered", func() {
		typeDef := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.injected_type")
		registry.Register("goldi.main_type", typeDef)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a direct circular type dependency exists", func() {
		injectedTypeID := "goldi.type_1"
		typeDef1 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.type_2")
		registry.Register(injectedTypeID, typeDef1)

		otherTypeID := "goldi.type_2"
		typeDef2 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.type_1")
		registry.Register(otherTypeID, typeDef2)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should return an error when a transitive circular type dependency exists", func() {
		typeID1 := "goldi.type_1"
		typeDef1 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.type_2")
		registry.Register(typeID1, typeDef1)

		typeID2 := "goldi.type_2"
		typeDef2 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.type_3")
		registry.Register(typeID2, typeDef2)

		typeID3 := "goldi.type_3"
		typeDef3 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.type_1")
		registry.Register(typeID3, typeDef3)

		Expect(validator.Validate(container)).NotTo(Succeed())
	})

	It("should not return an error when everything is OK", func() {
		config["param"] = true
		injectedTypeID := "goldi.injected_type"
		typeDef1 := goldi.NewType(testAPI.NewMockTypeWithArgs, "hello world", "%param%")
		registry.Register(injectedTypeID, typeDef1)

		otherTypeID := "goldi.main_type"
		typeDef2 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.injected_type")
		registry.Register(otherTypeID, typeDef2)

		Expect(validator.Validate(container)).To(Succeed())
	})
})
