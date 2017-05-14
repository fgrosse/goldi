package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Test Suite")
}

// The types below act as mocks in the tests

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

func NewMockTypeWithArgs(stringParameter string, boolParameter bool) *MockType {
	return &MockType{stringParameter, boolParameter}
}

type TypeForServiceInjection struct {
	InjectedType *MockType
}

func NewTypeForServiceInjection(injectedType *MockType) *TypeForServiceInjection {
	return &TypeForServiceInjection{injectedType}
}

type TypeForServiceInjectionMultiple struct {
	InjectedTypes []*TypeForServiceInjectionMultiple
}

func NewTypeForServiceInjectionMultipleArgs(injectedTypes ...*TypeForServiceInjectionMultiple) *TypeForServiceInjectionMultiple {
	t := &TypeForServiceInjectionMultiple{}
	t.InjectedTypes = injectedTypes
	return t
}
