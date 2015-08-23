package matchers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/onsi/gomega/types"
)

func ContainCode(expected string) types.GomegaMatcher {
	return &ContainCodeMatcher{ExpectedCode: unindent(expected)}
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

	return fmt.Sprintf("Expected output:\n%s\nto contain the code:\n%s", m.indentSource(), indent([]byte(m.ExpectedCode)))
}

func (m *ContainCodeMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to contain the code\n%s", m.indentSource(), indent([]byte(m.ExpectedCode)))
}

func unindent(input string) string {
	indent := getFirstIndent(input)
	result := &bytes.Buffer{}
	indentIndex := 0
	for _, c := range input {
		if c == '\n' {
			indentIndex = 0
		} else if indentIndex < len(indent) {
			indentIndex++
			continue
		}

		result.WriteRune(c)
	}

	return result.String()
}

func getFirstIndent(input string) string {
	// search for the first real indent
	indent := &bytes.Buffer{}
	for _, c := range input {
		switch c {
		case ' ':
			fallthrough
		case '\t':
			indent.WriteRune(c)
		case '\n':
			indent.Reset()
		default:
			return indent.String()
		}
	}

	return indent.String()
}
