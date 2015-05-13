package testAPI

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
