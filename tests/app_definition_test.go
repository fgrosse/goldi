package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("AppDefinition", func() {

	var definition goldi.AppDefinition

	BeforeEach(func() {
		definition = goldi.NewAppDefinition()
	})

	Describe("RegisterType", func() {
		It("should store the type generator", func() {
			typeID := "goldi.test_type"
			generator := &testAPI.MockTypeGenerator{}
			Expect(definition.RegisterType(typeID, generator.NewMockType)).To(Succeed())

			generatorWrapper, typeIsRegistered := definition[typeID]
			Expect(typeIsRegistered).To(BeTrue())
			Expect(generatorWrapper).NotTo(BeNil())

			generatorWrapper.Generate()
			Expect(generator.HasBeenUsed).To(BeTrue())
		})

		It("should return an error if the type has been defined previously", func() {
			typeID := "goldi.test_type"
			generator := &testAPI.MockTypeGenerator{}
			Expect(definition.RegisterType(typeID, generator.NewMockType)).To(Succeed())
			Expect(definition.RegisterType(typeID, generator.NewMockType)).NotTo(Succeed())
		})
	})
})
