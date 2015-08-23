package goldi

import (
	"fmt"
	"reflect"
)

// The TypeConfigurator is used to configure a type after its instantiation.
// You can specify a function in another type that is known to the container.
// The type instance is passed to the configurator type, allowing the
// configurator to do whatever it needs to configure the type after its creation.
//
// A TypeConfigurator can be used, for example, when you have a type that requires
// complex setup based on configuration settings coming from different sources.
// Using an external configurator, you can decouple the setup logic from the business logic
// of the corresponding type to keep it DRY and easy to maintain. Also this way its easy to
// exchange setup logic at run time for example on different environments.
//
// Another interesting use case is when you have multiple objects that share a common
// configuration or that should be configured in a similar way at runtime.
type TypeConfigurator struct {
	ConfiguratorTypeID string
	MethodName         string
}

func NewTypeConfigurator(configuratorTypeID, methodName string) *TypeConfigurator {
	return &TypeConfigurator{
		ConfiguratorTypeID: configuratorTypeID,
		MethodName:         methodName,
	}
}

// Configure will get the configurator type and ass `thing` its configuration function.
// The method returns an error if thing is nil, the configurator type is not defined or
// the configurators function does not exist.
func (c *TypeConfigurator) Configure(thing interface{}, container *Container) error {
	if thing == nil {
		return fmt.Errorf("can not configure nil")
	}

	configurator, typeDefined := container.get(c.ConfiguratorTypeID)
	if typeDefined == false {
		return NewUnknownTypeReferenceError(c.ConfiguratorTypeID, `the configurator type "@%s" has not been defined`, c.ConfiguratorTypeID)
	}

	configuratorType := reflect.TypeOf(configurator)
	configuratorKind := configuratorType.Kind()
	if configuratorKind != reflect.Struct && !(configuratorKind == reflect.Ptr && configuratorType.Elem().Kind() == reflect.Struct) {
		return NewTypeReferenceError(c.ConfiguratorTypeID, configuratorType, "the configurator instance is no struct or pointer to struct but a %T", configurator)
	}

	configuratorValue := reflect.ValueOf(configurator)
	configuratorMethod := configuratorValue.MethodByName(c.MethodName)
	if configuratorMethod.IsValid() == false {
		return NewTypeReferenceError(c.ConfiguratorTypeID, configuratorType, "the configurator does not have a method %q", c.MethodName)
	}

	result := configuratorMethod.Call([]reflect.Value{reflect.ValueOf(thing)})
	if len(result) > 0 {
		lastResult := result[len(result)-1]

		errType := reflect.TypeOf((*error)(nil)).Elem()
		if lastResult.Type().AssignableTo(errType) {
			return lastResult.Interface().(error)
		}
	}

	return nil
}
