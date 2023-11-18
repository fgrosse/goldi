package goldi_test

import (
	"fmt"

	"github.com/fgrosse/goldi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleNewStructType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	// all of the following types are semantically identical
	container.Register("foo_1", goldi.NewStructType(Foo{}))
	container.Register("foo_2", goldi.NewStructType(&Foo{}))
	container.Register("foo_3", goldi.NewStructType(new(Foo)))

	// each reference to the "logger" type will now be resolved to that instance
	fmt.Printf("foo_1: %T\n", container.MustGet("foo_1"))
	fmt.Printf("foo_2: %T\n", container.MustGet("foo_2"))
	fmt.Printf("foo_3: %T\n", container.MustGet("foo_3"))
	// Output:
	// foo_1: *goldi_test.Foo
	// foo_2: *goldi_test.Foo
	// foo_3: *goldi_test.Foo
}

// ExampleNewStructType_ prevents godoc from printing the whole content of this file as example
func ExampleNewStructType_() {}

var _ = Describe("structType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewStructType(MockType{})
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("goldi.NewStructType()", func() {
		Context("with invalid arguments", func() {
			It("should return an invalid type if the generator is no struct or pointer to a struct", func() {
				Expect(goldi.IsValid(goldi.NewStructType(42))).To(BeFalse())
			})

			It("should return an invalid type if the generator is a pointer to something other than a struct", func() {
				something := "Hello Pointer World!"
				Expect(goldi.IsValid(goldi.NewStructType(&something))).To(BeFalse())
			})
		})

		Context("with first argument beeing a struct", func() {
			It("should create the type", func() {
				typeDef := goldi.NewStructType(MockType{})
				Expect(typeDef).NotTo(BeNil())
			})
		})

		Context("with first argument beeing a pointer to struct", func() {
			It("should create the type", func() {
				typeDef := goldi.NewStructType(&MockType{})
				Expect(typeDef).NotTo(BeNil())
			})
		})

		It("should return an invalid type if more factory arguments were provided than the struct has fields", func() {
			t := goldi.NewStructType(&MockType{}, "foo", true, "bar")
			Expect(goldi.IsValid(t)).To(BeFalse())
			Expect(t).To(MatchError("the struct MockType has only 2 fields but 3 arguments where provided"))
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := goldi.NewStructType(MockType{}, args...)
			Expect(typeDef.Arguments()).To(Equal(args))
		})
	})

	Describe("Generate()", func() {
		var (
			config    = map[string]interface{}{}
			container *goldi.Container
			resolver  *goldi.ParameterResolver
		)

		BeforeEach(func() {
			container = goldi.NewContainer(goldi.NewTypeRegistry(), config)
			resolver = goldi.NewParameterResolver(container)
		})

		Context("without struct arguments", func() {
			Context("when the factory is a struct (no pointer)", func() {
				It("should generate the type", func() {
					typeDef := goldi.NewStructType(MockType{})
					Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&MockType{}))
				})
			})

			It("should generate the type", func() {
				typeDef := goldi.NewStructType(&MockType{})
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&MockType{}))
			})

			It("should generate a new type each time", func() {
				typeDef := goldi.NewStructType(&MockType{})
				t1, err1 := typeDef.Generate(resolver)
				t2, err2 := typeDef.Generate(resolver)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				Expect(t1).NotTo(BeNil())
				Expect(t2).NotTo(BeNil())
				Expect(t1 == t2).To(BeFalse(), fmt.Sprintf("t1 (%p) should not point to the same instance as t2 (%p)", t1, t2))

				// Just to make the whole issue more explicit:
				t1Mock := t1.(*MockType)
				t2Mock := t2.(*MockType)
				t1Mock.StringParameter = "CHANGED"
				Expect(t2Mock.StringParameter).NotTo(Equal(t1Mock.StringParameter),
					"Changing two indipendently generated types should not affect both at the same time",
				)
			})
		})

		Context("with one or more arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewStructType(&MockType{}, "foo", true)

				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

				generatedMock := generatedType.(*MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			It("should use the given parameters", func() {
				typeDef := goldi.NewStructType(&MockType{}, "%param1%", "%param2%")
				config["param1"] = "TEST"
				config["param2"] = true
				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

				generatedMock := generatedType.(*MockType)
				Expect(generatedMock.StringParameter).To(Equal("TEST"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the struct field type", func() {
					It("should generate the type", func() {
						container.RegisterType("foo", NewMockType)
						typeDef := goldi.NewStructType(TypeForServiceInjection{}, "@foo")

						generatedType, err := typeDef.Generate(resolver)
						Expect(err).NotTo(HaveOccurred())
						Expect(generatedType).To(BeAssignableToTypeOf(&TypeForServiceInjection{}))

						generatedMock := generatedType.(*TypeForServiceInjection)
						Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&MockType{}))
					})
				})

				Context("and its type does not match the function signature", func() {
					It("should return an error", func() {
						container.RegisterType("foo", NewFoo)
						typeDef := goldi.NewStructType(TypeForServiceInjection{}, "@foo")

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *goldi_test.Foo) can not be used as field 1 for struct type goldi_test.TypeForServiceInjection`))
					})
				})
			})
		})
	})
})
