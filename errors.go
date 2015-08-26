package goldi

import "fmt"

// A TypeReferenceError occurs if you tried to inject a type that does not match the function declaration of the corresponding method.
type TypeReferenceError struct {
	error
	TypeID       string
	TypeInstance interface{}
}

// The UnknownTypeReferenceError occurs if you try to get a type by an unknown type id (type has not been registered).
type UnknownTypeReferenceError struct {
	error
	TypeID string
}

// newTypeReferenceError creates a new TypeReferenceError
func newTypeReferenceError(typeID string, typeInstance interface{}, message string, printfParameters ...interface{}) TypeReferenceError {
	return TypeReferenceError{
		error:        fmt.Errorf(message, printfParameters...),
		TypeID:       typeID,
		TypeInstance: typeInstance,
	}
}

// newUnknownTypeReferenceError creates a new UnknownTypeReferenceError
func newUnknownTypeReferenceError(typeID, message string, printfParameters ...interface{}) UnknownTypeReferenceError {
	return UnknownTypeReferenceError{
		error:  fmt.Errorf(message, printfParameters...),
		TypeID: typeID,
	}
}
