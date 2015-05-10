package matchers

import (
	"fmt"
	"github.com/fgrosse/goldi/tests/testAPI"
	"github.com/onsi/gomega/types"
	"strings"
)

func ContainCode(expected string) types.GomegaMatcher {
	return &ContainCodeMatcher{ExpectedCode: testAPI.Unindent(expected)}
}

type ContainCodeMatcher struct {
	BeValidGoCodeMatcher
	ExpectedCode string
}

func (m *ContainCodeMatcher) Match(actual interface{}) (success bool, err error) {
	isCompilable, err := m.BeValidGoCodeMatcher.Match(actual)
	if isCompilable == false || err != nil {
		return isCompilable, err
	}

	return strings.Contains(string(m.source), m.ExpectedCode), nil
}

func (m *ContainCodeMatcher) FailureMessage(actual interface{}) (message string) {
	if m.outIsEmpty || m.compileError != nil {
		return m.BeValidGoCodeMatcher.FailureMessage(actual)
	}

	return fmt.Sprintf("Expected output:\n%s\nto contain the code:\n%s", m.indentSource(), m.ExpectedCode)
}

func (m *ContainCodeMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to contain the code\n%s", m.indentSource(), m.ExpectedCode)
}
