package goldi

import "fmt"

// The TypeRegistry is effectively a map of typeID strings to Types
type TypeRegistry map[string]TypeFactory

// NewTypeRegistry creates a new empty TypeRegistry
func NewTypeRegistry() TypeRegistry {
	return TypeRegistry{}
}

// RegisterType is convenience method for TypeRegistry.Register
// It creates a new Type from the given generatorFunction and arguments and passes this to TypeRegistry.Register
// Since the type is created in this function using NewType this can panic under the same conditions as NewType does
func (r TypeRegistry) RegisterType(typeID string, generatorFunction interface{}, arguments ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not register type: %s", r)
		}
	}()

	t := NewType(generatorFunction, arguments...)
	r.Register(typeID, t)
	return nil
}

// Register saves a type under the given symbolic typeID so it can be retrieved later.
// It is perfectly legal to call Register multiple times with the same typeID.
// In this case you overwrite existing type definitions with new once
func (r TypeRegistry) Register(typeID string, typeDef TypeFactory) {
	r[typeID] = typeDef
}

// InjectInstance enables you to inject type instances.
func (r TypeRegistry) InjectInstance(typeID string, instance interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not inject type: %s", r)
		}
	}()

	factory := NewTypeInstanceFactory(instance)
	r.Register(typeID, factory)
	return nil
}
