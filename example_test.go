package goldi_test

import (
	"net/http"
	"time"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/validation"
)

func Example() {
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
	validator := validation.NewContainerValidator()
	validator.MustValidate(container) // will panic, use validator.Validate to get the error

	// whoever has access to the container can request these types now
	logger := container.MustGet("logger").(LoggerInterface)
	logger.DoStuff("...")

	// in the tests you might want to exchange the registered types with mocks or other implementations
	container.RegisterType("logger", NewNullLogger)

	// if you already have an instance you want to be used you can inject it directly
	myLogger := NewNullLogger()
	container.InjectInstance("logger", myLogger)
}

// Example_ prevents godoc from printing the whole content of this file as example
func Example_() {}

type LoggerInterface interface {
	DoStuff(message string) string
}

type LoggerProvider struct{}

func (p *LoggerProvider) GetLogger(name string) LoggerInterface {
	return &SimpleLogger{name}
}

type SimpleLogger struct{
	Name string
}

func (l *SimpleLogger) DoStuff(input string) string { return input }

type NullLogger struct{}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

func (l *NullLogger) DoStuff(input string) string { return input }

type AwesomeMailer struct {
	arg1, arg2 string
}

func NewAwesomeMailer(arg1, arg2 string) *AwesomeMailer {
	return &AwesomeMailer{arg1, arg2}
}

type Renderer struct {
	logger *LoggerInterface
}

func NewRenderer(logger *LoggerInterface) *Renderer {
	return &Renderer{logger}
}

type GeoClient struct {
	BaseURL string
}

type TimePackageMock struct{}

func (M *TimePackageMock) NewSystemClock() *time.Time {
	now := time.Now()
	return &now
}

type ExamplePackageMock struct{}

func (M *ExamplePackageMock) HandleHTTP() {}
