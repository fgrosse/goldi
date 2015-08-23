package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("Type", func() {
	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = NewType(tests.NewFoo)
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewType()", func() {
		Context("with invalid factory function", func() {
			It("should return an invalid type if the generator is no function", func() {
				Expect(IsValid(NewType(42))).To(BeFalse())
			})

			It("should return an invalid type if the generator has no output parameters", func() {
				Expect(IsValid(NewType(func() {}))).To(BeFalse())
			})

			It("should return an invalid type if the generator has more than one output parameter", func() {
				Expect(IsValid(NewType(func() (*tests.MockType, *tests.MockType) { return nil, nil }))).To(BeFalse())
			})

			It("should return an invalid type if the return parameter is no pointer", func() {
				Expect(IsValid(NewType(func() tests.MockType { return tests.MockType{} }))).To(BeFalse())
			})

			It("should not return an invalid type if the return parameter is an interface", func() {
				Expect(IsValid(NewType(func() interface{} { return tests.MockType{} }))).To(BeTrue())
			})
		})

		Context("without factory function arguments", func() {
			Context("when no factory argument is given", func() {
				It("should create the type", func() {
					typeDef := NewType(tests.NewMockType)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when any argument is given", func() {
				It("should return an invalid type", func() {
					Expect(IsValid(NewType(tests.NewMockType, "foo"))).To(BeFalse())
				})
			})
		})

		Context("with one or more factory function arguments", func() {
			Context("when an invalid number of arguments is given", func() {
				It("should return an invalid type", func() {
					Expect(IsValid(NewType(tests.NewMockTypeWithArgs))).To(BeFalse())
					Expect(IsValid(NewType(tests.NewMockTypeWithArgs, "foo"))).To(BeFalse())
					Expect(IsValid(NewType(tests.NewMockTypeWithArgs, "foo", false, 42))).To(BeFalse())
				})
			})

			Context("when the wrong argument types are given", func() {
				It("should return an invalid type", func() {
					Expect(IsValid(NewType(tests.NewMockTypeWithArgs, "foo", "bar"))).To(BeFalse())
					Expect(IsValid(NewType(tests.NewMockTypeWithArgs, true, "bar"))).To(BeFalse())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should create the type", func() {
					typeDef := NewType(tests.NewMockTypeWithArgs, "foo", true)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when the arguments are variadic", func() {
				It("should create the type", func() {
					typeDef := NewType(tests.NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())
				})

				It("should return an invalid type if not enough arguments where given", func() {
					t := NewType(tests.NewVariadicMockType, true)
					Expect(t).To(BeAssignableToTypeOf(&invalidType{}))
					Expect(t.(*invalidType).Err).To(MatchError("invalid number of input parameters for variadic function: got 1 but expected at least 3"))
				})
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := NewType(tests.NewMockTypeWithArgs, args...)
			Expect(typeDef.Arguments()).To(Equal(args))
		})
	})

	Describe("Generate()", func() {
		var (
			config    = map[string]interface{}{}
			container *Container
			resolver  *ParameterResolver
		)

		BeforeEach(func() {
			container = NewContainer(NewTypeRegistry(), config)
			resolver = NewParameterResolver(container)
		})

		Context("without factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := NewType(tests.NewMockType)
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&tests.MockType{}))
			})
		})

		Context("with one or more factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := NewType(tests.NewMockTypeWithArgs, "foo", true)

				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

				generatedMock := generatedType.(*tests.MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the function signature", func() {
					It("should generate the type", func() {
						container.RegisterType("foo", tests.NewMockType)
						typeDef := NewType(tests.NewTypeForServiceInjection, "@foo")

						generatedType, err := typeDef.Generate(resolver)
						Expect(err).NotTo(HaveOccurred())
						Expect(generatedType).To(BeAssignableToTypeOf(&tests.TypeForServiceInjection{}))

						generatedMock := generatedType.(*tests.TypeForServiceInjection)
						Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&tests.MockType{}))
					})
				})

				Context("and its type does not match the function signature", func() {
					It("should return an error", func() {
						container.RegisterType("foo", tests.NewFoo)
						typeDef := NewType(tests.NewTypeForServiceInjectionWithArgs, "@foo", "arg1", "arg2", true)

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *tests.Foo) can not be passed as argument 1 to the function signature tests.NewTypeForServiceInjectionWithArgs(*tests.MockType, string, string, bool)`))
					})
				})
			})

			Context("when the arguments are variadic", func() {
				It("should generate the type", func() {
					typeDef := NewType(tests.NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

					generatedMock := generatedType.(*tests.MockType)
					Expect(generatedMock.BoolParameter).To(BeTrue())
					Expect(generatedMock.StringParameter).To(Equal("1, two, drei"))
				})
			})

			Context("when a func reference type is given", func() {
				It("should generate the type", func() {
					foo := &tests.MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := NewType(tests.NewMockTypeFromStringFunc, "YEAH", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(tests.NewMockType()))
					Expect(generatedType.(*tests.MockType).StringParameter).To(Equal("Success! YEAH"))
				})
			})

			Context("when a func reference type is given as variadic argument", func() {
				It("should generate the type", func() {
					foo := &tests.MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := NewType(tests.NewVariadicMockTypeFuncs, "@foo::ReturnString", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(tests.NewMockType()))
					Expect(generatedType.(*tests.MockType).StringParameter).To(Equal("Success! Success! "))
				})
			})
		})
	})
})
