package goldi

type aliasType struct {
	typeID string
}

// NewAliasType create a new TypeFactory which just serves as alias to the given type ID.
// A call to an alias type will retrieve the aliased type as if it was requested via container.Get(typeID)
// This method will always return a valid type and works bot for regular type references (without leading @) and
// references to type functions.
func NewAliasType(typeID string) TypeFactory {
	return &aliasType{typeID}
}

func (a *aliasType) Arguments() []interface{} {
	return []interface{}{"@" + a.typeID}
}

func (a *aliasType) Generate(resolver *ParameterResolver) (interface{}, error) {
	typeID := NewTypeID(a.typeID)
	if typeID.IsFuncReference {
		r := NewFuncReferenceType(typeID.ID, typeID.FuncReferenceMethod)
		return r.Generate(resolver)
	}

	return resolver.Container.Get(a.typeID)
}
