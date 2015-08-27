package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("invalidType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = newInvalidType(fmt.Errorf("Something bad happened"))
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := newInvalidType(fmt.Errorf("Something bad happened"))
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(BeEmpty())
		})
	})

	Describe("Generate()", func() {
		It("should return its error", func() {
			config := map[string]interface{}{}
			container := NewContainer(NewTypeRegistry(), config)
			resolver := NewParameterResolver(container)

			e := fmt.Errorf("Something bad happened")
			typeDef := newInvalidType(e)

			generated, err := typeDef.Generate(resolver)
			Expect(generated).To(BeNil())
			Expect(err).To(Equal(e))
		})
	})
})

var _ = Describe("IsValid", func() {
	It("should return false when given an instance of invalidType", func() {
		Expect(IsValid(&invalidType{})).To(BeFalse())
	})

	It("should return true when not given an instance of invalidType", func() {
		Expect(IsValid(NewAliasType("foo"))).To(BeTrue())
	})
})
