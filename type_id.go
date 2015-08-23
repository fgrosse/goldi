package goldi

import "strings"

type TypeID struct {
	ID, Raw             string
	FuncReferenceMethod string
	IsOptional          bool
	IsFuncReference     bool
}

func NewTypeID(s string) (t TypeID) {
	if s == "" {
		panic("can not create typeID from empty string")
	}

	t.ID = s
	t.Raw = s

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

	return
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
