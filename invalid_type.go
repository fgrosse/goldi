package goldi

// The invalidType is used to defer handling errors when a TypeFactory implementation is instantiated.
// Instead of returning an error or panicking the invalid type can be returned.
// It will usually be checked by the ContainerValidator or at least return an error when Generate is called
// on the corresponding type.
type invalidType struct {
	error
}

func newInvalidType(err error) *invalidType {
	return &invalidType{err}
}

func (t *invalidType) Generate(parameterResolver *ParameterResolver) (interface{}, error) {
	return nil, t.error
}

func (t *invalidType) Arguments() []interface{} {
	return []interface{}{}
}

// IsValid checks if a given type factory is valid.
// This function can be used to check the result of functions like NewType
func IsValid(t TypeFactory) bool {
	_, isInvalid := t.(*invalidType)
	return !isInvalid
}
