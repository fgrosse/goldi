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
		container.RegisterType("logger", NewSimpleLogger)
		container.RegisterType("acme_corp.mailer", NewAwesomeMailer, "first argument", "%some_parameter%")

		// whoever has access to the container can request these types now
		logger := container.Get("logger").(LoggerInterface)
		logger.DoStuff("...")

		// in the tests you might want to exchange the registered types with mocks or other implementations
		container.RegisterType("logger", NewNullLogger)
	})
})
