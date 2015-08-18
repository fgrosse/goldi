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
	return resolver.Container.Get(a.TypeID)
}
