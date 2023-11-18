package goldi_test

import (
	"fmt"

	"github.com/fgrosse/goldi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExampleNewType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	// register the type using the factory function NewMockTypeWithArgs and pass two arguments
	container.Register("my_type", goldi.NewType(NewMockTypeWithArgs, "Hello World", true))

	t := container.MustGet("my_type").(*MockType)
	fmt.Printf("%#v", t)
	// Output:
	// &goldi_test.MockType{StringParameter:"Hello World", BoolParameter:true}
}

// ExampleNewType_ prevents godoc from printing the whole content of this file as example
func ExampleNewType_() {}

var _ = Describe("type", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewType(NewFoo)
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("goldi.NewType()", func() {
		Context("with invalid factory function", func() {
			It("should return an invalid type if the generator is no function", func() {
				Expect(goldi.IsValid(goldi.NewType(42))).To(BeFalse())
			})

			It("should return an invalid type if the generator has no output parameters", func() {
				Expect(goldi.IsValid(goldi.NewType(func() {}))).To(BeFalse())
			})

			It("should return an invalid type if the generator has more than one output parameter", func() {
				Expect(goldi.IsValid(goldi.NewType(func() (*MockType, *MockType) { return nil, nil }))).To(BeFalse())
			})

			It("should return an invalid type if the return parameter is no pointer", func() {
				Expect(goldi.IsValid(goldi.NewType(func() MockType { return MockType{} }))).To(BeFalse())
			})

			It("should not return an invalid type if the return parameter is an interface", func() {
				Expect(goldi.IsValid(goldi.NewType(func() interface{} { return MockType{} }))).To(BeTrue())
			})

			It("should not return an invalid type if the return parameter is a function", func() {
				Expect(goldi.IsValid(goldi.NewType(func() func() { return func() {} }))).To(BeTrue())
			})
		})

		Context("without factory function arguments", func() {
			Context("when no factory argument is given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(NewMockType)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when any argument is given", func() {
				It("should return an invalid type", func() {
					Expect(goldi.IsValid(goldi.NewType(NewMockType, "foo"))).To(BeFalse())
				})
			})
		})

		Context("with one or more factory function arguments", func() {
			Context("when an invalid number of arguments is given", func() {
				It("should return an invalid type", func() {
					Expect(goldi.IsValid(goldi.NewType(NewMockTypeWithArgs))).To(BeFalse())
					Expect(goldi.IsValid(goldi.NewType(NewMockTypeWithArgs, "foo"))).To(BeFalse())
					Expect(goldi.IsValid(goldi.NewType(NewMockTypeWithArgs, "foo", false, 42))).To(BeFalse())
				})
			})

			Context("when the wrong argument types are given", func() {
				It("should return an invalid type", func() {
					Expect(goldi.IsValid(goldi.NewType(NewMockTypeWithArgs, "foo", "bar"))).To(BeFalse())
					Expect(goldi.IsValid(goldi.NewType(NewMockTypeWithArgs, true, "bar"))).To(BeFalse())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(NewMockTypeWithArgs, "foo", true)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when the arguments are variadic", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())
				})

				It("should return an invalid type if not enough arguments where given", func() {
					t := goldi.NewType(NewVariadicMockType, true)
					Expect(goldi.IsValid(t)).To(BeFalse())
					Expect(t).To(MatchError("invalid number of input parameters for variadic function: got 1 but expected at least 3"))
				})
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := goldi.NewType(NewMockTypeWithArgs, args...)
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

		Context("without factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewType(NewMockType)
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&MockType{}))
			})
		})

		Context("with one or more factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewType(NewMockTypeWithArgs, "foo", true)

				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

				generatedMock := generatedType.(*MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the function signature", func() {
					It("should generate the type", func() {
						container.RegisterType("foo", NewMockType)
						typeDef := goldi.NewType(NewTypeForServiceInjection, "@foo")

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
						typeDef := goldi.NewType(NewTypeForServiceInjectionWithArgs, "@foo", "arg1", "arg2", true)

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *goldi_test.Foo) can not be passed as argument 1 to the function signature goldi_test.NewTypeForServiceInjectionWithArgs(*goldi_test.MockType, string, string, bool)`))
					})
				})
			})

			Context("when the arguments are variadic", func() {
				It("should generate the type", func() {
					typeDef := goldi.NewType(NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(&MockType{}))

					generatedMock := generatedType.(*MockType)
					Expect(generatedMock.BoolParameter).To(BeTrue())
					Expect(generatedMock.StringParameter).To(Equal("1, two, drei"))
				})
			})

			Context("when a func reference type is given", func() {
				It("should generate the type", func() {
					foo := &MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := goldi.NewType(NewMockTypeFromStringFunc, "YEAH", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(NewMockType()))
					Expect(generatedType.(*MockType).StringParameter).To(Equal("Success! YEAH"))
				})
			})

			Context("when a func reference type is given as variadic argument", func() {
				It("should generate the type", func() {
					foo := &MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := goldi.NewType(NewVariadicMockTypeFuncs, "@foo::ReturnString", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(NewMockType()))
					Expect(generatedType.(*MockType).StringParameter).To(Equal("Success! Success! "))
				})
			})
		})
	})
})
