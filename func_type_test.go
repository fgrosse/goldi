package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("FuncType", func() {
	Describe("Short usage example", func() {
		container := NewContainer(NewTypeRegistry(), map[string]interface{}{})
		resolver := NewParameterResolver(container)

		It("should work with a defined function", func() {
			// define the type
			typeDef := NewFuncType(SomeFunctionForFuncTypeTest)

			// generate it
			result, err := typeDef.Generate(resolver)
			Expect(result).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			// call it
			f := result.(func(name string, age int) (bool, error))
			ok, err := f("foo", 42)

			Expect(ok).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should work with closures", func() {
			typeDef := NewFuncType(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "test" {
					w.WriteHeader(http.StatusAccepted)
				}
			})

			request, _ := http.NewRequest("GET", "test", nil)
			response := httptest.NewRecorder()
			result, err := typeDef.Generate(resolver)
			Expect(result).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			handler := result.(func(w http.ResponseWriter, r *http.Request))
			handler(response, request)
			Expect(response.Code).To(Equal(http.StatusAccepted))
		})
	})

	It("should implement the TypeFactory interface", func() {
		var factory TypeFactory
		factory = NewFuncType(SomeFunctionForFuncTypeTest)

		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("NewFuncType()", func() {
		Context("with invalid argument", func() {
			It("should return an invalid type if the argument is no function", func() {
				Expect(IsValid(NewFuncType(42))).To(BeFalse())
			})
		})

		Context("with argument beeing a function", func() {
			It("should create the type", func() {
				typeDef := NewFuncType(SomeFunctionForFuncTypeTest)
				Expect(typeDef).NotTo(BeNil())
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(BeEmpty())
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

		It("should just return the function", func() {
			typeDef := NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(SomeFunctionForFuncTypeTest))
		})
	})
})

func SomeFunctionForFuncTypeTest(name string, age int) (bool, error) {
	return true, nil
}
