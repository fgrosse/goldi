package goldi

type InvalidType struct {
	Err error
}

func NewInvalidType(err error) *InvalidType {
	return &InvalidType{err}
}

func (t *InvalidType) Generate(parameterResolver *ParameterResolver) (interface{}, error) {
	return nil, t.Err
}

// Arguments is part of the TypeFactory interface and does always return an empty list for the InvalidType.
func (t *InvalidType) Arguments() []interface{} {
	return []interface{}{}
}
