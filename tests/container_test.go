package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("Container", func() {
	var (
		registry  goldi.TypeRegistry
		config    map[string]interface{}
		container *goldi.Container
	)

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
		config = map[string]interface{}{}
		container = goldi.NewContainer(registry, config)
	})

	It("should panic if a type can not be resolved", func() {
		Expect(func() { container.Get("foo.bar") }).To(Panic())
	})

	It("should resolve simple types", func() {
		registry.RegisterType("goldi.test_type", testAPI.NewMockType)
		Expect(container.Get("goldi.test_type")).To(BeAssignableToTypeOf(&testAPI.MockType{}))
	})

	It("should build the types lazily", func() {
		typeID := "goldi.test_type"
		generator := &testAPI.MockTypeFactory{}
		registry.RegisterType(typeID, generator.NewMockType)

		generatorWrapper, typeIsRegistered := registry[typeID]
		Expect(typeIsRegistered).To(BeTrue())
		Expect(generatorWrapper).NotTo(BeNil())

		Expect(generator.HasBeenUsed).To(BeFalse())
		container.Get(typeID)
		Expect(generator.HasBeenUsed).To(BeTrue())
	})

	It("should build the types as singletons (one instance per type ID)", func() {
		typeID := "goldi.test_type"
		generator := &testAPI.MockTypeFactory{}
		registry.RegisterType(typeID, generator.NewMockType)

		generatorWrapper, typeIsRegistered := registry[typeID]
		Expect(typeIsRegistered).To(BeTrue())
		Expect(generatorWrapper).NotTo(BeNil())

		firstResult := container.Get(typeID)
		secondResult := container.Get(typeID)
		thirdResult := container.Get(typeID)
		Expect(firstResult == secondResult).To(BeTrue())
		Expect(firstResult == thirdResult).To(BeTrue())
	})

	It("should pass static parameters as arguments when generating types", func() {
		typeID := "goldi.test_type"
		typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, "parameter1", true)
		registry.Register(typeID, typeDef)

		generatedType := container.Get("goldi.test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

		generatedMock := generatedType.(*testAPI.MockType)
		Expect(generatedMock.StringParameter).To(Equal("parameter1"))
		Expect(generatedMock.BoolParameter).To(Equal(true))
	})

	It("should be able to use parameters as arguments when generating types", func() {
		typeID := "goldi.test_type"
		typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, "%parameter1%", "%parameter2%")
		registry.Register(typeID, typeDef)

		config["parameter1"] = "test"
		config["parameter2"] = true

		generatedType := container.Get("goldi.test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

		generatedMock := generatedType.(*testAPI.MockType)
		Expect(generatedMock.StringParameter).To(Equal(config["parameter1"]))
		Expect(generatedMock.BoolParameter).To(Equal(config["parameter2"]))
	})

	It("should be able to inject already defined types into other types", func() {
		injectedTypeID := "goldi.injected_type"
		typeDef1 := goldi.NewType(testAPI.NewMockType)
		registry.Register(injectedTypeID, typeDef1)

		otherTypeID := "goldi.main_type"
		typeDef2 := goldi.NewType(testAPI.NewTypeForServiceInjection, "@goldi.injected_type")
		registry.Register(otherTypeID, typeDef2)

		generatedType := container.Get("goldi.main_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.TypeForServiceInjection{}))

		generatedMock := generatedType.(*testAPI.TypeForServiceInjection)
		Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))
	})
})
