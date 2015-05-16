package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeInstanceFactory", func() {

	var (
		config   map[string]interface{}
		registry goldi.TypeRegistry
	)

	BeforeEach(func() {
		config = map[string]interface{}{}
		registry = goldi.NewTypeRegistry()
	})

	It("should panic if NewTypeInstanceFactory is called with nil", func() {
		Expect(func() { goldi.NewTypeInstanceFactory(nil) }).To(Panic())
	})

	Describe("Generate", func() {
		It("should always return the given instance", func() {
			instance := testAPI.NewFoo()
			factory := goldi.NewTypeInstanceFactory(instance)

			for i := 0; i < 3; i++ {
				generateResult := factory.Generate(config, registry)
				Expect(generateResult == instance).To(BeTrue(),
					fmt.Sprintf("generateResult (%p) should point to the same instance as instance (%p)", generateResult, instance),
				)
			}
		})

		It("should panic if is called with nil", func() {
			factory := &goldi.TypeInstanceFactory{}
			Expect(func() { factory.Generate(config, registry) }).To(Panic())
		})
	})
})
