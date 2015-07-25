package generator
import "strings"

const DefaultFunctionName = "RegisterTypes"

type Config struct {
	packageName  string
	functionName string
}

func NewConfig(packageName, functionName string) Config {
	if functionName == "" {
		functionName = DefaultFunctionName
	}

	packageParts := strings.Split(packageName, "/")

	return Config{packageParts[len(packageParts)-1], functionName}
}

func (c Config) PackageName() string {
	return c.packageName
}

func (c Config) FunctionName() string {
	return c.functionName
}
