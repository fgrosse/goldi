package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
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
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected Generate to panic")
				Expect(r).To(BeAssignableToTypeOf(errors.New("")))
				err := r.(error)
				Expect(err.Error()).To(Equal("could not generate type: this type is not initialized. Did you use NewType to create it?"))
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

				generatedType := typeDef.Generate(resolver)
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

						generatedType := typeDef.Generate(resolver)
						Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.TypeForServiceInjection{}))

						generatedMock := generatedType.(*testAPI.TypeForServiceInjection)
						Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))
					})
				})

				Context("and its type does not match the function signature", func() {
					It("should panic with a helpful error message", func() {
						container.RegisterType("foo", testAPI.NewFoo)
						typeDef := goldi.NewType(testAPI.NewTypeForServiceInjectionWithArgs, "@foo", "arg1", "arg2", true)

						defer func() {
							r := recover()
							Expect(r).NotTo(BeNil(), "Expected Generate to panic")
							Expect(r).To(BeAssignableToTypeOf(errors.New("")))
							err := r.(error)
							Expect(err.Error()).To(Equal("could not generate type: the referenced type \"@foo\" (type *testAPI.Foo) can not be passed as argument 1 to the function signature testAPI.NewTypeForServiceInjectionWithArgs(*testAPI.MockType, string, string, bool)"))
						}()

						typeDef.Generate(resolver)
					})
				})
			})
		})
	})
})
