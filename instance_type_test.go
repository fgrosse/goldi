package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("InstanceType", func() {
	var resolver *ParameterResolver

	BeforeEach(func() {
		container := NewContainer(NewTypeRegistry(), map[string]interface{}{})
		resolver = NewParameterResolver(container)
	})

	It("should return an invalid type if NewInstanceType is called with nil", func() {
		Expect(IsValid(NewInstanceType(nil))).To(BeFalse())
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := NewInstanceType(tests.NewFoo())
			Expect(typeDef.Arguments()).To(BeEmpty())
		})
	})

	Describe("Generate", func() {
		It("should always return the given instance", func() {
			instance := tests.NewFoo()
			factory := NewInstanceType(instance)

			for i := 0; i < 3; i++ {
				generateResult, err := factory.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generateResult == instance).To(BeTrue(),
					fmt.Sprintf("generateResult (%p) should point to the same instance as instance (%p)", generateResult, instance),
				)
			}
		})
	})

	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = NewInstanceType("foo")
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})
})
