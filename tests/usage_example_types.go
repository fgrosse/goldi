package tests

import "time"

type LoggerInterface interface {
	DoStuff(message string)
}

type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) DoStuff(_ string) {}

type NullLogger struct{}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

func (l *NullLogger) DoStuff(_ string) {}

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

func NewDefaultClient(baseURL string, logger LoggerInterface) *GeoClient {
	return &GeoClient{baseURL}
}

type TimePackageMock struct{}

func (M *TimePackageMock) NewSystemClock() *time.Time {
	now := time.Now()
	return &now
}

type ExamplePackageMock struct{}

func (M *ExamplePackageMock) HandleHTTP() {
	// foo
}
