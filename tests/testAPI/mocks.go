package testAPI

import "strings"

type MockType struct {
	StringParameter string
	BoolParameter   bool
}

func (t *MockType) ReturnString(suffix string) string {
	return t.StringParameter + " " + suffix
}

func NewMockType() *MockType {
	return &MockType{}
}

func NewMockTypeWithArgs(stringParameter string, boolParameter bool) *MockType {
	return &MockType{stringParameter, boolParameter}
}

func NewVariadicMockType(foo bool, bar string, parameters ...string) *MockType {
	return &MockType{
		BoolParameter:   foo,
		StringParameter: strings.Join(parameters, ", "),
	}
}

type MockTypeFactory struct {
	HasBeenUsed bool
}

func (g *MockTypeFactory) NewMockType() *MockType {
	g.HasBeenUsed = true
	return &MockType{}
}

type TypeForServiceInjection struct {
	InjectedType *MockType
}

func NewTypeForServiceInjection(injectedType *MockType) *TypeForServiceInjection {
	return &TypeForServiceInjection{injectedType}
}

func NewTypeForServiceInjectionWithArgs(injectedType *MockType, name, location string, flag bool) *TypeForServiceInjection {
	return &TypeForServiceInjection{injectedType}
}
