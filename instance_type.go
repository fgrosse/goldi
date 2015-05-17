package goldi

import "fmt"

// InstanceType is a trivial implementation of the TypeFactory interface.
// It will always `generate` the same instance of some previously instantiated type.
// StructType implements the TypeFactory interface.
type InstanceType struct {

	// The Instance that this factory is going to return on each call to Generate
	Instance interface{}
}

// NewInstanceType creates a new InstanceType which will return the given instance on each call to Generate.
// It will panic if the given instance is nil
func NewInstanceType(instance interface{}) *InstanceType {
	if instance == nil {
		panic(fmt.Errorf("refused to create a new InstanceType with instance being nil"))
	}

	return &InstanceType{instance}
}

// Generate fulfills the TypeFactory interface and will always return the type instance of this factory.
// It will panic if the instance is nil
func (t *InstanceType) Generate(_ map[string]interface{}, _ TypeRegistry) interface{} {
	if t.Instance == nil {
		panic(fmt.Errorf("refused to return nil on InstanceType.Generate. Seems like you did not use NewInstanceType"))
	}

	return t.Instance
}

// Arguments is part of the TypeFactory interface and does always return an empty list for the InstanceType.
func (t *InstanceType) Arguments() []interface{} {
	return []interface{}{}
}
