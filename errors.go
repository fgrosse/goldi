package goldi

import "fmt"

type TypeReferenceError struct {
	error
	TypeID       string
	TypeInstance interface{}
}

type UnknownTypeReferenceError struct {
	error
	TypeID string
}

func NewTypeReferenceError(typeID string, typeInstance interface{}, message string, printfParameters ...interface{}) TypeReferenceError {
	return TypeReferenceError{
		error:        fmt.Errorf(message, printfParameters...),
		TypeID:       typeID,
		TypeInstance: typeInstance,
	}
}

func NewUnknownTypeReferenceError(typeID, message string, printfParameters ...interface{}) UnknownTypeReferenceError {
	return UnknownTypeReferenceError{
		error:  fmt.Errorf(message, printfParameters...),
		TypeID: typeID,
	}
}
