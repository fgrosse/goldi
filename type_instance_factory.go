package goldi

import "fmt"

// TypeInstanceFactory is a trivial implementation of the TypeFactory interface.
// It will always `generate` the same instance of some previously instantiated type.
type TypeInstanceFactory struct {

	// The Instance that this factory is going to return on each call to Generate
	Instance interface{}
}

// NewTypeInstanceFactory creates a new TypeInstanceFactory which will return the given instance on each call to Generate.
// It will panic if the given instance is nil
func NewTypeInstanceFactory(instance interface{}) *TypeInstanceFactory {
	if instance == nil {
		panic(fmt.Errorf("refused to create a new TypeInstanceFactory with instance being nil"))
	}

	return &TypeInstanceFactory{instance}
}

// Generate fulfills the TypeFactory interface and will always return the type instance of this factory.
// It will panic if the instance is nil
func (f *TypeInstanceFactory) Generate(_ map[string]interface{}, _ TypeRegistry) interface{} {
	if f.Instance == nil {
		panic(fmt.Errorf("refused to return nil on TypeInstanceFactory.Generate. Seems like you did not use NewTypeInstanceFactory"))
	}

	return f.Instance
}
