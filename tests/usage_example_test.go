package tests

import (
	. "github.com/fgrosse/goldi/tests/testAPI"
	. "github.com/onsi/ginkgo"

	"github.com/fgrosse/goldi"
	"net/http"
)

var _ = Describe("Usage example from README.md", func() {
	It("should not crash horribly", func() {
		// create a new container when your application loads
		registry := goldi.NewTypeRegistry()
		config := map[string]interface{}{
			"some_parameter": "Hello World",
			"timeout":        42.7,
		}
		container := goldi.NewContainer(registry, config)

		// now define the types you want to build using the di container
		// you can use simple structs
		container.RegisterType("logger", &SimpleLogger{})
		container.RegisterType("api.geo.client", new(GeoClient), "http://example.com/geo:1234")

		// you can also use factory functions and parameters
		container.RegisterType("acme_corp.mailer", NewAwesomeMailer, "first argument", "%some_parameter%")

		// dynamic or static parameters and references to other services can be used as arguments
		container.RegisterType("renderer", NewRenderer, "@logger")

		// closures and functions are also possible
		container.Register("http_handler", goldi.NewFuncType(func(w http.ResponseWriter, r *http.Request) {
			// do amazing stuff
		}))

		// once you are done registering all your types you should probably validate the container
		validator := goldi.NewContainerValidator()
		validator.MustValidate(container)

		// whoever has access to the container can request these types now
		logger := container.Get("logger").(LoggerInterface)
		logger.DoStuff("...")

		// in the tests you might want to exchange the registered types with mocks or other implementations
		container.RegisterType("logger", NewNullLogger)

		// if you already have an instance you want to be used you can inject it directly
		myLogger := NewNullLogger()
		container.InjectInstance("logger", myLogger)
	})

	Describe("goldigen usage example", func() {
		It("should not crash horribly", func() {
			registry := goldi.NewTypeRegistry()
			RegisterTypes(registry)
		})
	})
})

// the following variables are just here to mock that we use code from other packages like in the README.md
var (
	mytime = new(TimePackageMock)
	example = new(ExamplePackageMock)
)

func RegisterTypes(types goldi.TypeRegistry) {
	types.RegisterType("logger", new(SimpleLogger))
	types.RegisterType("my_fancy.client", NewDefaultClient, "%client_base_url%", "@logger")
	types.RegisterType("time.clock", mytime.NewSystemClock)
	types.Register("http_handler", goldi.NewFuncType(example.HandleHTTP))
}
