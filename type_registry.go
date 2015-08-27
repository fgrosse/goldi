package goldi

import (
	"fmt"
	"reflect"
)

// The TypeRegistry is effectively a map of typeID strings to TypeFactory
type TypeRegistry map[string]TypeFactory

// NewTypeRegistry creates a new empty TypeRegistry
func NewTypeRegistry() TypeRegistry {
	return TypeRegistry{}
}

// RegisterType is convenience method for TypeRegistry.Register
// It tries to create the correct TypeFactory and passes this to TypeRegistry.Register
// This function panics if the given generator function and arguments can not be used to create a new type factory.
func (r TypeRegistry) RegisterType(typeID string, factory interface{}, arguments ...interface{}) {
	var typeFactory TypeFactory

	factoryType := reflect.TypeOf(factory)
	kind := factoryType.Kind()
	switch {
	case kind == reflect.Struct:
		fallthrough
	case kind == reflect.Ptr && factoryType.Elem().Kind() == reflect.Struct:
		typeFactory = NewStructType(factory, arguments...)
	case kind == reflect.Func:
		typeFactory = NewType(factory, arguments...)
	default:
		panic(fmt.Errorf("could not register type %q: could not determine TypeFactory for factory type %T", typeID, factory))
	}

	r.Register(typeID, typeFactory)
}

// Register saves a type under the given symbolic typeID so it can be retrieved later.
// It is perfectly legal to call Register multiple times with the same typeID.
// In this case you overwrite existing type definitions with new once
func (r TypeRegistry) Register(typeID string, typeDef TypeFactory) {
	r[typeID] = typeDef
}

// RegisterAll will register all given type factories under the mapped type ID
// It uses TypeRegistry.Register internally
func (r TypeRegistry) RegisterAll(factories map[string]TypeFactory) {
	for typeID, typeDef := range factories {
		r.Register(typeID, typeDef)
	}
}

// InjectInstance enables you to inject type instances.
// If instance is nil an error is returned
func (r TypeRegistry) InjectInstance(typeID string, instance interface{}) {
	factory := NewInstanceType(instance)
	r.Register(typeID, factory)
}
