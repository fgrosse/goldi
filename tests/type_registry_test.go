package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
			generator := &testAPI.MockTypeFactory{}
			Expect(registry.RegisterType(typeID, generator.NewMockType)).To(Succeed())

			generatorWrapper, typeIsRegistered := registry[typeID]
			Expect(typeIsRegistered).To(BeTrue())
			Expect(generatorWrapper).NotTo(BeNil())

			generatorWrapper.Generate(config)
			Expect(generator.HasBeenUsed).To(BeTrue())
		})

		It("should recover panics from NewType", func() {
			Expect(func() { registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs) }).NotTo(Panic())
			Expect(registry.RegisterType("goldi.test_type", testAPI.NewMockTypeWithArgs)).NotTo(Succeed())
		})

		It("should pass parameters to the new type", func() {
			typeID := "goldi.test_type"
			Expect(registry.RegisterType(typeID, testAPI.NewMockTypeWithArgs, "foo", true)).To(Succeed())
			Expect(registry).To(HaveKey(typeID))
			Expect(registry["goldi.test_type"].Generate(config).(*testAPI.MockType).StringParameter).To(Equal("foo"))
			Expect(registry["goldi.test_type"].Generate(config).(*testAPI.MockType).BoolParameter).To(Equal(true))
		})
	})
})
