package matchers

import (
	"fmt"
	"github.com/onsi/gomega/types"
)

func ImportPackage(expected string) types.GomegaMatcher {
	return &ImportPackageMatcher{ExpectedPackage: fmt.Sprintf("%q", expected)}
}

type ImportPackageMatcher struct {
	BeValidGoCodeMatcher
	ExpectedPackage string
}

func (m *ImportPackageMatcher) Match(actual interface{}) (success bool, err error) {
	isCompilable, err := m.BeValidGoCodeMatcher.Match(actual)
	if isCompilable == false || err != nil {
		return isCompilable, err
	}

	for _, importSpec := range m.astFile.Imports {
		if importSpec.Path.Value == m.ExpectedPackage {
			return true, nil
		}
	}
	return false, nil
}

func (m *ImportPackageMatcher) FailureMessage(actual interface{}) (message string) {
	if m.outIsEmpty || m.compileError != nil {
		return m.BeValidGoCodeMatcher.FailureMessage(actual)
	}

	return fmt.Sprintf("Expected output:\n%s\nto import package %q", m.indentSource(), m.ExpectedPackage)
}

func (m *ImportPackageMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to import package %q", m.indentSource(), m.ExpectedPackage)
}
