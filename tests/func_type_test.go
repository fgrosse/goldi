package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"errors"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("FuncType", func() {
	var typeDef *goldi.FuncType

	Describe("Short usage example", func() {
		typeRegistry := goldi.NewTypeRegistry()
		resolver := goldi.NewParameterResolver(map[string]interface{}{}, typeRegistry)

		It("should work with a defined function", func() {
			// define the type
			typeDef = goldi.NewFuncType(SomeFunctionForFuncTypeTest)

			// generate it
			f := typeDef.Generate(resolver).(func(name string, age int) (bool, error))

			// call it
			ok, err := f("foo", 42)

			Expect(ok).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should work with closures", func() {
			typeDef = goldi.NewFuncType(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "test" {
					w.WriteHeader(http.StatusAccepted)
				}
			})

			request, _ := http.NewRequest("GET", "test", nil)
			response := httptest.NewRecorder()
			handler := typeDef.Generate(resolver).(func(w http.ResponseWriter, r *http.Request))
			handler(response, request)

			Expect(response.Code).To(Equal(http.StatusAccepted))
		})
	})

	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewFuncType(SomeFunctionForFuncTypeTest)

		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewFuncType()", func() {
		Context("with invalid argument", func() {
			It("should panic if the argument is no function", func() {
				Expect(func() { goldi.NewFuncType(42) }).To(Panic())
			})
		})

		Context("with argument beeing a function", func() {
			It("should create the type", func() {
				typeDef = goldi.NewFuncType(SomeFunctionForFuncTypeTest)
				Expect(typeDef).NotTo(BeNil())
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef = goldi.NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(BeEmpty())
		})
	})

	Describe("Generate()", func() {
		var (
			config       = map[string]interface{}{}
			typeRegistry goldi.TypeRegistry
			resolver     *goldi.ParameterResolver
		)

		BeforeEach(func() {
			typeRegistry = goldi.NewTypeRegistry()
			resolver = goldi.NewParameterResolver(config, typeRegistry)
		})

		It("should panic if Generate is called on an uninitialized type", func() {
			typeDef = &goldi.FuncType{}
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil(), "Expected Generate to panic")
				Expect(r).To(BeAssignableToTypeOf(errors.New("")))
				err := r.(error)
				Expect(err.Error()).To(Equal("could not generate type: this func type is not initialized. Did you use NewFuncType to create it?"))
			}()

			typeDef.Generate(resolver)
		})

		It("should just return the function", func() {
			typeDef = goldi.NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(SomeFunctionForFuncTypeTest))
		})
	})
})

func SomeFunctionForFuncTypeTest(name string, age int) (bool, error) {
	return true, nil
}
