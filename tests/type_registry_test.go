package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeRegistry", func() {
	var (
		registry goldi.TypeRegistry
		config   = map[string]interface{}{} // for test convenience
		resolver *goldi.ParameterResolver
	)

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
		resolver = goldi.NewParameterResolver(config, registry)
	})

	Describe("RegisterType", func() {
		Context("with factory function type", func() {
			It("should store the type", func() {
				typeID := "goldi.test_type"
				factory := &testAPI.MockTypeFactory{}
				registry.RegisterType(typeID, factory.NewMockType)

				factoryWrapper, typeIsRegistered := registry[typeID]
				Expect(typeIsRegistered).To(BeTrue())
				Expect(factoryWrapper).NotTo(BeNil())

				factoryWrapper.Generate(resolver)
				Expect(factory.HasBeenUsed).To(BeTrue())
			})

			It("should pass parameters to the new type", func() {
				typeID := "goldi.test_type"
				registry.RegisterType(typeID, testAPI.NewMockTypeWithArgs, "foo", true)
				Expect(registry).To(HaveKey(typeID))
				Expect(registry["goldi.test_type"].Generate(resolver).(*testAPI.MockType).StringParameter).To(Equal("foo"))
				Expect(registry["goldi.test_type"].Generate(resolver).(*testAPI.MockType).BoolParameter).To(Equal(true))
			})
		})

		Context("with struct type", func() {
			It("should store the type", func() {
				typeID := "goldi.test_type"
				foo := testAPI.Foo{}
				registry.RegisterType(typeID, foo)

				fooType, typeIsRegistered := registry[typeID]
				Expect(typeIsRegistered).To(BeTrue())
				Expect(fooType).NotTo(BeNil())

				newFoo := fooType.Generate(resolver)
				Expect(newFoo).To(BeAssignableToTypeOf(&foo))
			})

			It("should pass parameters to the new type", func() {
				typeID := "goldi.test_type"
				registry.RegisterType(typeID, testAPI.Baz{}, "param1", "param2")
				Expect(registry).To(HaveKey(typeID))

				newBaz := registry["goldi.test_type"].Generate(resolver).(*testAPI.Baz)
				Expect(newBaz.Parameter1).To(Equal("param1"))
				Expect(newBaz.Parameter2).To(Equal("param2"))
			})
		})

		Context("with pointer to struct type", func() {
			It("should store the type", func() {
				typeID := "goldi.test_type"
				foo := &testAPI.Foo{}
				registry.RegisterType(typeID, foo)

				fooType, typeIsRegistered := registry[typeID]
				Expect(typeIsRegistered).To(BeTrue())
				Expect(fooType).NotTo(BeNil())

				newFoo := fooType.Generate(resolver)
				Expect(newFoo).To(BeAssignableToTypeOf(foo))
			})

			It("should pass parameters to the new type", func() {
				typeID := "goldi.test_type"
				registry.RegisterType(typeID, &testAPI.Baz{}, "param1", "param2")
				Expect(registry).To(HaveKey(typeID))

				newBaz := registry["goldi.test_type"].Generate(resolver).(*testAPI.Baz)
				Expect(newBaz.Parameter1).To(Equal("param1"))
				Expect(newBaz.Parameter2).To(Equal("param2"))
			})
		})

		Context("with invalid factory type", func() {
			It("should panic", func() {
				Expect(func() { registry.RegisterType("invalid_type", 42) }).To(Panic())
			})
		})
	})

	Describe("InjectInstance", func() {
		It("should store the type instance", func() {
			typeID := "goldi.test_type"
			fooInstance := testAPI.NewFoo()
			Expect(registry.InjectInstance(typeID, fooInstance)).To(Succeed())

			factory, typeIsRegistered := registry[typeID]
			Expect(typeIsRegistered).To(BeTrue())
			Expect(factory).NotTo(BeNil())

			generateResult := factory.Generate(resolver)
			Expect(generateResult == fooInstance).To(BeTrue(),
				fmt.Sprintf("generateResult (%p) should point to the same instance as fooInstance (%p)", generateResult, fooInstance),
			)
		})

		It("should recover panics from NewInstanceType", func() {
			Expect(registry.InjectInstance("goldi.test_type", nil)).NotTo(Succeed())
		})
	})
})
