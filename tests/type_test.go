package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("Type", func() {
	var (
		err     error // just for convenience
		typeDef *goldi.Type
	)

	Describe("NewType()", func() {
		Context("with invalid factory function", func() {
			It("should return an error if the generator is no function", func() {
				_, err = goldi.NewType(42)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has no output parameters", func() {
				_, err = goldi.NewType(func() {})
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the generator has more than one output parameter", func() {
				_, err = goldi.NewType(func() (*testAPI.MockType, *testAPI.MockType) { return nil, nil })
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if the return parameter is no pointer", func() {
				_, err = goldi.NewType(func() testAPI.MockType { return testAPI.MockType{} })
				Expect(err).To(HaveOccurred())
			})

			It("should not return an error if the return parameter is an interface", func() {
				_, err = goldi.NewType(func() interface{} { return testAPI.MockType{} })
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with factory functions without arguments", func() {
			Context("when no factory argument is given", func() {
				It("should create the type", func() {
					typeDef, err = goldi.NewType(testAPI.NewMockType)
					Expect(err).NotTo(HaveOccurred())
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when any argument is given", func() {
				It("should return an error", func() {
					typeDef, err = goldi.NewType(testAPI.NewMockType, "foo")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("with factory functions with one or more arguments", func() {
			Context("when an invalid number of arguments is given", func() {
				It("should return an error", func() {
					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs)
					Expect(err).To(HaveOccurred())

					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, "foo")
					Expect(err).To(HaveOccurred())

					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", false, 42)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when the wrong argument types are given", func() {
				It("should return an error", func() {
					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", "bar")
					Expect(err).To(HaveOccurred())

					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, true, "bar")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should create the type", func() {
					typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", true)
					Expect(err).NotTo(HaveOccurred())
					Expect(typeDef).NotTo(BeNil())
				})
			})
		})
	})

	Describe("Generate()", func() {
		config := map[string]interface{}{}

		Context("with factory functions without arguments", func() {
			It("should generate the type", func() {
				typeDef, err = goldi.NewType(testAPI.NewMockType)
				Expect(typeDef.Generate(config)).To(BeAssignableToTypeOf(&testAPI.MockType{}))
			})
		})

		Context("with factory functions with one or more arguments", func() {
			It("should generate the type", func() {
				typeDef, err = goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", true)
				Expect(err).NotTo(HaveOccurred())

				generatedType := typeDef.Generate(config)
				Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

				generatedMock := generatedType.(*testAPI.MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})
		})
	})
})
