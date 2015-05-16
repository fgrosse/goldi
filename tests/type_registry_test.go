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
	)

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
	})

	Describe("RegisterType", func() {
		It("should store the type", func() {
			typeID := "goldi.test_type"
			factory := &testAPI.MockTypeFactory{}
			Expect(registry.RegisterType(typeID, factory.NewMockType)).To(Succeed())

			factoryWrapper, typeIsRegistered := registry[typeID]
			Expect(typeIsRegistered).To(BeTrue())
			Expect(factoryWrapper).NotTo(BeNil())

			factoryWrapper.Generate(config, registry)
			Expect(factory.HasBeenUsed).To(BeTrue())
		})

		It("should recover panics from NewType", func() {
			Expect(func() { registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs) }).NotTo(Panic())
			Expect(registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs)).NotTo(Succeed())
		})

		It("should pass parameters to the new type", func() {
			typeID := "goldi.test_type"
			Expect(registry.RegisterType(typeID, testAPI.NewMockTypeWithArgs, "foo", true)).To(Succeed())
			Expect(registry).To(HaveKey(typeID))
			Expect(registry["goldi.test_type"].Generate(config, registry).(*testAPI.MockType).StringParameter).To(Equal("foo"))
			Expect(registry["goldi.test_type"].Generate(config, registry).(*testAPI.MockType).BoolParameter).To(Equal(true))
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

			generateResult := factory.Generate(config, registry)
			Expect(generateResult == fooInstance).To(BeTrue(),
				fmt.Sprintf("generateResult (%p) should point to the same instance as fooInstance (%p)", generateResult, fooInstance),
			)
		})

		It("should recover panics from NewTypeInstanceFactory", func() {
			Expect(registry.InjectInstance("goldi.test_type", nil)).NotTo(Succeed())
		})
	})
})
