package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("Container", func() {
	var (
		definition goldi.AppDefinition
		config     map[string]interface{}
		container  *goldi.Container
	)

	BeforeEach(func() {
		definition = goldi.NewAppDefinition()
		config = map[string]interface{}{}
		container = goldi.NewContainer(definition, config)
	})

	It("should panic if a type can not be resolved", func() {
		Expect(func() { container.Get("foo.bar") }).To(Panic())
	})

	It("should resolve simple types", func() {
		Expect(definition.RegisterType("goldi.test_type", testAPI.NewMockType)).To(Succeed())
		Expect(container.Get("goldi.test_type")).To(BeAssignableToTypeOf(&testAPI.MockType{}))
	})

	It("should be able to use parameters as arguments when generating types", func() {
		typeID := "goldi.test_type"
		typeDef, err := goldi.NewType(testAPI.NewMockTypeWithArgs, "static parameter", "%dynamic_parameter%")
		Expect(err).NotTo(HaveOccurred())
		Expect(definition.Register(typeID, typeDef)).To(Succeed())

		config["dynamic_parameter"] = true

		generatedType := container.Get("goldi.test_type")
		Expect(generatedType).NotTo(BeNil())
		Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

		generatedMock := generatedType.(*testAPI.MockType)
		Expect(generatedMock.StringParameter).To(Equal("static parameter"))
		Expect(generatedMock.BoolParameter).To(Equal(config["dynamic_parameter"]))
	})
})
