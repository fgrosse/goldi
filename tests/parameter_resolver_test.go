package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
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
				config["foo"] = "success"
				parameter := reflect.ValueOf("%foo%")
				expectedType := parameter.Type()

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Interface()).To(Equal(config["foo"]))
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
				container.RegisterType("foo", testAPI.NewFoo)
			})

			Context("when the type is assignable to the expected type", func() {
				It("should generate the type and return it", func() {
					parameter := reflect.ValueOf("@foo")
					expectedType := reflect.TypeOf(testAPI.NewFoo())

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Interface()).To(BeAssignableToTypeOf(testAPI.NewFoo()))
				})
			})

			Context("when the type is not assignable to the expected type", func() {
				It("should return an error", func() {
					parameter := reflect.ValueOf("@foo")
					expectedType := reflect.TypeOf(testAPI.NewBar())

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(goldi.NewTypeReferenceError("foo", testAPI.NewFoo(), `the referenced type "@foo" (type *testAPI.Foo) is not assignable to the expected type *testAPI.Bar`)))
					Expect(result.IsValid()).To(BeFalse())
				})
			})

			Context("when a func reference type is requested", func() {
				It("should generate the type and return the function", func() {
					foo := &testAPI.MockType{StringParameter: "Success!"}
					container.InjectInstance("foo", foo)

					parameter := reflect.ValueOf("@foo::ReturnString")
					expectedType := reflect.TypeOf(func(string) string { return "" })

					result, err := resolver.Resolve(parameter, expectedType)
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Interface()).To(BeAssignableToTypeOf(foo.ReturnString))
					Expect(result.Interface().(func(string) string )("YEAH")).To(Equal("Success! YEAH"))
				})
			})
		})

		Context("when the type has not been registered", func() {
			It("should return an error", func() {
				parameter := reflect.ValueOf("@foo")
				expectedType := reflect.TypeOf(testAPI.NewMockType())

				result, err := resolver.Resolve(parameter, expectedType)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(goldi.NewUnknownTypeReferenceError("foo", `the referenced type "@foo" has not been defined`)))
				Expect(result.IsValid()).To(BeFalse())
			})
		})
	})
})
