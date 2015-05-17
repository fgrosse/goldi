package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("InstanceType", func() {

	var (
		config   map[string]interface{}
		registry goldi.TypeRegistry
	)

	BeforeEach(func() {
		config = map[string]interface{}{}
		registry = goldi.NewTypeRegistry()
	})

	It("should panic if NewInstanceType is called with nil", func() {
		Expect(func() { goldi.NewInstanceType(nil) }).To(Panic())
	})

	Describe("Generate", func() {
		It("should always return the given instance", func() {
			instance := testAPI.NewFoo()
			factory := goldi.NewInstanceType(instance)

			for i := 0; i < 3; i++ {
				generateResult := factory.Generate(config, registry)
				Expect(generateResult == instance).To(BeTrue(),
					fmt.Sprintf("generateResult (%p) should point to the same instance as instance (%p)", generateResult, instance),
				)
			}
		})

		It("should panic if is called with nil", func() {
			factory := &goldi.InstanceType{}
			Expect(func() { factory.Generate(config, registry) }).To(Panic())
		})
	})

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewInstanceType("foo")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})
})
