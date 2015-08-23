package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("Type", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewType(testAPI.NewFoo)
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
				Expect(func() { goldi.NewType(func() (*testAPI.MockType, *testAPI.MockType) { return nil, nil }) }).To(Panic())
			})

			It("should panic if the return parameter is no pointer", func() {
				Expect(func() { goldi.NewType(func() testAPI.MockType { return testAPI.MockType{} }) }).To(Panic())
			})

			It("should not panic if the return parameter is an interface", func() {
				Expect(func() { goldi.NewType(func() interface{} { return testAPI.MockType{} }) }).NotTo(Panic())
			})
		})

		Context("without factory function arguments", func() {
			Context("when no factory argument is given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(testAPI.NewMockType)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when any argument is given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(testAPI.NewMockType, "foo") }).To(Panic())
				})
			})
		})

		Context("with one or more factory function arguments", func() {
			Context("when an invalid number of arguments is given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(testAPI.NewMockTypeWithArgs) }).To(Panic())
					Expect(func() { goldi.NewType(testAPI.NewMockTypeWithArgs, "foo") }).To(Panic())
					Expect(func() { goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", false, 42) }).To(Panic())
				})
			})

			Context("when the wrong argument types are given", func() {
				It("should panic", func() {
					Expect(func() { goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", "bar") }).To(Panic())
					Expect(func() { goldi.NewType(testAPI.NewMockTypeWithArgs, true, "bar") }).To(Panic())
				})
			})

			Context("when the correct argument number and types are given", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", true)
					Expect(typeDef).NotTo(BeNil())
				})
			})

			Context("when the arguments are variadic", func() {
				It("should create the type", func() {
					typeDef := goldi.NewType(testAPI.NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())
				})

				It("should return an error if not enough arguments where given", func() {
					defer func() {
						Expect(recover()).To(MatchError("could not register type: invalid number of input parameters for variadic function: got 1 but expected at least 3"))
					}()

					goldi.NewType(testAPI.NewVariadicMockType, true)
				})
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, args...)
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
				typeDef := goldi.NewType(testAPI.NewMockType)
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&testAPI.MockType{}))
			})
		})

		Context("with one or more factory function arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewType(testAPI.NewMockTypeWithArgs, "foo", true)

				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

				generatedMock := generatedType.(*testAPI.MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the function signature", func() {
					It("should generate the type", func() {
						container.RegisterType("foo", testAPI.NewMockType)
						typeDef := goldi.NewType(testAPI.NewTypeForServiceInjection, "@foo")

						generatedType, err := typeDef.Generate(resolver)
						Expect(err).NotTo(HaveOccurred())
						Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.TypeForServiceInjection{}))

						generatedMock := generatedType.(*testAPI.TypeForServiceInjection)
						Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))
					})
				})

				Context("and its type does not match the function signature", func() {
					It("should return an error", func() {
						container.RegisterType("foo", testAPI.NewFoo)
						typeDef := goldi.NewType(testAPI.NewTypeForServiceInjectionWithArgs, "@foo", "arg1", "arg2", true)

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *testAPI.Foo) can not be passed as argument 1 to the function signature testAPI.NewTypeForServiceInjectionWithArgs(*testAPI.MockType, string, string, bool)`))
					})
				})
			})

			Context("when the arguments are variadic", func() {
				It("should generate the type", func() {
					typeDef := goldi.NewType(testAPI.NewVariadicMockType, true, "ignored", "1", "two", "drei")
					Expect(typeDef).NotTo(BeNil())

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

					generatedMock := generatedType.(*testAPI.MockType)
					Expect(generatedMock.BoolParameter).To(BeTrue())
					Expect(generatedMock.StringParameter).To(Equal("1, two, drei"))
				})
			})

			Context("when a func reference type is given", func() {
				It("should generate the type", func() {
					foo := &testAPI.MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := goldi.NewType(testAPI.NewMockTypeFromStringFunc, "YEAH", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(testAPI.NewMockType()))
					Expect(generatedType.(*testAPI.MockType).StringParameter).To(Equal("Success! YEAH"))
				})
			})

			Context("when a func reference type is given as variadic argument", func() {
				It("should generate the type", func() {
					foo := &testAPI.MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)
					typeDef := goldi.NewType(testAPI.NewVariadicMockTypeFuncs, "@foo::ReturnString", "@foo::ReturnString")

					generatedType, err := typeDef.Generate(resolver)
					Expect(err).NotTo(HaveOccurred())
					Expect(generatedType).To(BeAssignableToTypeOf(testAPI.NewMockType()))
					Expect(generatedType.(*testAPI.MockType).StringParameter).To(Equal("Success! Success! "))
				})
			})
		})
	})
})
