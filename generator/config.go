package generator

const DefaultFunctionName = "RegisterTypes"

type Config struct {
	packageName  string
	functionName string
}

func NewConfig(packageName, functionName string) Config {
	if functionName == "" {
		functionName = DefaultFunctionName
	}
	return Config{packageName, functionName}
}

func (c Config) PackageName() string {
	return c.packageName
}

func (c Config) FunctionName() string {
	return c.functionName
}
