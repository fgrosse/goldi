package tests

import (
	"fmt"
	"strings"
)

type Foo struct{}

type Bar struct{}

type Baz struct {
	Parameter1, Parameter2 string
}

func NewFoo() *Foo {
	return &Foo{}
}

func NewBar() *Bar {
	return &Bar{}
}

func NewBaz(parameter1, parameter2 string) *Baz {
	return &Baz{parameter1, parameter2}
}

type MockType struct {
	StringParameter string
	BoolParameter   bool
}

func (t *MockType) DoStuff() string {
	return "I did stuff"
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

type someFunc func(string) string

func NewMockTypeFromStringFunc(s string, sf someFunc) *MockType {
	return &MockType{StringParameter: sf(s)}
}

func NewVariadicMockTypeFuncs(funcs ...someFunc) *MockType {
	m := &MockType{}
	for _, f := range funcs {
		m.StringParameter = f(m.StringParameter)
	}

	return m
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

type MockTypeConfigurator struct {
	StringParameter string
}

func NewMockTypeConfigurator(configuredStringParam string) *MockTypeConfigurator {
	return &MockTypeConfigurator{StringParameter: configuredStringParam}
}

func (c *MockTypeConfigurator) Configure(m *MockType) {
	m.StringParameter = c.StringParameter
}

type FailingMockTypeConfigurator struct {
	StringParameter string
}

func NewFailingMockTypeConfigurator() *FailingMockTypeConfigurator {
	return &FailingMockTypeConfigurator{}
}

func (c *FailingMockTypeConfigurator) Configure(m *MockType) error {
	return fmt.Errorf("this is the error message from the tests.MockTypeConfigurator")
}
