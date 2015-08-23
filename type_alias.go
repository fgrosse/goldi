package goldi

type TypeAlias struct {
	TypeID string
}

func NewTypeAlias(typeID string) *TypeAlias {
	return &TypeAlias{typeID}
}

func (a *TypeAlias) Arguments() []interface{} {
	return []interface{}{"@" + a.TypeID}
}

func (a *TypeAlias) Generate(resolver *ParameterResolver) interface{} {
	typeID := newTypeId(a.TypeID)
	if typeID.IsFuncReference {
		r := NewFuncReferenceType(typeID.ID, typeID.FuncReferenceMethod)
		return r.Generate(resolver)
	}

	return resolver.Container.Get(a.TypeID)
}
