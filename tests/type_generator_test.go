package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeGenerator", func() {
	var (
		err       error // just for convenience
		generator *goldi.TypeGenerator
	)

	Describe("Generate", func() {
		Context("with invalid generator function", func() {
			It("should return an error if the generator is no function", func() {
				_, err = goldi.NewTypeGenerator(42)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has no output parameters", func() {
				_, err = goldi.NewTypeGenerator(func() {})
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has more than one output parameter", func() {
				_, err = goldi.NewTypeGenerator(func() (*testAPI.MockType, *testAPI.MockType) { return nil, nil })
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the return parameter is no pointer", func() {
				_, err = goldi.NewTypeGenerator(func() testAPI.MockType { return testAPI.MockType{} })
				Expect(err).To(HaveOccurred())
			})

			It("should not return an error if the return parameter is an interface", func() {
				_, err = goldi.NewTypeGenerator(func() interface{} { return testAPI.MockType{} })
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with generator functions without arguments", func() {
			BeforeEach(func() {
				generator, err = goldi.NewTypeGenerator(testAPI.NewMockType)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when no argument is given in Generate", func() {
				It("should generate the type", func() {
					Expect(generator.Generate()).To(BeAssignableToTypeOf(&testAPI.MockType{}))
				})
			})

			Context("when any argument is given in Generate", func() {
				It("should panic", func() {
					Expect(func() { generator.Generate("foo") }).To(Panic())
				})
			})
		})

		Context("with generator functions with one or more arguments", func() {
			BeforeEach(func() {
				generator, err = goldi.NewTypeGenerator(testAPI.NewMockTypeWithArgs)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when an invalid number of arguments is given", func() {
				It("should panic", func() {
					Expect(func() { generator.Generate() }).To(Panic())
					Expect(func() { generator.Generate("foo") }).To(Panic())
					Expect(func() { generator.Generate("foo", false, 42) }).To(Panic())
				})
			})

			Context("when an the wrong argument types are given", func() {
				It("should panic", func() {
					Expect(func() { generator.Generate("foo", "bar") }).To(Panic())
					Expect(func() { generator.Generate(true, false) }).To(Panic())
					Expect(func() { generator.Generate(true, "bar") }).To(Panic())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should generate the type", func() {
					Expect(generator.Generate("foo", true)).To(BeAssignableToTypeOf(&testAPI.MockType{}))
				})
			})
		})
	})
})
