package goldi

import "fmt"

// instanceType is a trivial implementation of the TypeFactory interface.
// It will always `generate` the same instance of some previously instantiated type.
type instanceType struct {

	// The instance that this factory is going to return on each call to Generate
	Instance interface{}
}

// NewInstanceType creates a new TypeFactory which will return the given instance on each call to Generate.
// It will return an invalid type factory if the given instance is nil
//
// You can not generate this type using goldigen
func NewInstanceType(instance interface{}) TypeFactory {
	if instance == nil {
		return newInvalidType(fmt.Errorf("refused to create a new InstanceType with instance being nil"))
	}

	return &instanceType{instance}
}

func (t *instanceType) Generate(_ *ParameterResolver) (interface{}, error) {
	return t.Instance, nil
}

func (t *instanceType) Arguments() []interface{} {
	return []interface{}{}
}
