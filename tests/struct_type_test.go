package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("StructType", func() {
	var typeDef *goldi.StructType

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewStructType(testAPI.Foo{})
		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewStructType()", func() {
		Context("with invalid arguments", func() {
			It("should panic if the generator is no struct or pointer to a struct", func() {
				Expect(func() { goldi.NewStructType(42) }).To(Panic())
			})

			It("should panic if the generator is a pointer to something other than a struct", func() {
				something := "Hello Pointer World!"
				Expect(func() { goldi.NewStructType(&something) }).To(Panic())
			})
		})

		Context("with first argument beeing a struct", func() {
			It("should create the type", func() {
				typeDef = goldi.NewStructType(testAPI.MockType{})
				Expect(typeDef).NotTo(BeNil())
			})
		})

		Context("with first argument beeing a pointer to struct", func() {
			It("should create the type", func() {
				typeDef = goldi.NewStructType(&testAPI.MockType{})
				Expect(typeDef).NotTo(BeNil())
			})
		})

		It("should panic if more factory arguments where provided than the struct has fields", func() {
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected Generate to panic")
				Expect(r).To(BeAssignableToTypeOf(errors.New("")))
				err := r.(error)
				Expect(err.Error()).To(Equal("could not register struct type: the struct MockType has only 2 fields but 3 arguments where provided"))
			}()

			typeDef = goldi.NewStructType(&testAPI.MockType{}, "foo", true, "bar")
		})
	})

	Describe("Generate()", func() {
		var (
			config       = map[string]interface{}{}
			typeRegistry goldi.TypeRegistry
		)

		BeforeEach(func() {
			typeRegistry = goldi.NewTypeRegistry()
		})

		It("should panic if Generate is called on an uninitialized type", func() {
			typeDef = &goldi.StructType{}
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected Generate to panic")
				Expect(r).To(BeAssignableToTypeOf(errors.New("")))
				err := r.(error)
				Expect(err.Error()).To(Equal("could not generate type: this struct type is not initialized. Did you use NewStructType to create it?"))
			}()

			typeDef.Generate(config, typeRegistry)
		})

		Context("without struct arguments", func() {
			Context("when the factory is a struct (no pointer)", func() {
				It("should generate the type", func() {
					typeDef = goldi.NewStructType(testAPI.MockType{})
					Expect(typeDef.Generate(config, typeRegistry)).To(BeAssignableToTypeOf(&testAPI.MockType{}))
				})
			})

			It("should generate the type", func() {
				typeDef = goldi.NewStructType(&testAPI.MockType{})
				Expect(typeDef.Generate(config, typeRegistry)).To(BeAssignableToTypeOf(&testAPI.MockType{}))
			})

			It("should generate a new type each time", func() {
				typeDef = goldi.NewStructType(&testAPI.MockType{})
				t1 := typeDef.Generate(config, typeRegistry)
				t2 := typeDef.Generate(config, typeRegistry)

				Expect(t1).NotTo(BeNil())
				Expect(t2).NotTo(BeNil())
				Expect(t1 == t2).To(BeFalse(), fmt.Sprintf("t1 (%p) should not point to the same instance as t2 (%p)", t1, t2))

				// Just to make the whole issue more explicit:
				t1Mock := t1.(*testAPI.MockType)
				t2Mock := t2.(*testAPI.MockType)
				t1Mock.StringParameter = "CHANGED"
				Expect(t2Mock.StringParameter).NotTo(Equal(t1Mock.StringParameter),
					"Changing two indipendently generated types should not affect both at the same time",
				)
			})
		})

		Context("with one or more arguments", func() {
			It("should generate the type", func() {
				typeDef = goldi.NewStructType(&testAPI.MockType{}, "foo", true)

				generatedType := typeDef.Generate(config, typeRegistry)
				Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

				generatedMock := generatedType.(*testAPI.MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			It("should use the given parameters", func() {
				typeDef = goldi.NewStructType(&testAPI.MockType{}, "%param1%", "%param2%")
				config["param1"] = "TEST"
				config["param2"] = true
				generatedType := typeDef.Generate(config, typeRegistry)
				Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))

				generatedMock := generatedType.(*testAPI.MockType)
				Expect(generatedMock.StringParameter).To(Equal("TEST"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the struct field type", func() {
					It("should generate the type", func() {
						typeRegistry.RegisterType("foo", testAPI.NewMockType)
						typeDef = goldi.NewStructType(testAPI.TypeForServiceInjection{}, "@foo")

						generatedType := typeDef.Generate(config, typeRegistry)
						Expect(generatedType).To(BeAssignableToTypeOf(&testAPI.TypeForServiceInjection{}))

						generatedMock := generatedType.(*testAPI.TypeForServiceInjection)
						Expect(generatedMock.InjectedType).To(BeAssignableToTypeOf(&testAPI.MockType{}))
					})
				})

				Context("and its type does not match the function signature", func() {
					It("should panic with a helpful error message", func() {
						typeRegistry.RegisterType("foo", testAPI.NewFoo)
						typeDef = goldi.NewStructType(testAPI.TypeForServiceInjection{}, "@foo")

						defer func() {
							r := recover()
							Expect(r).NotTo(BeNil(), "Expected Generate to panic")
							Expect(r).To(BeAssignableToTypeOf(errors.New("")))
							err := r.(error)
							Expect(err.Error()).To(Equal("could not generate type: the referenced type \"@foo\" (type *testAPI.Foo) can not be used as field 1 for struct type testAPI.TypeForServiceInjection"))
						}()

						typeDef.Generate(config, typeRegistry)
					})
				})
			})
		})
	})
})
