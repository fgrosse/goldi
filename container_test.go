package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("Container", func() {
	var (
		registry  TypeRegistry
		config    map[string]interface{}
		container *Container
	)

	BeforeEach(func() {
		registry = NewTypeRegistry()
		config = map[string]interface{}{}
		container = NewContainer(registry, config)
	})

	It("should panic if a type can not be resolved", func() {
		Expect(func() { container.MustGet("foo.bar") }).To(Panic())
	})

	It("should resolve simple types", func() {
		registry.RegisterType("test_type", tests.NewMockType)
		Expect(container.MustGet("test_type")).To(BeAssignableToTypeOf(&tests.MockType{}))
	})

	It("should build the types lazily", func() {
		typeID := "test_type"
		generator := &tests.MockTypeFactory{}
		registry.RegisterType(typeID, generator.NewMockType)

		generatorWrapper, typeIsRegistered := registry[typeID]
		Expect(typeIsRegistered).To(BeTrue())
		Expect(generatorWrapper).NotTo(BeNil())

		Expect(generator.HasBeenUsed).To(BeFalse())
		container.MustGet(typeID)
		Expect(generator.HasBeenUsed).To(BeTrue())
	})

	It("should build the types as singletons (one instance per type ID)", func() {
		typeID := "test_type"
		generator := &tests.MockTypeFactory{}
		registry.RegisterType(typeID, generator.NewMockType)

		generatorWrapper, typeIsRegistered := registry[typeID]
		Expect(typeIsRegistered).To(BeTrue())
		Expect(generatorWrapper).NotTo(BeNil())

		firstResult := container.MustGet(typeID)
		secondResult := container.MustGet(typeID)
		thirdResult := container.MustGet(typeID)
		Expect(firstResult == secondResult).To(BeTrue())
		Expect(firstResult == thirdResult).To(BeTrue())
	})

	It("should pass static parameters as arguments when generating types", func() {
		typeID := "test_type"
		typeDef := NewType(tests.NewMockTypeWithArgs, "parameter1", true)
		registry.Register(typeID, typeDef)

		generatedType := container.MustGet("test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

		generatedMock := generatedType.(*tests.MockType)
		Expect(generatedMock.StringParameter).To(Equal("parameter1"))
		Expect(generatedMock.BoolParameter).To(Equal(true))
	})

	It("should be able to use parameters as arguments when generating types", func() {
		typeID := "test_type"
		typeDef := NewType(tests.NewMockTypeWithArgs, "%parameter1%", "%parameter2%")
		registry.Register(typeID, typeDef)

		config["parameter1"] = "test"
		config["parameter2"] = true

		generatedType := container.MustGet("test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

		generatedMock := generatedType.(*tests.MockType)
		Expect(generatedMock.StringParameter).To(Equal(config["parameter1"]))
		Expect(generatedMock.BoolParameter).To(Equal(config["parameter2"]))
	})

	It("should be able to inject already defined types into other types", func() {
		registry.Register("injected_type", NewType(tests.NewMockType))
		registry.Register("main_type", NewType(tests.NewTypeForServiceInjection, "@injected_type"))

		generatedType := container.MustGet("main_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&tests.TypeForServiceInjection{}))

		generatedMock := generatedType.(*tests.TypeForServiceInjection)
		Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&tests.MockType{}))
	})

	It("should inject the same instance when it is used by different services", func() {
		registry.RegisterType("foo", tests.NewMockType)
		registry.RegisterType("type1", tests.NewTypeForServiceInjection, "@foo")
		registry.RegisterType("type2", tests.NewTypeForServiceInjection, "@foo")

		generatedType1 := container.MustGet("type1")
		generatedType2 := container.MustGet("type2")
		Expect(generatedType1).To(BeAssignableToTypeOf(&tests.TypeForServiceInjection{}))
		Expect(generatedType2).To(BeAssignableToTypeOf(&tests.TypeForServiceInjection{}))

		generatedMock1 := generatedType1.(*tests.TypeForServiceInjection)
		generatedMock2 := generatedType2.(*tests.TypeForServiceInjection)

		Expect(generatedMock1.InjectedType == generatedMock2.InjectedType).To(BeTrue(), "Both generated types should have the same instance of @foo")
	})

	It("should inject nil when using optional types that are not defined", func() {
		registry.Register("main_type", NewType(tests.NewTypeForServiceInjection, "@?optional_type"))

		generatedType := container.MustGet("main_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&tests.TypeForServiceInjection{}))

		generatedMock := generatedType.(*tests.TypeForServiceInjection)
		Expect(generatedMock.InjectedType).To(BeNil())
	})
})
