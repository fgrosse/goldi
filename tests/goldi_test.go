package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("Goldi", func() {
	var (
		definition goldi.AppDefinition
		container  *goldi.Container
	)

	BeforeEach(func() {
		definition = goldi.NewAppDefinition()
		container = goldi.NewContainer(definition)
	})

	It("should panic if a type can not be resolved", func() {
		Expect(func() { container.Get("foo.bar") }).To(Panic())
	})

	It("should resolve simple types", func() {
		Expect(definition.RegisterType("goldi.test_type", testAPI.NewMockType)).To(Succeed())
		Expect(container.Get("goldi.test_type")).To(BeAssignableToTypeOf(&testAPI.MockType{}))
	})
})
