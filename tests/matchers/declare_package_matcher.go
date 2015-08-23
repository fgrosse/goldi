package matchers

import (
	"fmt"
	"github.com/onsi/gomega/types"
)

func DeclarePackage(packageName string) types.GomegaMatcher {
	return &DeclarePackageMatcher{ExpectedPackageName: packageName}
}

type DeclarePackageMatcher struct {
	BeValidGoCodeMatcher
	ExpectedPackageName string
}

func (m *DeclarePackageMatcher) Match(actual interface{}) (success bool, err error) {
	isCompilable, err := m.BeValidGoCodeMatcher.Match(actual)
	if isCompilable == false || err != nil {
		return isCompilable, err
	}

	return m.astFile.Name.Name == m.ExpectedPackageName, nil
}

func (m *DeclarePackageMatcher) FailureMessage(actual interface{}) (message string) {
	if m.outIsEmpty || m.compileError != nil {
		return m.BeValidGoCodeMatcher.FailureMessage(actual)
	}

	return fmt.Sprintf("Expected output:\n%s\nto be in package %q", m.indentSource(), m.ExpectedPackageName)
}

func (m *DeclarePackageMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to to be in package %q", m.indentSource(), m.ExpectedPackageName)
}
