package goldi

import "strings"

type TypeID struct {
	ID, Raw             string
	FuncReferenceMethod string
	IsOptional          bool
	IsFuncReference     bool
}

// NewTypeID creates a new TypeId. Trying to create a type ID from an empty string will panic
func NewTypeID(s string) *TypeID {
	if s == "" {
		panic("can not create typeID from empty string")
	}

	t := &TypeID{
		ID:  s,
		Raw: s,
	}

	if t.ID[0] == '@' {
		t.ID = t.ID[1:]
	}

	if t.ID[0] == '?' {
		t.IsOptional = true
		t.ID = t.ID[1:]
	}

	funcReferenceParts := strings.SplitN(t.ID, "::", 2)
	if len(funcReferenceParts) == 2 {
		t.IsFuncReference = true
		t.ID = funcReferenceParts[0]
		t.FuncReferenceMethod = funcReferenceParts[1]
	}

	return t
}

func (t *TypeID) String() string {
	if t.Raw != "" {
		return t.Raw
	}

	if t.FuncReferenceMethod != "" {
		return "@" + t.ID + "::" + t.FuncReferenceMethod
	}

	return "@" + t.ID
}

func isParameterOrTypeReference(p string) bool {
	return isParameter(p) || isTypeReference(p)
}

func isParameter(p string) bool {
	if len(p) < 3 {
		return false
	}

	return p[0] == '%' && p[len(p)-1] == '%'
}

func isTypeReference(p string) bool {
	if len(p) < 2 {
		return false
	}

	return p[0] == '@'
}
