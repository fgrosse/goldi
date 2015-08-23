package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("Type", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewType(tests.NewFoo)
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewType()", func() {
		Context("with invalid factory function", func() {
			It("should panic if the generator is no function", func() {
				Expect(func() { goldi.NewType(42) }).To(Panic())
			})

			It("should panic if the generator has no output parameters", func() {
				Expect(func() { goldi.NewType(func() {}) }).To(Panic())
			})

			It("should panic if the generator has more than one output parameter", func() {
				Expect(func() { goldi.NewType(func() (*tests.MockType, *tests.MockType) { return nil, nil }) }).To(Panic())
			})

			It("should panic if the return parameter is no pointer", func() {
				Expect(func() { goldi.NewType(func() tests.MockType { return tests.MockType{} }) }).To(Panic())
			})

			It("should not panic if the return parameter is an interface", func() {
				Expect(func() { goldi.NewType(func() interface{} { return tests.MockType{} }) }).NotTo(Panic())
			})
		})

		Context("without factory function arguments", func() {
			Context("when no factory argument is given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(tests.NewMockType)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when any argument is given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(tests.NewMockType, "foo") }).To(Panic())
				})
			})
		})

		Context("with one or more factory function arguments", func() {
			Context("when an invalid number of arguments is given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(tests.NewMockTypeWithArgs) }).To(Panic())
					Expect(func() { goldi.NewType(tests.NewMockTypeWithArgs, "foo") }).To(Panic())
					Expect(func() { goldi.NewType(tests.NewMockTypeWithArgs, "foo", false, 42) }).To(Panic())
				})
			})

			Context("when the wrong argument types are given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(tests.NewMockTypeWithArgs, "foo", "bar") }).To(Panic())
					Expect(func() { goldi.NewType(tests.NewMockTypeWithArgs, true, "bar") }).To(Panic())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(tests.NewMockTypeWithArgs, "foo", true)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when the arguments are variadic", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(tests.NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())
				})

				It("should return an error if not enough arguments where given", func() {
					defer func() {
						Expect(recover()).To(MatchError("could not register type: invalid number of input parameters for variadic function: got 1 but expected at least 3"))
					}()

					goldi.NewType(tests.NewVariadicMockType, true)
				})
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := goldi.NewType(tests.NewMockTypeWithArgs, args...)
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

		It("should panic if Generate is called on an uninitialized type", func() {
			typeDef := &goldi.Type{}
			defer func() {
				Expect(recover()).To(MatchError("this type is not initialized. Did you use NewType to create it?"))
			}()

			typeDef.Generate(resolver)
		})

		Context("without factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewType(tests.NewMockType)
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&tests.MockType{}))
			})
		})

		Context("with one or more factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewType(tests.NewMockTypeWithArgs, "foo", true)

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
						typeDef := goldi.NewType(tests.NewTypeForServiceInjection, "@foo")

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
						typeDef := goldi.NewType(tests.NewTypeForServiceInjectionWithArgs, "@foo", "arg1", "arg2", true)

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *tests.Foo) can not be passed as argument 1 to the function signature tests.NewTypeForServiceInjectionWithArgs(*tests.MockType, string, string, bool)`))
					})
				})
			})

			Context("when the arguments are variadic", func() {
				It("should generate the type", func() {
					typeDef := goldi.NewType(tests.NewVariadicMockType, true, "ignored", "1", "two", "drei")
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
					typeDef := goldi.NewType(tests.NewMockTypeFromStringFunc, "YEAH", "@foo::ReturnString")

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
					typeDef := goldi.NewType(tests.NewVariadicMockTypeFuncs, "@foo::ReturnString", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(tests.NewMockType()))
					Expect(generatedType.(*tests.MockType).StringParameter).To(Equal("Success! Success! "))
				})
			})
		})
	})
})
