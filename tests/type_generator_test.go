package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeGenerator", func() {
	Describe("Generate", func() {

		Context("with invalid generator function", func() {
			It("should return an error if the generator is no function", func() {
				_, err := goldi.NewTypeGenerator(42)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has no output parameters", func() {
				_, err := goldi.NewTypeGenerator(func() {})
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has more than one output parameter", func() {
				_, err := goldi.NewTypeGenerator(func() (*testAPI.MockType, *testAPI.MockType) { return nil, nil })
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the return parameter is no pointer", func() {
				_, err := goldi.NewTypeGenerator(func() testAPI.MockType { return testAPI.MockType{} })
				Expect(err).To(HaveOccurred())
			})

			It("should not return an error if the return parameter is an interface", func() {
				_, err := goldi.NewTypeGenerator(func() interface{} { return testAPI.MockType{} })
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("should generate simple types from a generator function", func() {
			generator, err := goldi.NewTypeGenerator(testAPI.NewMockType)
			Expect(err).NotTo(HaveOccurred())
			Expect(generator.Generate()).To(BeAssignableToTypeOf(&testAPI.MockType{}))
		})
	})
})
