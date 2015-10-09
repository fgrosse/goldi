package goldi

import "strings"

// TypeID represents a parsed type identifier and associated meta data.
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

// String implements the fmt.Stringer interface by returning the raw representation of this type ID.
func (t *TypeID) String() string {
	if t.Raw != "" {
		return t.Raw
	}

	if t.FuncReferenceMethod != "" {
		return "@" + t.ID + "::" + t.FuncReferenceMethod
	}

	return "@" + t.ID
}

// IsParameterOrTypeReference is a utility function that returns whether the given string represents a parameter or a reference to a type.
// See IsParameter and IsTypeReference for further details
func IsParameterOrTypeReference(p string) bool {
	return IsParameter(p) || IsTypeReference(p)
}

// IsParameter returns whether the given type ID represents a parameter.
// A goldi parameter is recognized by the leading and trailing percent sign
// Example: %foobar%
func IsParameter(p string) bool {
	if len(p) < 3 {
		return false
	}

	return p[0] == '%' && p[len(p)-1] == '%'
}

// IsTypeReference returns whether the given string represents a reference to a type.
// A goldi type reference is recognized by the leading @ sign.
// Example: @foobar
func IsTypeReference(p string) bool {
	if len(p) < 2 {
		return false
	}

	return p[0] == '@'
}
