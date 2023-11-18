package goldi_test

import (
	"fmt"

	"github.com/fgrosse/goldi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleContainer() {
	registry := goldi.NewTypeRegistry()
	config := map[string]interface{}{}
	container := goldi.NewContainer(registry, config)

	container.Register("logger", goldi.NewType(NewNullLogger))

	l := container.MustGet("logger")
	fmt.Printf("%T", l)
	// Output:
	// *goldi_test.NullLogger
}

func ExampleContainer_Get() {
	registry := goldi.NewTypeRegistry()
	config := map[string]interface{}{}
	container := goldi.NewContainer(registry, config)

	container.Register("logger", goldi.NewType(NewNullLogger))

	l, err := container.Get("logger")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// do stuff with the logger. usually you need a type assertion
	fmt.Printf("%T", l.(*NullLogger))

	// Output:
	// *goldi_test.NullLogger
}

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
		Expect(func() { container.MustGet("foo.bar") }).To(Panic())
	})

	It("should return an error if there was an issue generating the type", func() {
		container.Register("foo", goldi.NewStructType(nil))
		_, err := container.Get("foo")
		Expect(err).To(MatchError(`goldi: error while generating type "foo": the given struct is nil`))
	})

	It("should resolve simple types", func() {
		registry.RegisterType("test_type", NewMockType)
		Expect(container.MustGet("test_type")).To(BeAssignableToTypeOf(&MockType{}))
	})

	It("should build the types lazily", func() {
		typeID := "test_type"
		generator := &MockTypeFactory{}
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
		generator := &MockTypeFactory{}
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
		typeDef := goldi.NewType(NewMockTypeWithArgs, "parameter1", true)
		registry.Register(typeID, typeDef)

		generatedType := container.MustGet("test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

		generatedMock := generatedType.(*MockType)
		Expect(generatedMock.StringParameter).To(Equal("parameter1"))
		Expect(generatedMock.BoolParameter).To(Equal(true))
	})

	It("should be able to use parameters as arguments when generating types", func() {
		typeID := "test_type"
		typeDef := goldi.NewType(NewMockTypeWithArgs, "%parameter1%", "%parameter2%")
		registry.Register(typeID, typeDef)

		config["parameter1"] = "test"
		config["parameter2"] = true

		generatedType := container.MustGet("test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

		generatedMock := generatedType.(*MockType)
		Expect(generatedMock.StringParameter).To(Equal(config["parameter1"]))
		Expect(generatedMock.BoolParameter).To(Equal(config["parameter2"]))
	})

	It("should be able to inject already defined types into other types", func() {
		registry.Register("injected_type", goldi.NewType(NewMockType))
		registry.Register("main_type", goldi.NewType(NewTypeForServiceInjection, "@injected_type"))

		generatedType := container.MustGet("main_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&TypeForServiceInjection{}))

		generatedMock := generatedType.(*TypeForServiceInjection)
		Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&MockType{}))
	})

	It("should inject the same instance when it is used by different services", func() {
		registry.RegisterType("foo", NewMockType)
		registry.RegisterType("type1", NewTypeForServiceInjection, "@foo")
		registry.RegisterType("type2", NewTypeForServiceInjection, "@foo")

		generatedType1 := container.MustGet("type1")
		generatedType2 := container.MustGet("type2")
		Expect(generatedType1).To(BeAssignableToTypeOf(&TypeForServiceInjection{}))
		Expect(generatedType2).To(BeAssignableToTypeOf(&TypeForServiceInjection{}))

		generatedMock1 := generatedType1.(*TypeForServiceInjection)
		generatedMock2 := generatedType2.(*TypeForServiceInjection)

		Expect(generatedMock1.InjectedType == generatedMock2.InjectedType).To(BeTrue(), "Both generated types should have the same instance of @foo")
	})

	It("should inject nil when using optional types that are not defined", func() {
		registry.Register("main_type", goldi.NewType(NewTypeForServiceInjection, "@?optional_type"))

		generatedType := container.MustGet("main_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&TypeForServiceInjection{}))

		generatedMock := generatedType.(*TypeForServiceInjection)
		Expect(generatedMock.InjectedType).To(BeNil())
	})
})
