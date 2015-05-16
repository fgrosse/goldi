package tests

import (
	. "github.com/fgrosse/goldi/tests/testAPI"
	. "github.com/onsi/ginkgo"

	"github.com/fgrosse/goldi"
)

var _ = Describe("Usage example from README.md", func() {
	It("should not crash horribly", func() {
		registry := goldi.NewTypeRegistry()
		config := map[string]interface{}{
			"some_parameter": "Hello World",
			"timeout":        42.7,
		}
		container := goldi.NewContainer(registry, config)

		// now define the types you want to build using the di container
		// you can use simple structs in case you do not have a factory function
		container.RegisterType("logger", &SimpleLogger{})

		// you can also use factory functions and parameters
		container.RegisterType("acme_corp.mailer", NewAwesomeMailer, "first argument", "%some_parameter%")

		// dynamic or static parameters and references to other services can be used as arguments
		container.RegisterType("renderer", NewRenderer, "@logger")

		// once you are done registering all your types you should probably validate the container
		validator := goldi.NewContainerValidator()
		err := validator.Validate(container)
		if err != nil {
			panic(err)
		}

		// whoever has access to the container can request these types now
		logger := container.Get("logger").(LoggerInterface)
		logger.DoStuff("...")

		// in the tests you might want to exchange the registered types with mocks or other implementations
		container.RegisterType("logger", NewNullLogger)

		// if you already have an instance you want to be used you can inject it directly
		myLogger := NewNullLogger()
		container.InjectInstance("logger", myLogger)
	})
})
