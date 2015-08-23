package goldi

import (
	"fmt"
	"strings"
	"unicode"
)

type configuredType struct {
	*TypeConfigurator
	embeddedType TypeFactory
}

func NewConfiguredType(embeddedType TypeFactory, configuratorTypeID, configuratorMethod string) TypeFactory {
	if embeddedType == nil {
		return newInvalidType(fmt.Errorf("refusing to create a new ConfiguredType with nil as embedded type"))
	}

	configuratorTypeID = strings.TrimSpace(configuratorTypeID)
	configuratorMethod = strings.TrimSpace(configuratorMethod)

	if configuratorTypeID == "" || configuratorMethod == "" {
		return newInvalidType(fmt.Errorf("can not create a new ConfiguredType with empty configurator type or method (%q, %q)", configuratorTypeID, configuratorMethod))
	}

	if unicode.IsLower(rune(configuratorMethod[0])) {
		return newInvalidType(fmt.Errorf("can not create a new ConfiguredType with unexproted configurator method %q", configuratorMethod))
	}

	return &configuredType{
		TypeConfigurator: NewTypeConfigurator(configuratorTypeID, configuratorMethod),
		embeddedType:     embeddedType,
	}
}

func (t *configuredType) Arguments() []interface{} {
	return append(t.embeddedType.Arguments(), "@"+t.ConfiguratorTypeID)
}

func (t *configuredType) Generate(parameterResolver *ParameterResolver) (interface{}, error) {
	embedded, err := t.embeddedType.Generate(parameterResolver)
	if err != nil {
		return nil, fmt.Errorf("can not generate configured type: %s", err)
	}

	if err = t.Configure(embedded, parameterResolver.Container); err != nil {
		return nil, err
	}

	return embedded, nil
}
