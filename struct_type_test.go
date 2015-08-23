package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"fmt"
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests"
)

var _ = Describe("StructType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewStructType(tests.Foo{})
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
				typeDef := goldi.NewStructType(tests.MockType{})
				Expect(typeDef).NotTo(BeNil())
			})
		})

		Context("with first argument beeing a pointer to struct", func() {
			It("should create the type", func() {
				typeDef := goldi.NewStructType(&tests.MockType{})
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

			goldi.NewStructType(&tests.MockType{}, "foo", true, "bar")
		})
	})

	Describe("Arguments()", func() {
		It("should return all factory arguments", func() {
			args := []interface{}{"foo", true}
			typeDef := goldi.NewStructType(tests.MockType{}, args...)
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
			typeDef := &goldi.StructType{}
			defer func() {
				Expect(recover()).To(Equal("this struct type is not initialized. Did you use NewStructType to create it?"))
			}()

			typeDef.Generate(resolver)
		})

		Context("without struct arguments", func() {
			Context("when the factory is a struct (no pointer)", func() {
				It("should generate the type", func() {
					typeDef := goldi.NewStructType(tests.MockType{})
					Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&tests.MockType{}))
				})
			})

			It("should generate the type", func() {
				typeDef := goldi.NewStructType(&tests.MockType{})
				Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(&tests.MockType{}))
			})

			It("should generate a new type each time", func() {
				typeDef := goldi.NewStructType(&tests.MockType{})
				t1, err1 := typeDef.Generate(resolver)
				t2, err2 := typeDef.Generate(resolver)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				Expect(t1).NotTo(BeNil())
				Expect(t2).NotTo(BeNil())
				Expect(t1 == t2).To(BeFalse(), fmt.Sprintf("t1 (%p) should not point to the same instance as t2 (%p)", t1, t2))

				// Just to make the whole issue more explicit:
				t1Mock := t1.(*tests.MockType)
				t2Mock := t2.(*tests.MockType)
				t1Mock.StringParameter = "CHANGED"
				Expect(t2Mock.StringParameter).NotTo(Equal(t1Mock.StringParameter),
					"Changing two indipendently generated types should not affect both at the same time",
				)
			})
		})

		Context("with one or more arguments", func() {
			It("should generate the type", func() {
				typeDef := goldi.NewStructType(&tests.MockType{}, "foo", true)

				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

				generatedMock := generatedType.(*tests.MockType)
				Expect(generatedMock.StringParameter).To(Equal("foo"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			It("should use the given parameters", func() {
				typeDef := goldi.NewStructType(&tests.MockType{}, "%param1%", "%param2%")
				config["param1"] = "TEST"
				config["param2"] = true
				generatedType, err := typeDef.Generate(resolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(generatedType).To(BeAssignableToTypeOf(&tests.MockType{}))

				generatedMock := generatedType.(*tests.MockType)
				Expect(generatedMock.StringParameter).To(Equal("TEST"))
				Expect(generatedMock.BoolParameter).To(Equal(true))
			})

			Context("when a type reference is given", func() {
				Context("and its type matches the struct field type", func() {
					It("should generate the type", func() {
						container.RegisterType("foo", tests.NewMockType)
						typeDef := goldi.NewStructType(tests.TypeForServiceInjection{}, "@foo")

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
						typeDef := goldi.NewStructType(tests.TypeForServiceInjection{}, "@foo")

						_, err := typeDef.Generate(resolver)
						Expect(err).To(MatchError(`the referenced type "@foo" (type *tests.Foo) can not be used as field 1 for struct type tests.TypeForServiceInjection`))
					})
				})
			})
		})
	})
})