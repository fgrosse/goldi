package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"reflect"
)

var _ = Describe("ParameterResolver", func() {

	var (
		config    map[string]interface{}
		container *goldi.Container
		resolver  *goldi.ParameterResolver
	)

	BeforeEach(func() {
		config = map[string]interface{}{}
		container = goldi.NewContainer(goldi.NewTypeRegistry(), config)
		resolver = goldi.NewParameterResolver(container)
	})

	It("should return static parameters", func() {
		parameter := reflect.ValueOf(42)
		expectedType := parameter.Type()
		Expect(resolver.Resolve(parameter, expectedType)).To(Equal(parameter))
	})

	Context("with dynamic parameters", func() {
		Context("with invalid parameter name", func() {
			It("should not try to resolve `%`", func() {
				parameter := reflect.ValueOf("%")
				result, err := resolver.Resolve(parameter, parameter.Type())
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal("%"))
			})

			It("should not try to resolve `%%`", func() {
				parameter := reflect.ValueOf("%%")
				result, err := resolver.Resolve(parameter, parameter.Type())
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal("%%"))
			})
		})

		Context("when parameter has not been defined", func() {
			It("should return the parameter as is", func() {
				parameter := reflect.ValueOf("%foo%")
				expectedType := parameter.Type()

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal("%foo%"))
			})
		})

		Context("when the parameter has been defined", func() {
			It("should resolve string parameters using the configuration", func() {
				config["foo"] = "success"
				parameter := reflect.ValueOf("%foo%")
				expectedType := reflect.TypeOf(config["foo"])

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal(config["foo"]))
			})

			It("should resolve float parameters using the configuration", func() {
				config["bar"] = 3.1415
				parameter := reflect.ValueOf("%bar%")
				expectedType := reflect.TypeOf(config["bar"])

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal(config["bar"]))
			})
		})
	})

	Context("with type references", func() {
		Context("when the type has been registered", func() {
			BeforeEach(func() {
				container.RegisterType("foo", NewFoo)
			})

			Context("when the type is assignable to the expected type", func() {
				It("should generate the type and return it", func() {
					parameter := reflect.ValueOf("@foo")
					expectedType := reflect.TypeOf(NewFoo())

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Interface()).To(BeAssignableToTypeOf(NewFoo()))
				})
			})

			Context("when the type is not assignable to the expected type", func() {
				It("should return an error", func() {
					parameter := reflect.ValueOf("@foo")
					expectedType := reflect.TypeOf(NewBar())

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(result.IsValid()).To(BeFalse())
					Expect(err).To(MatchError(`the referenced type "@foo" (type *goldi_test.Foo) is not assignable to the expected type *goldi_test.Bar`))
					Expect(err.(goldi.TypeReferenceError).TypeID).To(Equal("foo"))
				})
			})

			Context("when a func reference type is requested", func() {
				It("should generate the type and return the function", func() {
					foo := &Foo{Value: "Success!"}
					container.InjectInstance("foo", foo)

					parameter := reflect.ValueOf("@foo::ReturnString")
					expectedType := reflect.TypeOf(func(string) string { return "" })

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Interface()).To(BeAssignableToTypeOf(foo.ReturnString))
					Expect(result.Interface().(func(string) string)("YEAH")).To(Equal("Success! YEAH"))
				})

				It("should return an error if the method does not exist", func() {
					container.InjectInstance("foo", new(Foo))

					parameter := reflect.ValueOf("@foo::ThisMethodDoesNotExist")
					expectedType := reflect.TypeOf(func() {})

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(result.IsValid()).To(BeFalse())
					Expect(err).To(MatchError(`the referenced method "@foo::ThisMethodDoesNotExist" does not exist or is not exported`))
				})
			})
		})

		Context("when the type has not been registered", func() {
			It("should return an error", func() {
				parameter := reflect.ValueOf("@foo")
				expectedType := reflect.TypeOf(&Foo{})

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).To(HaveOccurred())
				Expect(result.IsValid()).To(BeFalse())
				Expect(err).To(MatchError(`the referenced type "@foo" has not been defined`))
				Expect(err.(goldi.UnknownTypeReferenceError).TypeID).To(Equal("foo"))
			})
		})

		Context("when the type has is invalid", func() {
			It("should return an error", func() {
				parameter := reflect.ValueOf("@foo")
				expectedType := reflect.TypeOf(&Foo{})

				container.Register("foo", goldi.NewType(nil)) // foo will be invalid
				result, err := resolver.Resolve(parameter, expectedType)
				Expect(result).To(Equal(reflect.Zero(expectedType)))
				Expect(err).To(MatchError(`goldi: error while generating type "foo": the given factoryFunction is nil`))
			})
		})
	})
})
