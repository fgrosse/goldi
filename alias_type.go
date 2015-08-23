package goldi

type AliasType struct {
	TypeID string
}

func NewAliasType(typeID string) *AliasType {
	return &AliasType{typeID}
}

func (a *AliasType) Arguments() []interface{} {
	return []interface{}{"@" + a.TypeID}
}

func (a *AliasType) Generate(resolver *ParameterResolver) (interface{}, error) {
	typeID := NewTypeID(a.TypeID)
	if typeID.IsFuncReference {
		r := NewFuncReferenceType(typeID.ID, typeID.FuncReferenceMethod)
		return r.Generate(resolver)
	}

	return resolver.Container.Get(a.TypeID), nil
}
