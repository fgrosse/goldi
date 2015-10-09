package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"net/http"
)

func ExampleNewFuncType() {
	container := goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})

	// define the type
	container.Register("my_func", goldi.NewFuncType(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "test" {
			w.WriteHeader(http.StatusAccepted)
		}
	}))

	// generate it
	result, err := container.Get("my_func")
	if err != nil {
		return
	}

	// call it
	f := result.(func(name string, age int) (bool, error))
	ok, err := f("foo", 42)
	if ok != true || err != nil {
		panic("!!!")
	}
}

// ExampleNewFuncType_ prevents godoc from printing the whole content of this file as example
func ExampleNewFuncType_() {}

var _ = Describe("funcType", func() {
	It("should implement the TypeFactory interface", func() {
		var factory goldi.TypeFactory
		factory = goldi.NewFuncType(SomeFunctionForFuncTypeTest)

		// if this compiles the test passes (next expectation only to make compiler happy)
		Expect(factory).NotTo(BeNil())
	})

	Describe("goldi.NewFuncType()", func() {
		Context("with invalid argument", func() {
			It("should return an invalid type if the argument is no function", func() {
				Expect(goldi.IsValid(goldi.NewFuncType(42))).To(BeFalse())
			})
		})

		Context("with argument beeing a function", func() {
			It("should create the type", func() {
				typeDef := goldi.NewFuncType(SomeFunctionForFuncTypeTest)
				Expect(typeDef).NotTo(BeNil())
			})
		})
	})

	Describe("Arguments()", func() {
		It("should return an empty list", func() {
			typeDef := goldi.NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Arguments()).NotTo(BeNil())
			Expect(typeDef.Arguments()).To(BeEmpty())
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

		It("should just return the function", func() {
			typeDef := goldi.NewFuncType(SomeFunctionForFuncTypeTest)
			Expect(typeDef.Generate(resolver)).To(BeAssignableToTypeOf(SomeFunctionForFuncTypeTest))
		})
	})
})

func SomeFunctionForFuncTypeTest(name string, age int) (bool, error) {
	return true, nil
}
