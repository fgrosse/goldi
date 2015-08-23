package factories

import "github.com/fgrosse/goldi"

type InvalidType struct {
	Err error
}

func NewInvalidType(err error) *InvalidType {
	return &InvalidType{err}
}

// Arguments is part of the TypeFactory interface and does always return an empty list for the InstanceType.
func (t *InvalidType) Arguments() (args []interface{}) { return }

func (t *InvalidType) Generate(parameterResolver *goldi.ParameterResolver) (interface{}, error) {
	return nil, t.Err
}
