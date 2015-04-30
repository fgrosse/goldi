package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeRegistry", func() {

	var registry goldi.TypeRegistry

	BeforeEach(func() {
		registry = goldi.NewTypeRegistry()
	})

	Describe("RegisterType", func() {
		It("should store the type generator", func() {
			typeID := "goldi.test_type"
			generator := &testAPI.MockTypeFactory{}
			Expect(registry.RegisterType(typeID, generator.NewMockType)).To(Succeed())

			generatorWrapper, typeIsRegistered := registry[typeID]
			Expect(typeIsRegistered).To(BeTrue())
			Expect(generatorWrapper).NotTo(BeNil())

			config := map[string]interface{}{}
			generatorWrapper.Generate(config)
			Expect(generator.HasBeenUsed).To(BeTrue())
		})

		It("should recover panics from NewType", func() {
			Expect(func() { registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs) }).NotTo(Panic())
			Expect(registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs)).NotTo(Succeed())
		})

		It("should return an error if the type has been defined previously", func() {
			typeID := "goldi.test_type"
			generator := &testAPI.MockTypeFactory{}
			Expect(registry.RegisterType(typeID, generator.NewMockType)).To(Succeed())
			Expect(registry.RegisterType(typeID, generator.NewMockType)).NotTo(Succeed())
		})
	})
})
