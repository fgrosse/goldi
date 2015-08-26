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

// NewConfiguredType creates a new TypeFactory that decorates a given TypeFactory.
// The returned configurator will use the decorated type factory first to create a type and then use
// the resolve the configurator by the given type ID and call the configured method with the instance.
//
// Internally the goldi.TypeConfigurator is used.
//
// The method removes any leading or trailing whitespace from configurator type ID and method.
// NewConfiguredType will return an invalid type when embeddedType is nil or the trimmed configurator typeID or method is empty.
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
